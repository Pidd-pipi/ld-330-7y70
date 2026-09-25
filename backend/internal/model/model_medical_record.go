package model

import "time"

type MedicalRecord struct {
	ID             uint        `gorm:"primaryKey" json:"id"`
	PatientID      uint        `gorm:"index;not null" json:"patient_id"`
	Patient        *Patient    `json:"patient,omitempty"`
	DoctorID       uint        `gorm:"index;not null" json:"doctor_id"`
	Doctor         *User       `gorm:"foreignKey:DoctorID" json:"doctor,omitempty"`
	DepartmentID   uint        `gorm:"index;not null" json:"department_id"`
	Department     *Department `json:"department,omitempty"`
	RecordType     string      `gorm:"size:20;not null" json:"record_type"`
	ChiefComplaint string      `gorm:"type:text" json:"chief_complaint"`
	PresentIllness string      `gorm:"type:text" json:"present_illness"`
	PastHistory    string      `gorm:"type:text" json:"past_history"`
	PhysicalExam   string      `gorm:"type:text" json:"physical_exam"`
	AuxiliaryExam  string      `gorm:"type:text" json:"auxiliary_exam"`
	Diagnosis      string      `gorm:"type:text;index" json:"diagnosis"`
	TreatmentPlan  string      `gorm:"type:text" json:"treatment_plan"`
	RichContent    string      `gorm:"type:text" json:"rich_content"`
	Status         string      `gorm:"size:20;not null;default:'draft';index" json:"status"`
	ReviewerID     *uint       `json:"reviewer_id"`
	ReviewedAt     *time.Time  `json:"reviewed_at"`
	ArchivedAt     *time.Time  `json:"archived_at"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}
type RecordChangeRequest struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	MedicalRecordID uint      `gorm:"index" json:"medical_record_id"`
	ApplicantID     uint      `json:"applicant_id"`
	Reason          string    `gorm:"type:text" json:"reason"`
	Status          string    `gorm:"size:20;default:'pending'" json:"status"`
	CreatedAt       time.Time `json:"created_at"`
}
