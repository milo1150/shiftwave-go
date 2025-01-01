package model

import "gorm.io/gorm"

type Rating struct {
	gorm.Model
	Remark string `json:"remark"`
	Score  uint   `json:"score"`

	BranchID uint   `gorm:"not null"`                                       // Foreign key referencing the Branch model
	Branch   Branch `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` // One-to-One relationship
}
