package middleware

import (
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func JWT(e *echo.Echo, c echo.Context) error {
	e.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte("secret"),
	}))

	token, ok := c.Get("user").(*jwt.Token) // by default token is stored under `user` key
	if !ok {
		return errors.New("JWT token missing or invalid")
	}

	claims, ok := token.Claims.(jwt.MapClaims) // by default claims is of type `jwt.MapClaims`
	if !ok {
		return errors.New("failed to cast claims as jwt.MapClaims")
	}

	return c.JSON(http.StatusOK, claims)
}
