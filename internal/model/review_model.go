package model

import (
	"errors"
	"shiftwave-go/internal/types"

	"gorm.io/gorm"
)

type Review struct {
	gorm.Model
	Remark   string     `json:"remark"`
	Score    uint       `json:"score"`
	Lang     types.Lang `json:"lang" gorm:"default='TH'"`
	RemarkEn string     `json:"remark_en"`

	BranchID uint   `gorm:"not null"`                                       // Foreign key referencing the Branch model
	Branch   Branch `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` // One-to-One relationship
}

func (r *Review) BeforeCreate(tx *gorm.DB) (err error) {
	if !r.Lang.IsValid() {
		return errors.New("BeforeCreate Review: invalid Lang value")
	}

	return nil
}
