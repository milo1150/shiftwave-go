package main

import (
	"shiftwave-go/internal/database"
	"shiftwave-go/internal/handlers"
	"shiftwave-go/internal/types"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	// Initialize the database
	db := database.InitDatabase()

	// Initialize Echo
	e := echo.New()

	// Initialize state
	app := &types.App{DB: db}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	handlers.SetupRoutes(e, app)
	// e.GET("/", func(c echo.Context) error {
	// 	return c.String(http.StatusOK, "Hello, World")
	// })
	// e.POST("/assessment", app.createProduct)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
