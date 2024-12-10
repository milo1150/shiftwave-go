package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func connectDB() {
	// Define the correct PostgreSQL connection string
	dsn := "host=postgres user=postgres password=postgres dbname=mydb port=5432 sslmode=disable TimeZone=UTC"

	// Connect to the database using GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Ping the database to verify the connection
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}
	defer sqlDB.Close()

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	fmt.Println("Connected to PostgreSQL using GORM!")
	fmt.Println("RUN OK")
}

func startServer() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World")
	})
	e.Logger.Fatal(e.Start(":8080"))
}

func main() {
	connectDB()
	startServer()
}
