package repository

import (
	"errors"
	"shiftwave-go/internal/auth"
	"shiftwave-go/internal/enum"
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

	// Check is branches existed
	branches, err := FindBranchesByUUIDs(db, payload.Branches)
	if err != nil {
		return err
	}

	// Validate Role enum
	role, ok := enum.ParseRole(payload.Role)
	if !ok {
		return errors.New("invalid role")
	}

	// Encrypt password
	hashPassword, err := auth.HashedPassword(payload.Password)
	if err != nil {
		return err
	}

	// Create user
	result := db.Create(&model.User{
		Username: payload.Username,
		Password: hashPassword,
		Branches: *branches,
		Role:     *role,
	})
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

func GetAllUsers(db *gorm.DB) (*[]model.User, error) {
	users := &[]model.User{}

	if err := db.Preload("Branches").Order("id DESC").Find(users).Error; err != nil {
		return nil, err
	}

	return users, nil
}
