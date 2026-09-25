package model

import "time"

type Patient struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	RecordNo       string     `gorm:"uniqueIndex;size:32;not null" json:"record_no"`
	Name           string     `gorm:"size:64;not null;index" json:"name"`
	Gender         string     `gorm:"size:10;not null" json:"gender"`
	Age            int        `json:"age"`
	IDCard         string     `gorm:"uniqueIndex;size:32;not null" json:"id_card"`
	Phone          string     `gorm:"size:32;index" json:"phone"`
	Allergies      string     `gorm:"type:text" json:"allergies"`
	MedicalHistory string     `gorm:"type:text" json:"medical_history"`
	IsMerged       bool       `gorm:"not null;default:false;index" json:"is_merged"`
	MergedIntoID   *uint      `gorm:"index" json:"merged_into_id,omitempty"`
	MergedAt       *time.Time `json:"merged_at,omitempty"`
	RecordCount    int64      `gorm:"->" json:"record_count,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type PatientMergeLog struct {
	ID                     uint      `gorm:"primaryKey" json:"id"`
	KeepPatientID          uint      `gorm:"index;not null" json:"keep_patient_id"`
	KeepRecordNo           string    `gorm:"size:32;not null" json:"keep_record_no"`
	KeepPatientName        string    `gorm:"size:64;not null" json:"keep_patient_name"`
	MergedPatientID        uint      `gorm:"index;not null" json:"merged_patient_id"`
	MergedRecordNo         string    `gorm:"size:32;not null" json:"merged_record_no"`
	MergedPatientName      string    `gorm:"size:64;not null" json:"merged_patient_name"`
	MovedRecordCount       int64     `gorm:"not null" json:"moved_record_count"`
	MovedPrescriptionCount int64     `gorm:"not null" json:"moved_prescription_count"`
	Reason                 string    `gorm:"type:text;not null" json:"reason"`
	OperatorID             uint      `gorm:"index;not null" json:"operator_id"`
	OperatorName           string    `gorm:"size:64;not null" json:"operator_name"`
	MergedAt               time.Time `gorm:"index" json:"merged_at"`
}
