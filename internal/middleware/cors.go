package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func ConfigCORS() echo.MiddlewareFunc {
	config := middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:4321",
			"https://be-shiftwave-dev.mijio.app",
			"https://shiftwave-dev.mijio.app",
			"https://be-shiftwave-dev.mijio.app:8081",
		},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	})
	return config
}
