package handler

import (
	"context"
	"shiftwave-go/internal/middleware"
	"shiftwave-go/internal/types"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

func SetupRoutes(e *echo.Echo, app *types.App, rdb *redis.Client, ctx context.Context) {
	e.POST("/review", func(ctx echo.Context) error {
		return CreateReviewHandler(ctx, app)
	}, middleware.IpRateLimiterMiddleware(rdb, ctx, 1))

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
