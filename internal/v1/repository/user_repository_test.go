package repository_test

import (
	"shiftwave-go/internal/v1/repository"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestFindUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create SQL mock: %v", err)
	}
	defer db.Close()

	// Init GORM with mock db
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to init gormDB: %v", err)
	}

	// Set up mock to return the desired result
	mockId, mockUsername := 1, "cc"
	mock.ExpectQuery(".*").
		WithArgs(mockUsername, mockId).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).
			AddRow(mockId, mockUsername))

	// Call the repository function
	user, err := repository.FindUser(gormDB, mockUsername)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Validate the result
	if user.ID != uint(mockId) || user.Username != mockUsername {
		t.Fatalf("unexpected user: %+v", user)
	}

	// Ensure all mock expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}
}
