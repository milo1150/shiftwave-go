package handler

import (
	"shiftwave-go/internal/middleware"
	"shiftwave-go/internal/types"

	"github.com/casbin/casbin/v2"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, app *types.App, enforcer *casbin.Enforcer) {
	// Swagger endpoint
	e.GET("/swagger/*", echoSwagger.WrapHandler, middleware.ValidateJwt(e, app))

	e.GET("/generate-pdf", func(c echo.Context) error {
		return GenerateQRCodeHandler(c, app.DB)
	}, middleware.JwtAndPermissionMiddlewares(e, app, enforcer)...)

	e.GET("/generate-random-reviews", func(c echo.Context) error {
		return GenerateRandomReviews(c, app)
	}, middleware.JwtAndPermissionMiddlewares(e, app, enforcer)...)
}
