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

	"github.com/labstack/echo/v4"
)

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
	app := &types.App{DB: db, ENV: env, RDB: rdb, Context: ctx}

	// Load json data and mapping to database
	setup.MasterDataLoader(app.DB)

	// Middlewares
	middleware.SetupMiddlewares(e, ctx)

	// Routes
	baseHandler.SetupRoutes(e, app)
	v1.SetupRoutes(e, app)

	// Cronjob
	scheduler.OpenAITranslateScheduler(app)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
