package main

import (
	"context"
	"shiftwave-go/internal/auth"
	"shiftwave-go/internal/database"
	baseHandler "shiftwave-go/internal/handler"
	"shiftwave-go/internal/middleware"
	"shiftwave-go/internal/scheduler"
	"shiftwave-go/internal/setup"
	"shiftwave-go/internal/types"
	v1handler "shiftwave-go/internal/v1/handler"

	_ "shiftwave-go/docs" // Import Swagger docs package

	"github.com/labstack/echo/v4"
)

// @title Shiftwave API
// @version 1.0
// @description This is a sample API documentation for Echo with Swagger.
// @BasePath /v1
func main() {
	// Load and construct env
	env := setup.EnvLoader()

	// Initialize the database
	db := database.InitDatabase()

	// Initialize Redis client
	rdb := database.NewRedisClient(env.RedisPassword)

	// Initialize Permission (Casbin)
	enforcer := auth.InitPermission(db)

	// Create application context
	ctx := context.Background()

	// Initialize Echo
	e := echo.New()

	// Initialize state
	app := &types.App{
		DB:      db,
		ENV:     env,
		RDB:     rdb,
		Context: ctx,
	}

	// Load json data and mapping to database
	setup.MasterDataLoader(app)

	// Middlewares
	middleware.SetupMiddlewares(e, app.ENV)

	// Base Route
	baseHandler.SetupRoutes(e, app)
	// V1 Route
	routeV1 := v1handler.RouteV1(v1handler.RouteV1{Echo: e, App: app, Enforcer: enforcer})
	routeV1.SetupRoutes()

	// Cronjob - Translate MY to EN
	scheduler.InitializeOpenAiTranslateScheduler(app)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
