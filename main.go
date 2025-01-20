package main

import (
	"context"
	"shiftwave-go/internal/database"
	baseHandler "shiftwave-go/internal/handler"
	"shiftwave-go/internal/middleware"
	"shiftwave-go/internal/scheduler"
	"shiftwave-go/internal/setup"
	"shiftwave-go/internal/types"
	v1 "shiftwave-go/internal/v1/handler"

	_ "shiftwave-go/docs" // Import Swagger docs package

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Shiftwave API
// @version 1.0
// @description This is a sample API documentation for Echo with Swagger.
// @BasePath /v1
func main() {
	// Initialize the database
	db := database.InitDatabase()

	// Initialize Redis client
	rdb := database.NewRedisClient()

	// Load and construct env
	env := setup.EnvLoader()

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
	setup.MasterDataLoader(app.DB)

	// Middlewares
	middleware.SetupMiddlewares(e, ctx, app.ENV)

	// Swagger endpoint
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Routes
	baseHandler.SetupRoutes(e, app)
	v1.SetupRoutes(e, app)

	// Cronjob - Translate MY to EN
	scheduler.InitializeOpenAiTranslateScheduler(app)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
