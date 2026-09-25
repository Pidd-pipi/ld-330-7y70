package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/blueship581/gbemr/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PatientRepository interface {
	Create(*model.Patient) error
	FindByID(uint) (*model.Patient, error)
	FindIncludingMergedByID(uint) (*model.Patient, error)
	CountsByID(uint) (int64, int64, error)
	Search(string, int, int) ([]model.PatientWithCounts, int64, error)
	Update(*model.Patient) error
	Merge(retainedID, mergedID uint, log *model.PatientMergeRecord) error
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
func (r *patientRepository) FindIncludingMergedByID(id uint) (*model.Patient, error) {
	var p model.Patient
	if err := r.db.Unscoped().First(&p, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("find patient including merged: %w", err)
	}
	return &p, nil
}
func (r *patientRepository) CountsByID(id uint) (int64, int64, error) {
	var recordCount, prescriptionCount int64
	if err := r.db.Model(&model.MedicalRecord{}).Where("patient_id = ?", id).Count(&recordCount).Error; err != nil {
		return 0, 0, fmt.Errorf("count patient records: %w", err)
	}
	if err := r.db.Model(&model.Prescription{}).Where("patient_id = ?", id).Count(&prescriptionCount).Error; err != nil {
		return 0, 0, fmt.Errorf("count patient prescriptions: %w", err)
	}
	return recordCount, prescriptionCount, nil
}
func (r *patientRepository) Search(q string, page, size int) ([]model.PatientWithCounts, int64, error) {
	var ps []model.PatientWithCounts
	applyFilters := func(db *gorm.DB) *gorm.DB {
		if q = strings.TrimSpace(q); q != "" {
			like := "%" + q + "%"
			db = db.Where("name LIKE ? OR id_card LIKE ? OR phone LIKE ?", like, like, like)
		}
		return db
	}

	var total int64
	if err := applyFilters(r.db.Model(&model.Patient{})).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count patients: %w", err)
	}
	if err := applyFilters(r.db.Model(&model.Patient{})).
		Select("patients.*, " +
			"(SELECT COUNT(*) FROM medical_records WHERE medical_records.patient_id = patients.id) AS record_count, " +
			"(SELECT COUNT(*) FROM prescriptions WHERE prescriptions.patient_id = patients.id) AS prescription_count").
		Order("patients.created_at desc").
		Offset((page - 1) * size).
		Limit(size).
		Find(&ps).Error; err != nil {
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
func lockPatient(tx *gorm.DB, p *model.Patient, id uint) error {
	query := tx.Unscoped().Model(&model.Patient{})
	if tx.Dialector.Name() == "postgres" {
		query = query.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	if e := query.First(p, id).Error; e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			return fmt.Errorf("patient %d: %w", id, ErrNotFound)
		}
		return fmt.Errorf("lock patient %d: %w", id, e)
	}
	if p.DeletedAt.Valid {
		return fmt.Errorf("patient %d already merged: %w", id, ErrConflict)
	}
	return nil
}
func (r *patientRepository) Merge(retainedID, mergedID uint, mergeLog *model.PatientMergeRecord) (err error) {
	tx := r.db.Begin(&sql.TxOptions{Isolation: sql.LevelSerializable})
	if tx.Error != nil {
		return fmt.Errorf("begin patient merge transaction: %w", tx.Error)
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	var retained model.Patient
	if e := lockPatient(tx, &retained, retainedID); e != nil {
		return e
	}
	var merged model.Patient
	if e := lockPatient(tx, &merged, mergedID); e != nil {
		return e
	}

	var sourceRecordCount, sourcePrescriptionCount int64
	var retainedRecordCount, retainedPrescriptionCount int64
	if e := tx.Model(&model.MedicalRecord{}).Where("patient_id = ?", merged.ID).Count(&sourceRecordCount).Error; e != nil {
		return fmt.Errorf("count source records: %w", e)
	}
	if e := tx.Model(&model.Prescription{}).Where("patient_id = ?", merged.ID).Count(&sourcePrescriptionCount).Error; e != nil {
		return fmt.Errorf("count source prescriptions: %w", e)
	}
	if e := tx.Model(&model.MedicalRecord{}).Where("patient_id = ?", retained.ID).Count(&retainedRecordCount).Error; e != nil {
		return fmt.Errorf("count retained records: %w", e)
	}
	if e := tx.Model(&model.Prescription{}).Where("patient_id = ?", retained.ID).Count(&retainedPrescriptionCount).Error; e != nil {
		return fmt.Errorf("count retained prescriptions: %w", e)
	}

	result := tx.Model(&model.MedicalRecord{}).Where("patient_id = ?", merged.ID).Update("patient_id", retained.ID)
	if result.Error != nil {
		return fmt.Errorf("move medical records: %w", result.Error)
	}
	if result.RowsAffected != sourceRecordCount {
		return fmt.Errorf("move medical records verification: %w", ErrConflict)
	}
	result = tx.Model(&model.Prescription{}).Where("patient_id = ?", merged.ID).Update("patient_id", retained.ID)
	if result.Error != nil {
		return fmt.Errorf("move prescriptions: %w", result.Error)
	}
	if result.RowsAffected != sourcePrescriptionCount {
		return fmt.Errorf("move prescriptions verification: %w", ErrConflict)
	}

	var remainingRecords, remainingPrescriptions int64
	if e := tx.Model(&model.MedicalRecord{}).Where("patient_id = ?", merged.ID).Count(&remainingRecords).Error; e != nil {
		return fmt.Errorf("verify moved records: %w", e)
	}
	if e := tx.Model(&model.Prescription{}).Where("patient_id = ?", merged.ID).Count(&remainingPrescriptions).Error; e != nil {
		return fmt.Errorf("verify moved prescriptions: %w", e)
	}
	if remainingRecords != 0 || remainingPrescriptions != 0 {
		return fmt.Errorf("patient merge incomplete: %w", ErrConflict)
	}

	var finalRecordCount, finalPrescriptionCount int64
	if e := tx.Model(&model.MedicalRecord{}).Where("patient_id = ?", retained.ID).Count(&finalRecordCount).Error; e != nil {
		return fmt.Errorf("verify retained records: %w", e)
	}
	if e := tx.Model(&model.Prescription{}).Where("patient_id = ?", retained.ID).Count(&finalPrescriptionCount).Error; e != nil {
		return fmt.Errorf("verify retained prescriptions: %w", e)
	}
	if finalRecordCount != retainedRecordCount+sourceRecordCount ||
		finalPrescriptionCount != retainedPrescriptionCount+sourcePrescriptionCount {
		return fmt.Errorf("patient merge totals incomplete: %w", ErrConflict)
	}

	mergeLog.RetainedPatientID = retained.ID
	mergeLog.RetainedRecordNo = retained.RecordNo
	mergeLog.MergedPatientID = merged.ID
	mergeLog.MergedRecordNo = merged.RecordNo
	mergeLog.MedicalRecordCount = sourceRecordCount
	mergeLog.PrescriptionCount = sourcePrescriptionCount
	mergeLog.RetainedRecordCount = finalRecordCount
	mergeLog.RetainedPrescriptionCount = finalPrescriptionCount
	if e := tx.Create(mergeLog).Error; e != nil {
		return fmt.Errorf("create patient merge log: %w", e)
	}

	if e := tx.Model(&merged).Update("merged_into_id", retained.ID).Error; e != nil {
		return fmt.Errorf("mark patient merge target: %w", e)
	}
	if e := tx.Delete(&merged).Error; e != nil {
		return fmt.Errorf("hide merged patient: %w", e)
	}
	if e := tx.Commit().Error; e != nil {
		return fmt.Errorf("commit patient merge: %w", e)
	}
	committed = true
	return nil
}
