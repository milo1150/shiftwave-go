package main

import (
	"shiftwave-go/internal/database"
	"shiftwave-go/internal/handler"
	"shiftwave-go/internal/middleware"
	"shiftwave-go/internal/types"

	"github.com/labstack/echo/v4"
)

func main() {
	// Initialize the database
	db := database.InitDatabase()

	// Initialize Echo
	e := echo.New()

	// Initialize state
	app := &types.App{DB: db}

	// Load json data and mapping to database
	database.MasterDataLoader(app.DB)

	// Middlewares
	middleware.SetupMiddlewares(e)

	// Routes
	handler.SetupRoutes(e, app)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
