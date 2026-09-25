package repository

import (
	"errors"
	"fmt"
	"github.com/blueship581/gbemr/internal/model"
	"gorm.io/gorm"
	"strings"
)

type PatientRepository interface {
	Create(*model.Patient) error
	FindByID(uint) (*model.Patient, error)
	Search(string, int, int) ([]model.Patient, int64, error)
	Update(*model.Patient) error
}
type patientRepository struct{ db *gorm.DB }

func NewPatientRepository(db *gorm.DB) PatientRepository { return &patientRepository{db} }
func (r *patientRepository) Create(p *model.Patient) error {
	if err := r.db.Create(p).Error; err != nil {
		return fmt.Errorf("create patient: %w", err)
	}
	return nil
}
func (r *patientRepository) FindByID(id uint) (*model.Patient, error) {
	var p model.Patient
	if err := r.db.First(&p, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("find patient: %w", err)
	}
	return &p, nil
}
func (r *patientRepository) Search(q string, page, size int) ([]model.Patient, int64, error) {
	var ps []model.Patient
	db := r.db.Model(&model.Patient{})
	if q = strings.TrimSpace(q); q != "" {
		like := "%" + q + "%"
		db = db.Where("name LIKE ? OR id_card LIKE ? OR phone LIKE ?", like, like, like)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count patients: %w", err)
	}
	if err := db.Order("created_at desc").Offset((page - 1) * size).Limit(size).Find(&ps).Error; err != nil {
		return nil, 0, fmt.Errorf("search patients: %w", err)
	}
	return ps, total, nil
}
func (r *patientRepository) Update(p *model.Patient) error {
	if err := r.db.Save(p).Error; err != nil {
		return fmt.Errorf("update patient: %w", err)
	}
	return nil
}
