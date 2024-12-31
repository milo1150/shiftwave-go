package model

import (
	"gorm.io/gorm"
)

type Branch struct {
	gorm.Model
	Name string `json:"name"`
}
