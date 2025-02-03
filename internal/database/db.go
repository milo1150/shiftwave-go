package database

import (
	"fmt"
	"log"
	"shiftwave-go/internal/model"

	"github.com/google/uuid"
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

	if db.Migrator().HasColumn(&model.Branch{}, "uuid") {
		branches := &[]model.Branch{}

		query := db.Model(&model.Branch{}).Where("uuid IS NULL").Find(&branches)
		if err := query.Error; err != nil {
			log.Fatalf("Failed to fetch branches with null uuid: %v", err)
		}

		if len(*branches) > 0 {
			for _, branch := range *branches {
				query := db.Model(&model.Branch{}).Where("id = ?", branch.ID).Update("uuid", uuid.New())
				if err := query.Error; err != nil {
					log.Fatalf("Failed to update branch uuid: %v", err)
				}
			}
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

	if db.Migrator().HasColumn(&model.Review{}, "branch_id") && db.Migrator().HasColumn(&model.Review{}, "branch_uuid") {
		reviews := []model.Review{}
		if err := db.Model(&model.Review{}).Where("branch_uuid IS NULL").Find(&reviews).Error; err != nil {
			log.Fatalf("Failed to count reviews: %v", err.Error())
		}

		// If all data already has branch_uuid value then do nothing.
		if len(reviews) == 0 {
			return
		}

		for _, review := range reviews {
			branch := model.Branch{}
			if err := db.Where("id = ?", review.BranchID).Find(&branch).Error; err != nil {
				log.Fatalf("Error branch not found: %v", err.Error())
			}

			if err := db.Model(&review).Update("branch_uuid", branch.Uuid).Error; err != nil {
				log.Fatalf("Error update branch_uuid: %v", err.Error())
			}
		}
	}
}

func migrateUserTable(db *gorm.DB) {
	if db.Migrator().HasColumn(&model.User{}, "uuid") {
		users := &[]model.User{}

		query := db.Where("uuid IS NULL").Find(&users)
		if err := query.Error; err != nil {
			log.Fatalf("Failed to fetch users with null uuid: %v", err)
		}

		if len(*users) > 0 {
			for _, user := range *users {
				if err := db.Model(&model.User{}).Where("id = ?", user.ID).Update("uuid", uuid.New()).Error; err != nil {
					log.Fatalf("Failed to update user uuid: %v", err.Error())
				}
			}
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
	migrateUserTable(db)

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
