package model

import (
	"errors"

	"gorm.io/gorm"
)

type Branch struct {
	gorm.Model
	Name     string `json:"name"`
	IsActive bool   `json:"is_active" gorm:"default:true"`
}

func (B *Branch) BeforeDelete(tx *gorm.DB) error {
	return errors.New("deletion of Branch is not allowed")
}
