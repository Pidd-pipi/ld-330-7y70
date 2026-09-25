package model

import "time"

type MedicalOrder struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	MedicalRecordID uint      `gorm:"index;not null" json:"medical_record_id"`
	Type            string    `gorm:"size:20;not null" json:"type"`
	Content         string    `gorm:"type:text;not null" json:"content"`
	Status          string    `gorm:"size:20;not null;default:'pending'" json:"status"`
	CreatorID       uint      `json:"creator_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
