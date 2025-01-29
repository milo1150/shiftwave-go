package model

import (
	"errors"
	"shiftwave-go/internal/enum"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string    `gorm:"size:255;not null;unique"`
	Password     string    `gorm:"not null"`
	ActiveStatus bool      `gorm:"default:true"`
	Role         enum.Role `gorm:"default:'user'"`
	Branches     []Branch  `gorm:"many2many:users_branches"`
}

func (u *User) BeforeDelete(tx *gorm.DB) error {
	return errors.New("deletion of user records is not allowed")
}
