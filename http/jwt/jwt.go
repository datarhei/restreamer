package jwt

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/datarhei/core/v16/app"
	"github.com/datarhei/core/v16/http/api"

	jwtgo "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// The Config type holds information that is required to create a new JWT provider
type Config struct {
	Realm         string
	Secret        string
	SkipLocalhost bool
}

// JWT provides access to a JWT provider
type JWT interface {
	AddValidator(iss string, issuer Validator) error

	ClearValidators()

	Validators() []string

	// Middleware returns an echo middleware
	AccessMiddleware() echo.MiddlewareFunc
	RefreshMiddleware() echo.MiddlewareFunc

	// LoginHandler is an echo route handler for retrieving a JWT
	LoginHandler(c echo.Context) error

	// RefreshHandle is an echo route handler for refreshing a JWT
	RefreshHandler(c echo.Context) error
}

type jwt struct {
	realm             string
	skipLocalhost     bool
	secret            []byte
	accessValidFor    time.Duration
	accessConfig      echojwt.Config
	accessMiddleware  echo.MiddlewareFunc
	refreshValidFor   time.Duration
	refreshConfig     echojwt.Config
	refreshMiddleware echo.MiddlewareFunc
	// Validators is a map of all recognized issuers to their specific validators. The key is the value of
	// the "iss" field in the claims. Somewhat required because otherwise the token cannot be verified.
	validators map[string]Validator
	lock       sync.RWMutex
}

// New returns a new JWT provider
func New(config Config) (JWT, error) {
	j := &jwt{
		realm:           config.Realm,
		skipLocalhost:   config.SkipLocalhost,
		secret:          []byte(config.Secret),
		accessValidFor:  time.Minute * 10,
		refreshValidFor: time.Hour * 24,
	}

	if len(j.secret) == 0 {
		return nil, fmt.Errorf("the JWT secret must not be empty")
	}

	skipperFunc := func(c echo.Context) bool {
		if j.skipLocalhost {
			ip := c.RealIP()

			if ip == "127.0.0.1" || ip == "::1" {
				return true
			}
		}

		return false
	}

	j.accessConfig = echojwt.Config{
		Skipper:        skipperFunc,
		SigningMethod:  middleware.AlgorithmHS256,
		ContextKey:     "user",
		TokenLookup:    "header:Authorization:Bearer ",
		ErrorHandler:   j.ErrorHandler,
		ParseTokenFunc: j.parseToken("access"),
	}

	j.refreshConfig = echojwt.Config{
		Skipper:        skipperFunc,
		SigningMethod:  middleware.AlgorithmHS256,
		ContextKey:     "user",
		TokenLookup:    "header:Authorization:Bearer ",
		ErrorHandler:   j.ErrorHandler,
		ParseTokenFunc: j.parseToken("refresh"),
	}

	return j, nil
}

func (j *jwt) parseToken(use string) func(c echo.Context, auth string) (interface{}, error) {
	keyFunc := func(*jwtgo.Token) (interface{}, error) { return j.secret, nil }

	return func(c echo.Context, auth string) (interface{}, error) {
		var token *jwtgo.Token
		var err error

		token, err = jwtgo.Parse(auth, keyFunc)
		if err != nil {
			return nil, err
		}

		if !token.Valid {
			return nil, errors.New("invalid token")
		}

		if _, ok := token.Claims.(jwtgo.MapClaims)["usefor"]; !ok {
			return nil, fmt.Errorf("usefor claim is required")
		}

		claimuse := token.Claims.(jwtgo.MapClaims)["usefor"].(string)

		if claimuse != use {
			return nil, fmt.Errorf("invalid token claim")
		}

		return token, nil
	}
}

func (j *jwt) Validators() []string {
	j.lock.RLock()
	defer j.lock.RUnlock()

	values := []string{}

	for _, v := range j.validators {
		values = append(values, v.String())
	}

	return values
}

func (j *jwt) AddValidator(iss string, issuer Validator) error {
	j.lock.Lock()
	defer j.lock.Unlock()

	if j.validators == nil {
		j.validators = make(map[string]Validator)
	}

	if _, ok := j.validators[iss]; ok {
		return fmt.Errorf("a validator for %s is already registered", iss)
	}

	j.validators[iss] = issuer

	return nil
}

func (j *jwt) ClearValidators() {
	j.lock.Lock()
	defer j.lock.Unlock()

	if j.validators == nil {
		return
	}

	for _, v := range j.validators {
		v.Cancel()
	}

	j.validators = nil
}

func (j *jwt) ErrorHandler(c echo.Context, err error) error {
	if c.Request().URL.Path == "/api" {
		return c.JSON(http.StatusOK, api.MinimalAbout{
			App:   app.Name,
			Auths: j.Validators(),
			Version: api.VersionMinimal{
				Number: app.Version.MajorString(),
			},
		})
	}

	return api.Err(http.StatusUnauthorized, "", "Missing or invalid JWT token")
}

func (j *jwt) AccessMiddleware() echo.MiddlewareFunc {
	if j.accessMiddleware == nil {
		j.accessMiddleware = echojwt.WithConfig(j.accessConfig)
	}

	return j.accessMiddleware
}

func (j *jwt) RefreshMiddleware() echo.MiddlewareFunc {
	if j.refreshMiddleware == nil {
		j.refreshMiddleware = echojwt.WithConfig(j.refreshConfig)
	}

	return j.refreshMiddleware
}

// LoginHandler returns an access token and a refresh token
// @Summary Retrieve an access and a refresh token
// @Description Retrieve valid JWT access and refresh tokens to use for accessing the API. Login either by username/password or Auth0 token
// @ID jwt-login
// @Produce json
// @Param data body api.Login true "Login data"
// @Success 200 {object} api.JWT
// @Failure 400 {object} api.Error
// @Failure 403 {object} api.Error
// @Failure 500 {object} api.Error
// @Security Auth0KeyAuth
// @Router /api/login [post]
func (j *jwt) LoginHandler(c echo.Context) error {
	var ok bool
	var subject string
	var err error

	j.lock.RLock()
	for _, validator := range j.validators {
		ok, subject, err = validator.Validate(c)
		if ok {
			break
		}
	}
	j.lock.RUnlock()

	if ok {
		if err != nil {
			time.Sleep(5 * time.Second)
			return api.Err(http.StatusUnauthorized, "", "Invalid authorization credentials: %s", err.Error())
		}
	} else {
		time.Sleep(5 * time.Second)
		return api.Err(http.StatusBadRequest, "", "Missing authorization credentials")
	}

	at, rt, err := j.createToken(subject)
	if err != nil {
		return api.Err(http.StatusInternalServerError, "", "Failed to create JWT: %s", err.Error())
	}

	return c.JSON(http.StatusOK, api.JWT{
		AccessToken:  at,
		RefreshToken: rt,
	})
}

// RefreshHandler returns a new refresh token
// @Summary Retrieve a new access token
// @Description Retrieve a new access token by providing the refresh token
// @ID jwt-refresh
// @Produce json
// @Success 200 {object} api.JWTRefresh
// @Failure 500 {object} api.Error
// @Security ApiRefreshKeyAuth
// @Router /api/login/refresh [get]
func (j *jwt) RefreshHandler(c echo.Context) error {
	token, ok := c.Get("user").(*jwtgo.Token)
	if !ok {
		return api.Err(http.StatusForbidden, "", "Invalid token")
	}

	subject, err := token.Claims.GetSubject()
	if err != nil {
		return api.Err(http.StatusForbidden, "", "Invalid subject: %s", err.Error())
	}

	at, _, err := j.createToken(subject)
	if err != nil {
		return api.Err(http.StatusInternalServerError, "", "Failed to create JWT: %s", err.Error())
	}

	return c.JSON(http.StatusOK, api.JWTRefresh{
		AccessToken: at,
	})
}

// Already assigned claims: https://www.iana.org/assignments/jwt/jwt.xhtml

func (j *jwt) createToken(username string) (string, string, error) {
	now := time.Now()
	accessExpires := now.Add(j.accessValidFor)
	refreshExpires := now.Add(j.refreshValidFor)

	// Create access token
	accessToken := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, jwtgo.MapClaims{
		"iss":    j.realm,
		"sub":    username,
		"usefor": "access",
		"iat":    now.Unix(),
		"exp":    accessExpires.Unix(),
		"exi":    uint64(accessExpires.Sub(now).Seconds()),
		"jti":    uuid.New().String(),
	})

	// Generate encoded access token
	at, err := accessToken.SignedString(j.secret)
	if err != nil {
		return "", "", err
	}

	// Create refresh token
	refreshToken := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, jwtgo.MapClaims{
		"iss":    j.realm,
		"sub":    username,
		"usefor": "refresh",
		"iat":    now.Unix(),
		"exp":    refreshExpires.Unix(),
		"exi":    uint64(refreshExpires.Sub(now).Seconds()),
		"jti":    uuid.New().String(),
	})

	// Generate encoded refresh token
	rt, err := refreshToken.SignedString(j.secret)
	if err != nil {
		return "", "", err
	}

	return at, rt, nil
}
