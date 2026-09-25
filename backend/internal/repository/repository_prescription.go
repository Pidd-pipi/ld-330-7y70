package repository

import (
	"fmt"
	"github.com/blueship581/gbemr/internal/model"
	"gorm.io/gorm"
)

type PrescriptionRepository interface {
	Create(*model.Prescription) error
	FindByID(uint) (*model.Prescription, error)
	ListByRecord(uint) ([]model.Prescription, error)
	Update(*model.Prescription) error
}
type prescriptionRepository struct{ db *gorm.DB }

func NewPrescriptionRepository(db *gorm.DB) PrescriptionRepository {
	return &prescriptionRepository{db}
}
func (r *prescriptionRepository) Create(v *model.Prescription) error {
	if err := r.db.Create(v).Error; err != nil {
		return fmt.Errorf("create prescription: %w", err)
	}
	return nil
}
func (r *prescriptionRepository) FindByID(id uint) (*model.Prescription, error) {
	var v model.Prescription
	if err := r.db.Preload("Items").First(&v, id).Error; err != nil {
		return nil, fmt.Errorf("find prescription: %w", err)
	}
	return &v, nil
}
func (r *prescriptionRepository) ListByRecord(id uint) ([]model.Prescription, error) {
	var vs []model.Prescription
	if err := r.db.Preload("Items").Where("medical_record_id=?", id).Find(&vs).Error; err != nil {
		return nil, fmt.Errorf("list prescriptions: %w", err)
	}
	return vs, nil
}
func (r *prescriptionRepository) Update(v *model.Prescription) error {
	if err := r.db.Save(v).Error; err != nil {
		return fmt.Errorf("update prescription: %w", err)
	}
	return nil
}
