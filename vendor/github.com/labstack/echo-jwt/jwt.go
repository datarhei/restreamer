// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: © 2016 LabStack and Echo contributors

package echojwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

// Config defines the config for JWT middleware.
type Config struct {
	// Skipper defines a function to skip middleware.
	Skipper middleware.Skipper

	// BeforeFunc defines a function which is executed just before the middleware.
	BeforeFunc middleware.BeforeFunc

	// SuccessHandler defines a function which is executed for a valid token.
	SuccessHandler func(c echo.Context)

	// ErrorHandler defines a function which is executed when all lookups have been done and none of them passed Validator
	// function. ErrorHandler is executed with last missing (ErrExtractionValueMissing) or an invalid key.
	// It may be used to define a custom JWT error.
	//
	// Note: when error handler swallows the error (returns nil) middleware continues handler chain execution towards handler.
	// This is useful in cases when portion of your site/api is publicly accessible and has extra features for authorized users
	// In that case you can use ErrorHandler to set default public JWT token value to request and continue with handler chain.
	ErrorHandler func(c echo.Context, err error) error

	// ContinueOnIgnoredError allows the next middleware/handler to be called when ErrorHandler decides to
	// ignore the error (by returning `nil`).
	// This is useful when parts of your site/api allow public access and some authorized routes provide extra functionality.
	// In that case you can use ErrorHandler to set a default public JWT token value in the request context
	// and continue. Some logic down the remaining execution chain needs to check that (public) token value then.
	ContinueOnIgnoredError bool

	// Context key to store user information from the token into context.
	// Optional. Default value "user".
	ContextKey string

	// Signing key to validate token.
	// This is one of the three options to provide a token validation key.
	// The order of precedence is a user-defined KeyFunc, SigningKeys and SigningKey.
	// Required if neither user-defined KeyFunc nor SigningKeys is provided.
	SigningKey interface{}

	// Map of signing keys to validate token with kid field usage.
	// This is one of the three options to provide a token validation key.
	// The order of precedence is a user-defined KeyFunc, SigningKeys and SigningKey.
	// Required if neither user-defined KeyFunc nor SigningKey is provided.
	SigningKeys map[string]interface{}

	// Signing method used to check the token's signing algorithm.
	// Optional. Default value HS256.
	SigningMethod string

	// KeyFunc defines a user-defined function that supplies the public key for a token validation.
	// The function shall take care of verifying the signing algorithm and selecting the proper key.
	// A user-defined KeyFunc can be useful if tokens are issued by an external party.
	// Used by default ParseTokenFunc implementation.
	//
	// When a user-defined KeyFunc is provided, SigningKey, SigningKeys, and SigningMethod are ignored.
	// This is one of the three options to provide a token validation key.
	// The order of precedence is a user-defined KeyFunc, SigningKeys and SigningKey.
	// Required if neither SigningKeys nor SigningKey is provided.
	// Not used if custom ParseTokenFunc is set.
	// Default to an internal implementation verifying the signing algorithm and selecting the proper key.
	KeyFunc jwt.Keyfunc

	// TokenLookup is a string in the form of "<source>:<name>" or "<source>:<name>,<source>:<name>" that is used
	// to extract token from the request.
	// Optional. Default value "header:Authorization".
	// Possible values:
	// - "header:<name>" or "header:<name>:<cut-prefix>"
	// 			`<cut-prefix>` is argument value to cut/trim prefix of the extracted value. This is useful if header
	//			value has static prefix like `Authorization: <auth-scheme> <authorisation-parameters>` where part that we
	//			want to cut is `<auth-scheme> ` note the space at the end.
	//			In case of JWT tokens `Authorization: Bearer <token>` prefix we cut is `Bearer `.
	// If prefix is left empty the whole value is returned.
	// - "query:<name>"
	// - "param:<name>"
	// - "cookie:<name>"
	// - "form:<name>"
	// Multiple sources example:
	// - "header:Authorization:Bearer ,cookie:myowncookie"
	TokenLookup string

	// TokenLookupFuncs defines a list of user-defined functions that extract JWT token from the given context.
	// This is one of the two options to provide a token extractor.
	// The order of precedence is user-defined TokenLookupFuncs, and TokenLookup.
	// You can also provide both if you want.
	TokenLookupFuncs []middleware.ValuesExtractor

	// ParseTokenFunc defines a user-defined function that parses token from given auth. Returns an error when token
	// parsing fails or parsed token is invalid.
	// Defaults to implementation using `github.com/golang-jwt/jwt` as JWT implementation library
	ParseTokenFunc func(c echo.Context, auth string) (interface{}, error)

	// Claims are extendable claims data defining token content. Used by default ParseTokenFunc implementation.
	// Not used if custom ParseTokenFunc is set.
	// Optional. Defaults to function returning jwt.MapClaims
	NewClaimsFunc func(c echo.Context) jwt.Claims
}

const (
	// AlgorithmHS256 is token signing algorithm
	AlgorithmHS256 = "HS256"
)

// ErrJWTMissing denotes an error raised when JWT token value could not be extracted from request
var ErrJWTMissing = echo.NewHTTPError(http.StatusUnauthorized, "missing or malformed jwt")

// ErrJWTInvalid denotes an error raised when JWT token value is invalid or expired
var ErrJWTInvalid = echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")

// JWT returns a JSON Web Token (JWT) auth middleware.
//
// For valid token, it sets the user in context and calls next handler.
// For invalid token, it returns "401 - Unauthorized" error.
// For missing token, it returns "400 - Bad Request" error.
//
// See: https://jwt.io/introduction
func JWT(signingKey interface{}) echo.MiddlewareFunc {
	return WithConfig(Config{SigningKey: signingKey})
}

// WithConfig returns a JSON Web Token (JWT) auth middleware or panics if configuration is invalid.
//
// For valid token, it sets the user in context and calls next handler.
// For invalid token, it returns "401 - Unauthorized" error.
// For missing token, it returns "400 - Bad Request" error.
//
// See: https://jwt.io/introduction
func WithConfig(config Config) echo.MiddlewareFunc {
	mw, err := config.ToMiddleware()
	if err != nil {
		panic(err)
	}
	return mw
}

// ToMiddleware converts Config to middleware or returns an error for invalid configuration
func (config Config) ToMiddleware() (echo.MiddlewareFunc, error) {
	if config.Skipper == nil {
		config.Skipper = middleware.DefaultSkipper
	}
	if config.ContextKey == "" {
		config.ContextKey = "user"
	}
	if config.TokenLookup == "" && len(config.TokenLookupFuncs) == 0 {
		config.TokenLookup = "header:Authorization:Bearer "
	}
	if config.SigningMethod == "" {
		config.SigningMethod = AlgorithmHS256
	}

	if config.NewClaimsFunc == nil {
		config.NewClaimsFunc = func(c echo.Context) jwt.Claims {
			return jwt.MapClaims{}
		}
	}
	if config.SigningKey == nil && len(config.SigningKeys) == 0 && config.KeyFunc == nil && config.ParseTokenFunc == nil {
		return nil, errors.New("jwt middleware requires signing key")
	}
	if config.KeyFunc == nil {
		config.KeyFunc = config.defaultKeyFunc
	}
	if config.ParseTokenFunc == nil {
		config.ParseTokenFunc = config.defaultParseTokenFunc
	}
	extractors, err := middleware.CreateExtractors(config.TokenLookup)
	if err != nil {
		return nil, err
	}
	if len(config.TokenLookupFuncs) > 0 {
		extractors = append(config.TokenLookupFuncs, extractors...)
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			if config.BeforeFunc != nil {
				config.BeforeFunc(c)
			}
			var lastExtractorErr error
			var lastTokenErr error
			for _, extractor := range extractors {
				auths, extrErr := extractor(c)
				if extrErr != nil {
					lastExtractorErr = extrErr
					continue
				}
				for _, auth := range auths {
					token, err := config.ParseTokenFunc(c, auth)
					if err != nil {
						lastTokenErr = err
						continue
					}
					// Store user information from token into context.
					c.Set(config.ContextKey, token)
					if config.SuccessHandler != nil {
						config.SuccessHandler(c)
					}
					return next(c)
				}
			}

			// prioritize token errors over extracting errors
			err = lastTokenErr
			if err == nil {
				err = lastExtractorErr
			}
			if config.ErrorHandler != nil {
				tmpErr := config.ErrorHandler(c, err)
				if config.ContinueOnIgnoredError && tmpErr == nil {
					return next(c)
				}
				return tmpErr
			}
			if lastTokenErr == nil {
				return ErrJWTMissing.WithInternal(err)
			}
			return ErrJWTInvalid.WithInternal(err)
		}
	}, nil
}

func (config Config) defaultKeyFunc(token *jwt.Token) (interface{}, error) {
	if token.Method.Alg() != config.SigningMethod {
		return nil, fmt.Errorf("unexpected jwt signing method=%v", token.Header["alg"])
	}
	if len(config.SigningKeys) == 0 {
		return config.SigningKey, nil
	}

	if kid, ok := token.Header["kid"].(string); ok {
		if key, ok := config.SigningKeys[kid]; ok {
			return key, nil
		}
	}
	return nil, fmt.Errorf("unexpected jwt key id=%v", token.Header["kid"])
}

// defaultParseTokenFunc creates JWTGo implementation for ParseTokenFunc
func (config Config) defaultParseTokenFunc(c echo.Context, auth string) (interface{}, error) {
	token, err := jwt.ParseWithClaims(auth, config.NewClaimsFunc(c), config.KeyFunc)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return token, nil
}
