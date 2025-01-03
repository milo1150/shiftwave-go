package handler

import (
	"shiftwave-go/internal/types"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, app *types.App) {
	e.POST("/review", func(ctx echo.Context) error {
		return CreateReviewHandler(ctx, app)
	})
	e.GET("/reviews", func(ctx echo.Context) error {
		return GetReviewsHandler(ctx, app)
	})
	e.GET("/review/:id", func(ctx echo.Context) error {
		// if err := middleware.JWT(e, ctx); err != nil {
		// 	return err
		// }
		return GetReviewHandler(ctx, app)
	})
}
