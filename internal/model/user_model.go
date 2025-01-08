package model

import (
	"errors"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string `gorm:"not null"`
	Password     string `gorm:"not null"`
	ActiveStatus bool
	Branches     []Branch `gorm:"many2many:users_branches"`
}

func (u *User) BeforeDelete() error {
	return errors.New("deletion of user records is not allowed")
}
