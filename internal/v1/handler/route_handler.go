package handler

import (
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
	userRoute(r.Echo, r.App, r.Enforcer)

	// /v1/reviews
	reviewsRoute(r.Echo, r.App, r.Enforcer)

	// /v1/review
	reviewRoute(r.Echo, r.App)

	// /v1/branch
	branchRoute(r.Echo, r.App, r.Enforcer)

	// /v1/branches
	branchesRoute(r.Echo, r.App)
}

func userRoute(e *echo.Echo, app *types.App, enforcer *casbin.Enforcer) {
	// Public
	e.POST("/v1/login", func(c echo.Context) error {
		return LoginHandler(c, app)
	})

	// Private
	userGroup := e.Group("/v1/user", middleware.AdminMiddlewares(e, app, enforcer)...)

	userGroup.GET("/get-users", func(c echo.Context) error {
		return GetAllUsersHandler(c, app)
	})

	userGroup.POST("/create-user", func(c echo.Context) error {
		return CreateUserHandler(c, app)
	})
}

func reviewsRoute(e *echo.Echo, app *types.App, enforcer *casbin.Enforcer) {
	reviewsGroup := e.Group("/v1/reviews", middleware.AdminMiddlewares(e, app, enforcer)...)

	reviewsGroup.GET("", func(c echo.Context) error {
		return GetReviewsHandler(c, app)
	})

	reviewsGroup.GET("/average-rating", func(c echo.Context) error {
		return GetAverageRatingHandler(c, app)
	})

	reviewsGroup.GET("/sse", func(c echo.Context) error {
		return ReviewSse(c)
	})

	// TODO: Grouping with JWT auth
	e.GET("/v1/reviews/s-ws", func(c echo.Context) error {
		return ReviewWsSingleConnection(c, app)
	})

	// TODO: Grouping with JWT auth
	e.GET("/v1/reviews/m-ws", func(c echo.Context) error {
		return ReviewWsMultipleConnection(c, app)
	})
}

func reviewRoute(e *echo.Echo, app *types.App) {
	reviewGroup := e.Group("/v1/review")

	reviewGroup.POST("", func(c echo.Context) error {
		return CreateReviewHandler(c, app.DB)
	}, middleware.IpRateLimiterMiddleware(app.RDB, 1))

	e.GET("/v1/review/limit", func(c echo.Context) error {
		return CheckDailyLimit(c, app.RDB)
	})
}

func branchRoute(e *echo.Echo, app *types.App, enforcer *casbin.Enforcer) {
	branchGroup := e.Group("/v1/branch", middleware.AdminMiddlewares(e, app, enforcer)...)

	branchGroup.POST("", func(c echo.Context) error {
		return CreateBranchHandler(c, app.DB)
	})

	branchGroup.PATCH("/:uuid", func(c echo.Context) error {
		return UpdateBranchHandler(c, app.DB)
	})
}

func branchesRoute(e *echo.Echo, app *types.App) {
	branchesGroup := e.Group("/v1/branches")

	branchesGroup.GET("", func(c echo.Context) error {
		return GetBranchesHandler(c, app.DB)
	})
}
