package model

import "time"

type Patient struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	RecordNo       string    `gorm:"uniqueIndex;size:32;not null" json:"record_no"`
	Name           string    `gorm:"size:64;not null;index" json:"name"`
	Gender         string    `gorm:"size:10;not null" json:"gender"`
	Age            int       `json:"age"`
	IDCard         string    `gorm:"uniqueIndex;size:32;not null" json:"id_card"`
	Phone          string    `gorm:"size:32;index" json:"phone"`
	Allergies      string    `gorm:"type:text" json:"allergies"`
	MedicalHistory string    `gorm:"type:text" json:"medical_history"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
