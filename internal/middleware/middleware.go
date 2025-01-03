package middleware

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
)

func SetupMiddlewares(e *echo.Echo, rdb *redis.Client, ctx context.Context) {
	e.Use(middleware.Secure())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(ConfigCORS())
	e.Use(ConfigRateLimiter())

	// Example of custom middleware
	// e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
	// 	return func(c echo.Context) error {
	// 		return next(c)
	// 	}
	// })
}
