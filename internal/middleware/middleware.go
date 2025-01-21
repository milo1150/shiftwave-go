package middleware

import (
	"shiftwave-go/internal/types"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetupMiddlewares(e *echo.Echo, env types.Env) {
	e.Use(middleware.Secure())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(ConfigCORS())
	e.Use(ConfigRateLimiter())
}
