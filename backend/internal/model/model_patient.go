package model

import (
	"time"

	"gorm.io/gorm"
)

type Patient struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	RecordNo       string         `gorm:"uniqueIndex;size:32;not null" json:"record_no"`
	Name           string         `gorm:"size:64;not null;index" json:"name"`
	Gender         string         `gorm:"size:10;not null" json:"gender"`
	Age            int            `json:"age"`
	IDCard         string         `gorm:"uniqueIndex;size:32;not null" json:"id_card"`
	Phone          string         `gorm:"size:32;index" json:"phone"`
	Allergies      string         `gorm:"type:text" json:"allergies"`
	MedicalHistory string         `gorm:"type:text" json:"medical_history"`
	MergedIntoID   *uint          `gorm:"index" json:"merged_into_id,omitempty"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

type PatientWithCounts struct {
	Patient           `gorm:"embedded"`
	RecordCount       int64 `gorm:"column:record_count" json:"record_count"`
	PrescriptionCount int64 `gorm:"column:prescription_count" json:"prescription_count"`
}

func (PatientWithCounts) TableName() string { return "patients" }
