package model

import (
	"errors"

	"gorm.io/gorm"
)

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type User struct {
	gorm.Model
	Username     string `gorm:"size:255;not null;unique"`
	Password     string `gorm:"not null"`
	ActiveStatus bool   `gorm:"default:true"`
	Role         string `gorm:"default:'user'"`

	Branches []Branch `gorm:"many2many:users_branches"`
}

func (u *User) BeforeDelete(tx *gorm.DB) error {
	return errors.New("deletion of user records is not allowed")
}
