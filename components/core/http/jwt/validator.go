package jwt

import (
	"fmt"
	"strings"

	"github.com/datarhei/core/v16/http/api"
	"github.com/datarhei/core/v16/http/handler/util"
	"github.com/datarhei/core/v16/http/jwt/jwks"

	jwtgo "github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type Validator interface {
	String() string

	// Validate returns true if it identified itself as validator for
	// that request. False if it doesn't handle this request. The string
	// is the username. An error is only returned if it identified itself
	// as validator but there was an error during validation.
	Validate(c echo.Context) (bool, string, error)
	Cancel()
}

type localValidator struct {
	username string
	password string
}

func NewLocalValidator(username, password string) (Validator, error) {
	v := &localValidator{
		username: username,
		password: password,
	}

	return v, nil
}

func (v *localValidator) String() string {
	return "localjwt"
}

func (v *localValidator) Validate(c echo.Context) (bool, string, error) {
	var login api.Login

	if err := util.ShouldBindJSON(c, &login); err != nil {
		return false, "", nil
	}

	if login.Username != v.username || login.Password != v.password {
		return true, "", fmt.Errorf("invalid username or password")
	}

	return true, v.username, nil
}

func (v *localValidator) Cancel() {}

type auth0Validator struct {
	domain   string
	issuer   string
	audience string
	clientID string
	users    []string
	certs    jwks.JWKS
}

func NewAuth0Validator(domain, audience, clientID string, users []string) (Validator, error) {
	v := &auth0Validator{
		domain:   domain,
		issuer:   "https://" + domain + "/",
		audience: audience,
		clientID: clientID,
		users:    users,
	}

	url := v.issuer + ".well-known/jwks.json"
	certs, err := jwks.NewFromURL(url, jwks.Config{})
	if err != nil {
		return nil, err
	}

	v.certs = certs

	return v, nil
}

func (v auth0Validator) String() string {
	return fmt.Sprintf("auth0 domain=%s audience=%s clientid=%s", v.domain, v.audience, v.clientID)
}

func (v *auth0Validator) Validate(c echo.Context) (bool, string, error) {
	// Look for an Auth header
	values := c.Request().Header.Values("Authorization")
	prefix := "Bearer "

	auth := ""
	for _, value := range values {
		if !strings.HasPrefix(value, prefix) {
			continue
		}

		auth = value[len(prefix):]

		break
	}

	if len(auth) == 0 {
		return false, "", nil
	}

	p := &jwtgo.Parser{}
	token, _, err := p.ParseUnverified(auth, jwtgo.MapClaims{})
	if err != nil {
		return false, "", nil
	}

	var issuer string
	if claims, ok := token.Claims.(jwtgo.MapClaims); ok {
		if iss, ok := claims["iss"]; ok {
			issuer = iss.(string)
		}
	}

	if issuer != v.issuer {
		return false, "", nil
	}

	token, err = jwtgo.Parse(auth, v.keyFunc)
	if err != nil {
		return true, "", err
	}

	if !token.Valid {
		return true, "", fmt.Errorf("invalid token")
	}

	var subject string
	if claims, ok := token.Claims.(jwtgo.MapClaims); ok {
		if sub, ok := claims["sub"]; ok {
			subject = sub.(string)
		}
	}

	return true, subject, nil
}

func (v *auth0Validator) keyFunc(token *jwtgo.Token) (interface{}, error) {
	// Verify 'aud' claim
	if _, err := token.Claims.GetAudience(); err != nil {
		return nil, fmt.Errorf("invalid audience: %w", err)
	}

	// Verify 'iss' claim
	if _, err := token.Claims.GetIssuer(); err != nil {
		return nil, fmt.Errorf("invalid issuer: %w", err)
	}

	// Verify 'sub' claim
	sub, err := token.Claims.GetSubject()
	if err != nil {
		return nil, fmt.Errorf("invalid subject: %w", err)
	}

	found := false
	for _, u := range v.users {
		if sub == u {
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("user not allowed")
	}

	// find the key
	if _, ok := token.Header["kid"]; !ok {
		return nil, fmt.Errorf("kid not found")
	}

	kid := token.Header["kid"].(string)

	key, err := v.certs.Key(kid)
	if err != nil {
		return nil, fmt.Errorf("no cert for kid found: %w", err)
	}

	// find algorithm
	if _, ok := token.Header["alg"]; !ok {
		return nil, fmt.Errorf("kid not found")
	}

	alg := token.Header["alg"].(string)

	if key.Alg() != alg {
		return nil, fmt.Errorf("signing method doesn't match")
	}

	// get the public key
	publicKey, err := key.PublicKey()
	if err != nil {
		return nil, fmt.Errorf("invalid public key: %w", err)
	}

	return publicKey, nil
}

func (v *auth0Validator) Cancel() {
	v.certs.Cancel()
}
