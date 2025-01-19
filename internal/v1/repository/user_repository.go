package repository

import (
	"errors"
	"shiftwave-go/internal/auth"
	"shiftwave-go/internal/model"
	"shiftwave-go/internal/types"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB, payload *types.CreateUserPayload) error {
	// Check if user is exists
	existedUser, _ := FindUser(db, payload.Username)
	if existedUser != nil {
		return errors.New("user already exists")
	}

	// Encrypt password
	hashPassword, err := auth.HashedPassword(payload.Password)
	if err != nil {
		return err
	}

	// Create user
	result := db.Create(&model.User{Username: payload.Username, Password: hashPassword})
	if result.Error != nil {
		return result.Error.(*pgconn.PgError)
	}

	return nil
}

func FindUser(db *gorm.DB, username string) (*model.User, error) {
	user := &model.User{}
	if err := db.First(user, &model.User{Username: username}).Error; err != nil {
		return nil, err
	}
	return user, nil
}
