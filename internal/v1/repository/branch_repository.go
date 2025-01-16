package repository

import (
	"fmt"
	"shiftwave-go/internal/model"
	v1types "shiftwave-go/internal/v1/types"

	"gorm.io/gorm"
)

func CreateBranch(db *gorm.DB, branchName string) error {
	query := db.Create(&model.Branch{Name: branchName})
	return query.Error
}

func GetBranches(db *gorm.DB) (*[]model.Branch, error) {
	branches := &[]model.Branch{}

	if err := db.Find(branches).Error; err != nil {
		return nil, err
	}

	return branches, nil
}

func UpdateBranch(db *gorm.DB, id int, payload *v1types.UpdateBranchPayload) error {
	result := db.Model(&model.Branch{}).
		Where("id = ?", id).
		Updates(&model.Branch{IsActive: payload.IsActive, Name: payload.BranchName})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no branch found with id %v", id)
	}

	return nil
}
