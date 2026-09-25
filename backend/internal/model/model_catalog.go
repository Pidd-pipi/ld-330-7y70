package model

import "time"

type Drug struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Name          string    `gorm:"uniqueIndex;size:128;not null" json:"name"`
	Specification string    `gorm:"size:128" json:"specification"`
	Unit          string    `gorm:"size:32" json:"unit"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
type DiagnosisCode struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Code      string    `gorm:"uniqueIndex;size:32;not null" json:"code"`
	Name      string    `gorm:"size:128;not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type RecordTemplate struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Name       string    `gorm:"uniqueIndex;size:128;not null" json:"name"`
	RecordType string    `gorm:"size:20;not null" json:"record_type"`
	Content    string    `gorm:"type:text" json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
type AuditLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	Username  string    `gorm:"size:64" json:"username"`
	Action    string    `gorm:"size:100;not null" json:"action"`
	Resource  string    `gorm:"size:100" json:"resource"`
	Detail    string    `gorm:"type:text" json:"detail"`
	CreatedAt time.Time `json:"created_at"`
}
