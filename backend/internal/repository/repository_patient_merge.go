package repository

import (
	"fmt"

	"github.com/blueship581/gbemr/internal/model"
	"gorm.io/gorm"
)

type PatientMergeRepository interface {
	List() ([]model.PatientMergeRecord, error)
}

type patientMergeRepository struct{ db *gorm.DB }

func NewPatientMergeRepository(db *gorm.DB) PatientMergeRepository {
	return &patientMergeRepository{db}
}

func (r *patientMergeRepository) List() ([]model.PatientMergeRecord, error) {
	var rows []model.PatientMergeRecord
	if err := r.db.Order("created_at desc").Limit(200).Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("list patient merge logs: %w", err)
	}
	return rows, nil
}
