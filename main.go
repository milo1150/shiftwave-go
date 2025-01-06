package main

import (
	"context"
	"shiftwave-go/internal/database"
	"shiftwave-go/internal/handler"
	"shiftwave-go/internal/middleware"
	"shiftwave-go/internal/resources"
	"shiftwave-go/internal/types"

	"github.com/labstack/echo/v4"
)

func main() {
	// Initialize the database
	db := database.InitDatabase()

	// Initialize Redis client
	rdb := database.NewRedisClient()

	// Load and construct env
	env := resources.EnvLoader()

	// Create application context
	ctx := context.Background()

	// Initialize Echo
	e := echo.New()

	// Initialize state
	app := &types.App{DB: db, ENV: env, RDB: rdb, Context: ctx}

	// Load json data and mapping to database
	resources.MasterDataLoader(app.DB)

	// Middlewares
	middleware.SetupMiddlewares(e, ctx)

	// Routes
	handler.SetupRoutes(e, app, ctx)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
