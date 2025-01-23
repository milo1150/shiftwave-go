package handler

import (
	"shiftwave-go/internal/auth"
	"shiftwave-go/internal/middleware"
	"shiftwave-go/internal/types"

	"github.com/casbin/casbin/v2"
	"github.com/labstack/echo/v4"
)

type RouteV1 struct {
	Echo     *echo.Echo
	App      *types.App
	Enforcer *casbin.Enforcer
}

func (r *RouteV1) SetupRoutes() {
	// /v1/user
	userRoute(r.Echo, r.App)

	// /v1/reviews
	reviewsRoute(r.Echo, r.App, r.Enforcer)

	// /v1/review
	reviewRoute(r.Echo, r.App)

	// /v1/branch
	branchRoute(r.Echo, r.App)

	// /v1/branches
	branchesRoute(r.Echo, r.App)
}

func userRoute(e *echo.Echo, app *types.App) {
	userGroup := e.Group("/v1/user")

	userGroup.POST("/login", func(c echo.Context) error {
		return LoginHandler(c, app)
	})

	userGroup.POST("/create-user", func(c echo.Context) error {
		return CreateUser(c, app)
	})
}

func reviewsRoute(e *echo.Echo, app *types.App, enforcer *casbin.Enforcer) {
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

	// TODO: JWT auth
	e.GET("/v1/reviews/sse", func(c echo.Context) error {
		return ReviewSse(c)
	})
}

func reviewRoute(e *echo.Echo, app *types.App) {
	reviewGroup := e.Group("/v1/review")

	reviewGroup.POST("", func(c echo.Context) error {
		return CreateReviewHandler(c, app.DB)
	}, middleware.IpRateLimiterMiddleware(app.RDB, app.Context, 1))

	reviewGroup.GET("/:id", func(c echo.Context) error {
		return GetReviewHandler(c, app)
	})
}

func branchRoute(e *echo.Echo, app *types.App) {
	branchGroup := e.Group("/v1/branch")

	branchGroup.POST("", func(c echo.Context) error {
		return CreateBranchHandler(c, app.DB)
	})

	branchGroup.PATCH("/:id", func(c echo.Context) error {
		return UpdateBranchHandler(c, app.DB)
	})
}

func branchesRoute(e *echo.Echo, app *types.App) {
	branchesGroup := e.Group("/v1/branches")

	branchesGroup.GET("/v1/branches", func(c echo.Context) error {
		return GetBranchesHandler(c, app.DB)
	})
}
