package repository

import (
	"errors"
	"fmt"
	"github.com/blueship581/gbemr/internal/model"
	"gorm.io/gorm"
	"strings"
)

type RecordRepository interface {
	Create(*model.MedicalRecord) error
	FindByID(uint) (*model.MedicalRecord, error)
	Search(map[string]string, int, int) ([]model.MedicalRecord, int64, error)
	Update(*model.MedicalRecord) error
	CreateChangeRequest(*model.RecordChangeRequest) error
}
type recordRepository struct{ db *gorm.DB }

func NewRecordRepository(db *gorm.DB) RecordRepository { return &recordRepository{db} }
func (r *recordRepository) Create(v *model.MedicalRecord) error {
	if err := r.db.Create(v).Error; err != nil {
		return fmt.Errorf("create record: %w", err)
	}
	return nil
}
func (r *recordRepository) FindByID(id uint) (*model.MedicalRecord, error) {
	var v model.MedicalRecord
	if err := r.db.Preload("Patient").Preload("Doctor").Preload("Department").First(&v, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("find record: %w", err)
	}
	return &v, nil
}
func (r *recordRepository) Search(f map[string]string, page, size int) ([]model.MedicalRecord, int64, error) {
	var rows []model.MedicalRecord
	db := r.db.Model(&model.MedicalRecord{}).Preload("Patient").Preload("Doctor").Preload("Department")
	if x := f["patient_id"]; x != "" {
		db = db.Where("patient_id = ?", x)
	}
	if x := f["department_id"]; x != "" {
		db = db.Where("department_id = ?", x)
	}
	if x := f["doctor_id"]; x != "" {
		db = db.Where("doctor_id = ?", x)
	}
	if x := strings.TrimSpace(f["keyword"]); x != "" {
		db = db.Where("diagnosis LIKE ? OR chief_complaint LIKE ? OR rich_content LIKE ?", "%"+x+"%", "%"+x+"%", "%"+x+"%")
	}
	if x := f["start_date"]; x != "" {
		db = db.Where("created_at >= ?", x)
	}
	if x := f["end_date"]; x != "" {
		db = db.Where("created_at <= ?", x+" 23:59:59")
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count records: %w", err)
	}
	if err := db.Order("created_at desc").Offset((page - 1) * size).Limit(size).Find(&rows).Error; err != nil {
		return nil, 0, fmt.Errorf("search records: %w", err)
	}
	return rows, total, nil
}
func (r *recordRepository) Update(v *model.MedicalRecord) error {
	if err := r.db.Save(v).Error; err != nil {
		return fmt.Errorf("update record: %w", err)
	}
	return nil
}
func (r *recordRepository) CreateChangeRequest(v *model.RecordChangeRequest) error {
	if err := r.db.Create(v).Error; err != nil {
		return fmt.Errorf("create change request: %w", err)
	}
	return nil
}
