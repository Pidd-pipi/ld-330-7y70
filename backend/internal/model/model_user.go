package model

import "time"

type User struct {
	ID           uint        `gorm:"primaryKey" json:"id"`
	Username     string      `gorm:"uniqueIndex;size:64;not null" json:"username"`
	PasswordHash string      `json:"-"`
	Name         string      `gorm:"size:64;not null" json:"name"`
	Role         string      `gorm:"size:20;not null;index" json:"role"`
	DepartmentID *uint       `json:"department_id"`
	Department   *Department `json:"department,omitempty"`
	Active       bool        `gorm:"default:true" json:"active"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}
