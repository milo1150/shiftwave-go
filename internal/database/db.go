package database

import (
	"fmt"
	"log"
	"shiftwave-go/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDatabase() *gorm.DB {
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

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	fmt.Println("Connected to PostgreSQL using GORM!")

	// Migrate the schema
	err = db.AutoMigrate(getModels()...)
	if err != nil {
		log.Fatalf("Failed to migrate db.")
	}

	return db
}

func getModels() []interface{} {
	return []interface{}{
		&model.Review{},
		&model.Branch{},
	}
}
