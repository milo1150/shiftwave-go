package handler

import (
	"context"
	"shiftwave-go/internal/middleware"
	"shiftwave-go/internal/types"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, app *types.App, ctx context.Context) {
	e.POST("/review", func(ctx echo.Context) error {
		return CreateReviewHandler(ctx, app.DB)
	}, middleware.IpRateLimiterMiddleware(app.RDB, ctx, 1))

	e.GET("/reviews", func(ctx echo.Context) error {
		return GetReviewsHandler(ctx, app)
	})

	e.GET("/review/:id", func(ctx echo.Context) error {
		return GetReviewHandler(ctx, app)
	})

	e.GET("/generate-pdf", func(ctx echo.Context) error {
		return GenerateQRCodeHandler(ctx)
	})

	e.GET("/generate-random-reviews", func(c echo.Context) error {
		return GenerateRandomReviews(c, app.DB)
	})
}
