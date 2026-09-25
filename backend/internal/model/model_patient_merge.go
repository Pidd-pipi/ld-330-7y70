package model

import "time"

// PatientMergeRecord is the immutable audit trail for an administrator merge.
type PatientMergeRecord struct {
	ID                        uint      `gorm:"primaryKey" json:"id"`
	RetainedPatientID         uint      `gorm:"index;not null" json:"retained_patient_id"`
	RetainedRecordNo          string    `gorm:"size:32;not null" json:"retained_record_no"`
	MergedPatientID           uint      `gorm:"uniqueIndex;not null" json:"merged_patient_id"`
	MergedRecordNo            string    `gorm:"size:32;not null" json:"merged_record_no"`
	MergedByID                uint      `gorm:"index;not null" json:"merged_by_id"`
	MergedByName              string    `gorm:"size:64;not null" json:"merged_by_name"`
	Reason                    string    `gorm:"type:text;not null" json:"reason"`
	MedicalRecordCount        int64     `json:"medical_record_count"`
	PrescriptionCount         int64     `json:"prescription_count"`
	RetainedRecordCount       int64     `json:"retained_record_count"`
	RetainedPrescriptionCount int64     `json:"retained_prescription_count"`
	CreatedAt                 time.Time `json:"created_at"`
}
