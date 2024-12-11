package models

import "gorm.io/gorm"

type Assessment struct {
	gorm.Model
	Remark string
	Score  uint
}
