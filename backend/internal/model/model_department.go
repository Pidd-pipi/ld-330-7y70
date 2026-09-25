package model

import "time"

type Department struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Code        string    `gorm:"uniqueIndex;size:32;not null" json:"code"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Description string    `gorm:"size:500" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
