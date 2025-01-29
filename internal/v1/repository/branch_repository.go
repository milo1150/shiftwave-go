package repository

import (
	"errors"
	"fmt"
	"shiftwave-go/internal/model"
	v1types "shiftwave-go/internal/v1/types"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateBranch(db *gorm.DB, branchName string) error {
	query := db.Create(&model.Branch{Name: branchName})
	return query.Error
}

func GetBranches(db *gorm.DB) (*[]model.Branch, error) {
	branches := &[]model.Branch{}

	query := db.Order("id DESC").Find(branches)

	if err := query.Error; err != nil {
		return nil, err
	}

	return branches, nil
}

func UpdateBranch(db *gorm.DB, uuid uuid.UUID, payload *v1types.UpdateBranchPayload) error {
	result := db.Model(&model.Branch{}).
		Where("uuid = ?", uuid).
		Update("is_active", payload.IsActive).           // non-zero value
		Updates(&model.Branch{Name: payload.BranchName}) // non-zero value

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no branch found with uuid %v", uuid)
	}

	return nil
}

func FindBranchByUUID(db *gorm.DB, uuid uuid.UUID) (*model.Branch, error) {
	branch := &model.Branch{}

	query := db.Where("uuid = ?", uuid).First(branch)
	if err := query.Error; err != nil {
		return nil, err
	}

	return branch, nil
}

func FindBranchesByUUIDs(db *gorm.DB, uuids []uuid.UUID) (*[]model.Branch, error) {
	if len(uuids) == 0 {
		return nil, errors.New("uuids should not be empty")
	}

	branches := &[]model.Branch{}
	if err := db.Where("uuid IN (?)", uuids).Find(branches).Error; err != nil {
		return nil, err
	}

	if len(uuids) != len(*branches) {
		return nil, fmt.Errorf("some branch not found")
	}

	return branches, nil
}
