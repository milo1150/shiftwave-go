package model

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Branch struct {
	gorm.Model
	Name     string    `json:"name"`
	IsActive bool      `json:"is_active" gorm:"default:true"`
	Uuid     uuid.UUID `json:"uuid" gorm:"type:uuid;unique"`
}

func (b *Branch) BeforeDelete(tx *gorm.DB) error {
	return errors.New("deletion of Branch is not allowed")
}

func (b *Branch) BeforeCreate(tx *gorm.DB) (err error) {
	if b.Uuid == uuid.Nil {
		b.Uuid = uuid.New() // Generates a new UUID before inserting
	}
	return nil
}
