package repository

import (
	"fmt"
	"github.com/blueship581/gbemr/internal/model"
	"gorm.io/gorm"
)

type CatalogRepository interface {
	CreateDrug(*model.Drug) error
	ListDrugs() ([]model.Drug, error)
	CreateDiagnosis(*model.DiagnosisCode) error
	ListDiagnoses() ([]model.DiagnosisCode, error)
	CreateTemplate(*model.RecordTemplate) error
	ListTemplates() ([]model.RecordTemplate, error)
	CreateAudit(*model.AuditLog) error
	ListAudit() ([]model.AuditLog, error)
}
type catalogRepository struct{ db *gorm.DB }

func NewCatalogRepository(db *gorm.DB) CatalogRepository    { return &catalogRepository{db} }
func (r *catalogRepository) CreateDrug(v *model.Drug) error { return r.db.Create(v).Error }
func (r *catalogRepository) ListDrugs() ([]model.Drug, error) {
	var v []model.Drug
	return v, r.db.Order("name").Find(&v).Error
}
func (r *catalogRepository) CreateDiagnosis(v *model.DiagnosisCode) error {
	return r.db.Create(v).Error
}
func (r *catalogRepository) ListDiagnoses() ([]model.DiagnosisCode, error) {
	var v []model.DiagnosisCode
	return v, r.db.Order("code").Find(&v).Error
}
func (r *catalogRepository) CreateTemplate(v *model.RecordTemplate) error {
	return r.db.Create(v).Error
}
func (r *catalogRepository) ListTemplates() ([]model.RecordTemplate, error) {
	var v []model.RecordTemplate
	return v, r.db.Order("id desc").Find(&v).Error
}
func (r *catalogRepository) CreateAudit(v *model.AuditLog) error {
	if err := r.db.Create(v).Error; err != nil {
		return fmt.Errorf("create audit log: %w", err)
	}
	return nil
}
func (r *catalogRepository) ListAudit() ([]model.AuditLog, error) {
	var v []model.AuditLog
	if err := r.db.Order("created_at desc").Limit(200).Find(&v).Error; err != nil {
		return nil, fmt.Errorf("list audit logs: %w", err)
	}
	return v, nil
}
