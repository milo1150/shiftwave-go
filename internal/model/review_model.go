package model

import (
	"errors"
	"log"
	"shiftwave-go/internal/connection"
	"shiftwave-go/internal/enum"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Review struct {
	gorm.Model
	Remark   string    `json:"remark"`
	Score    uint      `json:"score"`
	Lang     enum.Lang `json:"lang" gorm:"default='TH'"`
	RemarkEn string    `json:"remark_en"`

	// Depreated: Use BranchUUID instead
	BranchID uint

	// Foreign key field, Use when WHERE branch_uuid = ?
	BranchUUID uuid.UUID `gorm:"type:uuid"`

	// Explicit Foreign Key, One-to-One relationship, Use when Preload()
	Branch Branch `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:BranchUUID;references:Uuid"`
}

func (r *Review) BeforeCreate(tx *gorm.DB) (err error) {
	if !r.Lang.IsValid() {
		return errors.New("BeforeCreate Review: invalid Lang value")
	}

	return nil
}

func broadcastMessage(channel chan string, msg string) {
	select {
	// Attempt to send message without blocking when buffer is full
	case channel <- msg:
	default:
		log.Printf("Channel %v is full \n", channel)
	}
}

// Channel - Check if has any connections before send message for prevent spam x message for first incoming connection.
func (r *Review) AfterCreate(tx *gorm.DB) error {
	_, sseIsEmpty := connection.CheckActiveSseChannel()
	if !sseIsEmpty {
		broadcastMessage(connection.ReviewChannelSse, "AfterCreate Review")
	}

	_, wsIsEmpty := connection.CheckActiveWsChannel()
	if !wsIsEmpty {
		broadcastMessage(connection.ReviewChannelWs, "AfterCreate Review")
	}

	return nil
}

// Channel - Check if has any connections before send message for prevent spam x message for first incoming connection.
func (r *Review) AfterUpdate(tx *gorm.DB) error {
	_, sseIsEmpty := connection.CheckActiveSseChannel()
	if !sseIsEmpty {
		broadcastMessage(connection.ReviewChannelSse, "AfterCreate Review")
	}

	_, wsIsEmpty := connection.CheckActiveWsChannel()
	if !wsIsEmpty {
		broadcastMessage(connection.ReviewChannelWs, "AfterCreate Review")
	}

	return nil
}
