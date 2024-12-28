package model

import "gorm.io/gorm"

type Rating struct {
	gorm.Model
	Remark string `json:"remark"`
	Score  uint   `json:"score"`
}
