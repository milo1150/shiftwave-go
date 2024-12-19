package main

import (
	middleware "shiftwave-go/internal"
	"shiftwave-go/internal/database"
	"shiftwave-go/internal/handlers"
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

	// Middlewares
	middleware.SetupMiddlewares(e)

	// Routes
	handlers.SetupRoutes(e, app)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
