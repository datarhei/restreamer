# Echo JWT middleware

JWT middleware for [Echo](https://github.com/labstack/echo) framework. By default uses [golang-jwt/jwt/v4](https://github.com/golang-jwt/jwt) 
as JWT implementation.

## Usage

Add JWT middleware dependency with go modules
```bash
go get github.com/labstack/echo-jwt
```

Use as import statement
```go
import "github.com/labstack/echo-jwt"
```

Add middleware in simplified form, by providing only the secret key
```go
e.Use(echojwt.JWT([]byte("secret")))
```

Add middleware with configuration options
```go
e.Use(echojwt.WithConfig(echojwt.Config{
  // ...
  SigningKey:             []byte("secret"),
  // ...
}))
```

Extract token in handler
```go
e.GET("/", func(c echo.Context) error {
  token, ok := c.Get("user").(*jwt.Token) // by default token is stored under `user` key
  if !ok {
    return errors.New("JWT token missing or invalid")
  }
  claims, ok := token.Claims.(jwt.MapClaims) // by default claims is of type `jwt.MapClaims`
  if !ok {
    return errors.New("failed to cast claims as jwt.MapClaims")
  }
  return c.JSON(http.StatusOK, claims)
})
```

## Full example

```go
package main

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte("secret"),
	}))

	e.GET("/", func(c echo.Context) error {
		token, ok := c.Get("user").(*jwt.Token) // by default token is stored under `user` key
		if !ok {
			return errors.New("JWT token missing or invalid")
		}
		claims, ok := token.Claims.(jwt.MapClaims) // by default claims is of type `jwt.MapClaims`
		if !ok {
			return errors.New("failed to cast claims as jwt.MapClaims")
		}
		return c.JSON(http.StatusOK, claims)
	})

	if err := e.Start(":8080"); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
```

Test with
```bash
curl -v -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ" http://localhost:8080
```

Output should be
```bash
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET / HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.81.0
> Accept: */*
> Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=UTF-8
< Date: Sun, 27 Nov 2022 21:34:17 GMT
< Content-Length: 52
< 
{"admin":true,"name":"John Doe","sub":"1234567890"}
```
