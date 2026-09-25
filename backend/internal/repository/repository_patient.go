package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"strings"
	"time"

	"github.com/blueship581/gbemr/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PatientRepository interface {
	Create(*model.Patient) error
	FindByID(uint) (*model.Patient, error)
	FindForMergeByID(uint) (*model.Patient, error)
	Search(string, int, int) ([]model.Patient, int64, error)
	Update(*model.Patient) error
	CountRelated(uint) (int64, int64, error)
	MergePatients(keepID, mergedID, operatorID uint, operatorName, reason string) (*model.PatientMergeLog, error)
	ListMergeLogs() ([]model.PatientMergeLog, error)
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
	if err := r.db.Where("is_merged = ?", false).First(&p, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("find patient: %w", err)
	}
	return &p, nil
}

func (r *patientRepository) FindForMergeByID(id uint) (*model.Patient, error) {
	var p model.Patient
	if err := r.db.First(&p, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("find patient for merge: %w", err)
	}
	return &p, nil
}

func (r *patientRepository) Search(q string, page, size int) ([]model.Patient, int64, error) {
	var ps []model.Patient
	db := r.db.Model(&model.Patient{}).Where("is_merged = ?", false)
	if q = strings.TrimSpace(q); q != "" {
		like := "%" + q + "%"
		db = db.Where("name LIKE ? OR id_card LIKE ? OR phone LIKE ?", like, like, like)
	}
	var total int64
	if err := db.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count patients: %w", err)
	}
	countSubQuery := "(SELECT COUNT(*) FROM medical_records WHERE medical_records.patient_id = patients.id)"
	if err := db.Select("patients.*, " + countSubQuery + " AS record_count").Order("created_at desc").Offset((page - 1) * size).Limit(size).Find(&ps).Error; err != nil {
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

func (r *patientRepository) CountRelated(id uint) (int64, int64, error) {
	var recordCount, prescriptionCount int64
	if err := r.db.Model(&model.MedicalRecord{}).Where("patient_id = ?", id).Count(&recordCount).Error; err != nil {
		return 0, 0, fmt.Errorf("count patient records: %w", err)
	}
	if err := r.db.Model(&model.Prescription{}).Where("patient_id = ?", id).Count(&prescriptionCount).Error; err != nil {
		return 0, 0, fmt.Errorf("count patient prescriptions: %w", err)
	}
	return recordCount, prescriptionCount, nil
}

func (r *patientRepository) MergePatients(keepID, mergedID, operatorID uint, operatorName, reason string) (*model.PatientMergeLog, error) {
	var logEntry *model.PatientMergeLog
	txOptions := &sql.TxOptions{}
	if r.db.Dialector.Name() == "postgres" {
		txOptions.Isolation = sql.LevelSerializable
	}
	err := r.db.Transaction(func(tx *gorm.DB) error {
		ids := []uint{keepID, mergedID}
		if keepID > mergedID {
			ids[0], ids[1] = mergedID, keepID
		}
		query := tx.Model(&model.Patient{}).Where("id IN ?", ids)
		if tx.Dialector.Name() == "postgres" {
			query = query.Clauses(clause.Locking{Strength: "UPDATE"})
		}

		var patients []model.Patient
		if err := query.Find(&patients).Error; err != nil {
			return fmt.Errorf("lock patients for merge: %w", err)
		}
		if len(patients) != 2 {
			return ErrNotFound
		}

		var keep, merged *model.Patient
		for i := range patients {
			switch patients[i].ID {
			case keepID:
				keep = &patients[i]
			case mergedID:
				merged = &patients[i]
			}
		}
		if keep == nil || merged == nil {
			return ErrNotFound
		}
		if keep.IsMerged || merged.IsMerged {
			return ErrPatientAlreadyMerged
		}

		var recordCount, prescriptionCount, keepRecordCount, keepPrescriptionCount int64
		if err := tx.Model(&model.MedicalRecord{}).Where("patient_id = ?", mergedID).Count(&recordCount).Error; err != nil {
			return fmt.Errorf("count merged records: %w", err)
		}
		if err := tx.Model(&model.Prescription{}).Where("patient_id = ?", mergedID).Count(&prescriptionCount).Error; err != nil {
			return fmt.Errorf("count merged prescriptions: %w", err)
		}
		if err := tx.Model(&model.MedicalRecord{}).Where("patient_id = ?", keepID).Count(&keepRecordCount).Error; err != nil {
			return fmt.Errorf("count retained records: %w", err)
		}
		if err := tx.Model(&model.Prescription{}).Where("patient_id = ?", keepID).Count(&keepPrescriptionCount).Error; err != nil {
			return fmt.Errorf("count retained prescriptions: %w", err)
		}

		result := tx.Model(&model.MedicalRecord{}).Where("patient_id = ?", mergedID).Update("patient_id", keepID)
		if result.Error != nil {
			return fmt.Errorf("move medical records: %w", result.Error)
		}
		if result.RowsAffected != recordCount {
			return ErrMergeIncomplete
		}

		result = tx.Model(&model.Prescription{}).Where("patient_id = ?", mergedID).Update("patient_id", keepID)
		if result.Error != nil {
			return fmt.Errorf("move prescriptions: %w", result.Error)
		}
		if result.RowsAffected != prescriptionCount {
			return ErrMergeIncomplete
		}

		var remainingRecords, remainingPrescriptions, finalRecords, finalPrescriptions int64
		if err := tx.Model(&model.MedicalRecord{}).Where("patient_id = ?", mergedID).Count(&remainingRecords).Error; err != nil {
			return fmt.Errorf("verify remaining records: %w", err)
		}
		if err := tx.Model(&model.Prescription{}).Where("patient_id = ?", mergedID).Count(&remainingPrescriptions).Error; err != nil {
			return fmt.Errorf("verify remaining prescriptions: %w", err)
		}
		if err := tx.Model(&model.MedicalRecord{}).Where("patient_id = ?", keepID).Count(&finalRecords).Error; err != nil {
			return fmt.Errorf("verify moved records: %w", err)
		}
		if err := tx.Model(&model.Prescription{}).Where("patient_id = ?", keepID).Count(&finalPrescriptions).Error; err != nil {
			return fmt.Errorf("verify moved prescriptions: %w", err)
		}
		if remainingRecords != 0 || remainingPrescriptions != 0 || finalRecords != keepRecordCount+recordCount || finalPrescriptions != keepPrescriptionCount+prescriptionCount {
			return ErrMergeIncomplete
		}

		now := time.Now()
		merged.IsMerged = true
		merged.MergedIntoID = &keepID
		merged.MergedAt = &now
		if err := tx.Save(merged).Error; err != nil {
			return fmt.Errorf("mark merged patient: %w", err)
		}

		logEntry = &model.PatientMergeLog{
			KeepPatientID:          keep.ID,
			KeepRecordNo:           keep.RecordNo,
			KeepPatientName:        keep.Name,
			MergedPatientID:        merged.ID,
			MergedRecordNo:         merged.RecordNo,
			MergedPatientName:      merged.Name,
			MovedRecordCount:       recordCount,
			MovedPrescriptionCount: prescriptionCount,
			Reason:                 reason,
			OperatorID:             operatorID,
			OperatorName:           operatorName,
			MergedAt:               now,
		}
		if err := tx.Create(logEntry).Error; err != nil {
			return fmt.Errorf("create patient merge log: %w", err)
		}
		return nil
	}, txOptions)
	if err != nil {
		if errors.Is(err, ErrPatientAlreadyMerged) || errors.Is(err, ErrMergeIncomplete) {
			return nil, err
		}
		return nil, fmt.Errorf("merge patients: %w", err)
	}
	return logEntry, nil
}

func (r *patientRepository) ListMergeLogs() ([]model.PatientMergeLog, error) {
	var logs []model.PatientMergeLog
	if err := r.db.Order("merged_at desc").Limit(200).Find(&logs).Error; err != nil {
		return nil, fmt.Errorf("list patient merge logs: %w", err)
	}
	return logs, nil
}
