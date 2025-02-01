package model

import (
	"errors"
	"shiftwave-go/internal/enum"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string    `gorm:"size:255;not null;unique"`
	Password     string    `gorm:"not null"`
	ActiveStatus bool      `gorm:"default:true"`
	Role         enum.Role `gorm:"default:'user'"`
	Branches     []Branch  `gorm:"many2many:users_branches"`
	Uuid         uuid.UUID `gorm:"type:uuid;uniqueIndex"`
}

func (u *User) BeforeDelete(tx *gorm.DB) error {
	return errors.New("deletion of user records is not allowed")
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.Uuid == uuid.Nil {
		u.Uuid = uuid.New() // Generates a new UUID before inserting
	}
	return nil
}
