package model

import "time"

type Prescription struct {
	ID              uint               `gorm:"primaryKey" json:"id"`
	MedicalRecordID uint               `gorm:"index;not null" json:"medical_record_id"`
	PatientID       uint               `gorm:"index;not null" json:"patient_id"`
	DoctorID        uint               `json:"doctor_id"`
	Status          string             `gorm:"size:20;not null;default:'pending_review'" json:"status"`
	Items           []PrescriptionItem `json:"items"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
}
type PrescriptionItem struct {
	ID             uint   `gorm:"primaryKey" json:"id"`
	PrescriptionID uint   `gorm:"index" json:"prescription_id"`
	DrugName       string `gorm:"size:128;not null" json:"drug_name"`
	Specification  string `gorm:"size:128" json:"specification"`
	Dosage         string `gorm:"size:64;not null" json:"dosage"`
	Frequency      string `gorm:"size:64;not null" json:"frequency"`
	Duration       string `gorm:"size:64;not null" json:"duration"`
}
