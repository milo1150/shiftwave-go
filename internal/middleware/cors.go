package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func ConfigCORS() echo.MiddlewareFunc {
	config := middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost:4321",
			"https://shiftwave-dev.mijio.app/",
		},
		// Do not allow delete beacuse we don't delete data here.
		AllowMethods: []string{
			echo.GET,
			echo.POST,
			echo.PATCH,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
		},
	})
	return config
}
