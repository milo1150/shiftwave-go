package repository

import (
	"errors"
	"shiftwave-go/internal/auth"
	"shiftwave-go/internal/enum"
	"shiftwave-go/internal/model"
	"shiftwave-go/internal/types"

	"github.com/google/uuid"
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

func GetUserByUUID(db *gorm.DB, uuid uuid.UUID) (*model.User, error) {
	user := &model.User{}

	if err := db.Preload("Branches").Where("uuid = ?", uuid).First(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func UpdateUsers(db *gorm.DB, payloads *[]types.UpdateUserPayload) error {
	for _, payload := range *payloads {
		// Find User
		userModel := model.User{}
		findUserQuery := db.Where("uuid = ?", payload.UserUuid).First(&userModel)
		if findUserQuery.Error != nil {
			return findUserQuery.Error
		}

		// Parse role enum
		parseRole, ok := enum.ParseRole(payload.Role)
		if !ok {
			return errors.New("invalid user role")
		}

		// Update User detail
		updateUserQuery := db.Model(&userModel).
			Select("ActiveStatus", "Role").
			Updates(model.User{
				ActiveStatus: payload.ActiveStatus,
				Role:         *parseRole,
			})
		if updateUserQuery.Error != nil {
			return updateUserQuery.Error
		}

		// Find Branches
		branchesModels, err := FindBranchesByUUIDs(db, payload.Branches)
		if err != nil {
			return err
		}

		// Update Uesr Branches
		query := db.Model(&userModel).Association("Branches").Replace(branchesModels)
		if query != nil {
			return errors.New("failed up update user branchres")
		}
	}
	return nil
}
