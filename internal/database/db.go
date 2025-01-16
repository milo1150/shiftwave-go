package database

import (
	"fmt"
	"log"
	"shiftwave-go/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func migrateBranchTable(db *gorm.DB) {
	// Drop active column and use is_active column instead.
	if db.Migrator().HasColumn(&model.Branch{}, "active") {
		err := db.Migrator().DropColumn(&model.Branch{}, "active")
		if err != nil {
			log.Fatalf("Failed to drop Active column in Branch table")
		}
	}
}

func migrateReviewTable(db *gorm.DB) {
	if db.Migrator().HasColumn(&model.Review{}, "lang") {
		// Validate lang column in Reviews table
		err := db.Exec("UPDATE reviews SET lang = 'EN' WHERE lang IS NULL").Error
		if err != nil {
			log.Fatalf("Failed to update default lang column in Reviews table")
		}
	}
}

func InitDatabase() *gorm.DB {
	// Define the correct PostgreSQL connection string
	dsn := "host=postgres user=postgres password=postgres dbname=mydb port=5432 sslmode=disable TimeZone=UTC"

	// Connect to the database using GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}

	// Manual migrate branch model
	migrateReviewTable(db)
	migrateBranchTable(db)

	// Ping the database to verify the connection
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	// Migrate the schema
	err = db.AutoMigrate(getModels()...)
	if err != nil {
		log.Fatalf("Failed to migrate db.")
	}

	fmt.Println("Connected to PostgreSQL using GORM!")

	return db
}

func getModels() []interface{} {
	return []interface{}{
		&model.User{},
		&model.Review{},
		&model.Branch{},
	}
}
