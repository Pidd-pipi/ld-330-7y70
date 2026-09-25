package service

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/blueship581/gbemr/internal/dto"
	"github.com/blueship581/gbemr/internal/model"
	"github.com/blueship581/gbemr/internal/repository"
)

var ErrPatientAlreadyMerged = errors.New("patient profile has already been merged")

type PatientService struct {
	repo   repository.PatientRepository
	merges repository.PatientMergeRepository
	logger *slog.Logger
}

func NewPatientService(r repository.PatientRepository, m repository.PatientMergeRepository, l *slog.Logger) *PatientService {
	return &PatientService{repo: r, merges: m, logger: l}
}
func (s *PatientService) Create(in dto.PatientInput) (*model.Patient, error) {
	p := &model.Patient{RecordNo: fmt.Sprintf("EMR%s", time.Now().Format("20060102150405.000")), Name: in.Name, Gender: in.Gender, Age: in.Age, IDCard: in.IDCard, Phone: in.Phone, Allergies: in.Allergies, MedicalHistory: in.MedicalHistory}
	if e := s.repo.Create(p); e != nil {
		return nil, fmt.Errorf("create patient service: %w", e)
	}
	return p, nil
}
func (s *PatientService) List(q string, page, size int) ([]model.PatientWithCounts, int64, error) {
	return s.repo.Search(q, page, size)
}
func (s *PatientService) Get(id uint) (*model.Patient, error) { return s.repo.FindByID(id) }
func (s *PatientService) Update(id uint, in dto.PatientInput) (*model.Patient, error) {
	p, e := s.repo.FindByID(id)
	if e != nil {
		return nil, e
	}
	p.Name, p.Gender, p.Age, p.IDCard, p.Phone, p.Allergies, p.MedicalHistory = in.Name, in.Gender, in.Age, in.IDCard, in.Phone, in.Allergies, in.MedicalHistory
	if e = s.repo.Update(p); e != nil {
		return nil, e
	}
	return p, nil
}

func (s *PatientService) MergePreview(in dto.PatientMergePreview) (*dto.PatientMergePreviewData, error) {
	if in.RetainedPatientID == in.MergedPatientID {
		return nil, errors.New("retained and merged patient must be different")
	}
	retained, _, _, e := s.loadMergeSummary(in.RetainedPatientID)
	if e != nil {
		return nil, e
	}
	merged, mergedRecords, mergedPrescriptions, e := s.loadMergeSummary(in.MergedPatientID)
	if e != nil {
		return nil, e
	}
	return &dto.PatientMergePreviewData{
		Retained:              retained,
		Merged:                merged,
		WillMoveRecords:       mergedRecords,
		WillMovePrescriptions: mergedPrescriptions,
	}, nil
}

func (s *PatientService) Merge(in dto.PatientMergeRequest, operatorID uint, operatorName string) (*model.PatientMergeRecord, error) {
	if in.RetainedPatientID == in.MergedPatientID {
		return nil, errors.New("retained and merged patient must be different")
	}
	retained, e := s.repo.FindIncludingMergedByID(in.RetainedPatientID)
	if e != nil {
		return nil, fmt.Errorf("load retained patient: %w", e)
	}
	if retained.DeletedAt.Valid {
		return nil, fmt.Errorf("retained patient is merged: %w", ErrPatientAlreadyMerged)
	}
	merged, e := s.repo.FindIncludingMergedByID(in.MergedPatientID)
	if e != nil {
		return nil, fmt.Errorf("load merged patient: %w", e)
	}
	if merged.DeletedAt.Valid {
		return nil, fmt.Errorf("source patient is merged: %w", ErrPatientAlreadyMerged)
	}

	logRow := &model.PatientMergeRecord{
		MergedByID:   operatorID,
		MergedByName: operatorName,
		Reason:       in.Reason,
	}
	if e = s.repo.Merge(retained.ID, merged.ID, logRow); e != nil {
		return nil, e
	}
	return logRow, nil
}

func (s *PatientService) MergeHistory() ([]model.PatientMergeRecord, error) {
	return s.merges.List()
}

func (s *PatientService) loadMergeSummary(id uint) (dto.PatientMergeSummary, int64, int64, error) {
	p, e := s.repo.FindIncludingMergedByID(id)
	if e != nil {
		return dto.PatientMergeSummary{}, 0, 0, e
	}
	recordCount, prescriptionCount, e := s.repo.CountsByID(id)
	if e != nil {
		return dto.PatientMergeSummary{}, 0, 0, e
	}
	if p.DeletedAt.Valid {
		return dto.PatientMergeSummary{}, 0, 0, fmt.Errorf("%s is already merged: %w", p.RecordNo, ErrPatientAlreadyMerged)
	}
	return dto.PatientMergeSummary{
		ID:                p.ID,
		RecordNo:          p.RecordNo,
		Name:              p.Name,
		Gender:            p.Gender,
		Age:               p.Age,
		IDCard:            p.IDCard,
		Phone:             p.Phone,
		RecordCount:       recordCount,
		PrescriptionCount: prescriptionCount,
		MergedIntoID:      p.MergedIntoID,
	}, recordCount, prescriptionCount, nil
}
