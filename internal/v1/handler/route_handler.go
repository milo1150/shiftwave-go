package handler

import (
	"shiftwave-go/internal/middleware"
	"shiftwave-go/internal/types"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, app *types.App) {
	e.GET("/v1/reviews", func(c echo.Context) error {
		return GetReviewsHandler(c, app)
	})

	e.GET("/v1/reviews/average-rating", func(c echo.Context) error {
		return GetAverageRatingHandler(c, app)
	})

	e.POST("/v1/review", func(c echo.Context) error {
		return CreateReviewHandler(c, app.DB)
	}, middleware.IpRateLimiterMiddleware(app.RDB, app.Context, 1))

	e.GET("/v1/review/:id", func(c echo.Context) error {
		return GetReviewHandler(c, app)
	})

	// TODO: permission
	e.POST("/v1/branch", func(c echo.Context) error {
		return CreateBranchHandler(c, app.DB)
	})

	// TODO: permission
	e.GET("/v1/branches", func(c echo.Context) error {
		return GetBranchesHandler(c, app.DB)
	})

	// TODO: permission
	e.PATCH("/v1/branch/:id", func(c echo.Context) error {
		return UpdateBranchHandler(c, app.DB)
	})

	e.GET("/v1/reviews/s-ws", func(c echo.Context) error {
		return ReviewWsSingleConnection(c, app)
	})

	e.GET("/v1/reviews/m-ws", func(c echo.Context) error {
		return ReviewWsMultipleConnection(c, app)
	})
}
