package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base contains common columns for all tables.
type Base struct {
	ID        int64          `json:"id"`
	UUID      string         `gorm:"column:uuid" json:"uuid"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"` //add soft delete in gorm
}

func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
	// UUID shouldn't be passed while creating data with apis
	// This is required for data restore only
	uuid := uuid.NewString()
	if b.UUID == "" {
		b.UUID = uuid
	}
	return
}
