package handler

import (
	"shiftwave-go/internal/types"
	v1handler "shiftwave-go/internal/v1/handler"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, app *types.App) {
	e.GET("/generate-pdf", func(c echo.Context) error {
		return GenerateQRCodeHandler(c)
	})

	e.GET("/generate-random-reviews", func(c echo.Context) error {
		return GenerateRandomReviews(c, app)
	})

	e.POST("/v1/login", func(c echo.Context) error {
		return v1handler.LoginHandler(c, app)
	})

	e.POST("/v1/user", func(c echo.Context) error {
		return v1handler.CreateUser(c, app)
	})
}
