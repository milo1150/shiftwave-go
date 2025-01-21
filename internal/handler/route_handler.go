package handler

import (
	"shiftwave-go/internal/types"

	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, app *types.App) {
	// Swagger endpoint
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/generate-pdf", func(c echo.Context) error {
		return GenerateQRCodeHandler(c)
	})

	e.GET("/generate-random-reviews", func(c echo.Context) error {
		return GenerateRandomReviews(c, app)
	})
}
