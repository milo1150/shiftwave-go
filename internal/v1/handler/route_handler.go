package handler

import (
	"shiftwave-go/internal/middleware"
	"shiftwave-go/internal/types"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, app *types.App) {
	e.POST("/v1/review", func(c echo.Context) error {
		return CreateReviewHandler(c, app.DB)
	}, middleware.IpRateLimiterMiddleware(app.RDB, app.Context, 1))

	e.GET("/v1/reviews", func(c echo.Context) error {
		return GetReviewsHandler(c, app)
	})

	e.GET("/v1/review/:id", func(c echo.Context) error {
		return GetReviewHandler(c, app)
	})
}
