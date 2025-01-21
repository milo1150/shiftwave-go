package handler

import (
	"shiftwave-go/internal/auth"
	"shiftwave-go/internal/middleware"
	"shiftwave-go/internal/types"

	"github.com/casbin/casbin/v2"
	"github.com/labstack/echo/v4"
)

// TODO: refactor
func SetupRoutes(e *echo.Echo, app *types.App, enforcer *casbin.Enforcer) {
	e.POST("/v1/login", func(c echo.Context) error {
		return LoginHandler(c, app)
	})

	e.POST("/v1/user", func(c echo.Context) error {
		return CreateUser(c, app)
	})

	// Group: /v1/reviews
	reviewsGroup := e.Group("/v1/reviews", auth.Jwt(e, app.ENV), middleware.RoutePermission(app.ENV.JWT, enforcer))

	reviewsGroup.GET("", func(c echo.Context) error {
		return GetReviewsHandler(c, app)
	})

	reviewsGroup.GET("/average-rating", func(c echo.Context) error {
		return GetAverageRatingHandler(c, app)
	})

	// TODO: Grouping with JWT auth
	e.GET("/v1/reviews/s-ws", func(c echo.Context) error {
		return ReviewWsSingleConnection(c, app)
	})
	e.GET("/v1/reviews/m-ws", func(c echo.Context) error {
		return ReviewWsMultipleConnection(c, app)
	})

	// Group: /v1/review
	reviewGroup := e.Group("/v1/review")

	reviewGroup.POST("", func(c echo.Context) error {
		return CreateReviewHandler(c, app.DB)
	}, middleware.IpRateLimiterMiddleware(app.RDB, app.Context, 1))

	reviewGroup.GET("/:id", func(c echo.Context) error {
		return GetReviewHandler(c, app)
	})

	// Group: /v1/branch
	branchGroup := e.Group("/v1/branch")

	branchGroup.POST("", func(c echo.Context) error {
		return CreateBranchHandler(c, app.DB)
	})

	branchGroup.PATCH("/:id", func(c echo.Context) error {
		return UpdateBranchHandler(c, app.DB)
	})

	// Group: /v1/branches
	branchesGroup := e.Group("/v1/branches")

	branchesGroup.GET("/v1/branches", func(c echo.Context) error {
		return GetBranchesHandler(c, app.DB)
	})
}
