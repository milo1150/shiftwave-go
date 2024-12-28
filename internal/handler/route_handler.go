package handler

import (
	"shiftwave-go/internal/middleware"
	"shiftwave-go/internal/types"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, app *types.App) {
	e.POST("/rating", func(ctx echo.Context) error {
		return CreateRatingHandler(ctx, app)
	})
	e.GET("/ratings", func(ctx echo.Context) error {
		return GetRatingsHandler(ctx, app)
	})
	e.GET("/rating/:id", func(ctx echo.Context) error {
		if err := middleware.JWT(e, ctx); err != nil {
			return err
		}
		return GetRatingHandler(ctx, app)
	})
}
