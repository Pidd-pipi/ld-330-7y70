package service

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/blueship581/gbemr/internal/dto"
	"github.com/blueship581/gbemr/internal/model"
	"github.com/blueship581/gbemr/internal/repository"
)

type PatientService struct {
	repo   repository.PatientRepository
	logger *slog.Logger
}

func NewPatientService(r repository.PatientRepository, l *slog.Logger) *PatientService {
	return &PatientService{r, l}
}

func (s *PatientService) Create(in dto.PatientInput) (*model.Patient, error) {
	p := &model.Patient{RecordNo: fmt.Sprintf("EMR%s", time.Now().Format("20060102150405.000")), Name: in.Name, Gender: in.Gender, Age: in.Age, IDCard: in.IDCard, Phone: in.Phone, Allergies: in.Allergies, MedicalHistory: in.MedicalHistory}
	if e := s.repo.Create(p); e != nil {
		return nil, fmt.Errorf("create patient service: %w", e)
	}
	return p, nil
}

func (s *PatientService) List(q string, page, size int) ([]model.Patient, int64, error) {
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

func (s *PatientService) MergePreview(in dto.PatientMergeRequest) (*dto.PatientMergePreview, error) {
	if in.KeepPatientID == in.MergedPatientID {
		return nil, errors.New("保留档案和被合并档案不能相同")
	}
	if in.Reason = strings.TrimSpace(in.Reason); len([]rune(in.Reason)) < 2 {
		return nil, errors.New("请填写至少 2 个字的合并原因")
	}
	keep, e := s.repo.FindForMergeByID(in.KeepPatientID)
	if e != nil {
		return nil, e
	}
	merged, e := s.repo.FindForMergeByID(in.MergedPatientID)
	if e != nil {
		return nil, e
	}
	if keep.IsMerged || merged.IsMerged {
		return nil, repository.ErrPatientAlreadyMerged
	}
	keepRecords, keepPrescriptions, e := s.repo.CountRelated(keep.ID)
	if e != nil {
		return nil, e
	}
	mergedRecords, mergedPrescriptions, e := s.repo.CountRelated(merged.ID)
	if e != nil {
		return nil, e
	}
	return &dto.PatientMergePreview{
		Keep:             toMergeProfile(keep, keepRecords, keepPrescriptions),
		Merged:           toMergeProfile(merged, mergedRecords, mergedPrescriptions),
		TotalRecordCount: keepRecords + mergedRecords,
	}, nil
}

func (s *PatientService) Merge(in dto.PatientMergeRequest, operatorID uint, operatorName string) (*model.PatientMergeLog, error) {
	in.Reason = strings.TrimSpace(in.Reason)
	if _, e := s.MergePreview(in); e != nil {
		return nil, e
	}
	return s.repo.MergePatients(in.KeepPatientID, in.MergedPatientID, operatorID, operatorName, in.Reason)
}

func (s *PatientService) MergeLogs() ([]model.PatientMergeLog, error) {
	return s.repo.ListMergeLogs()
}

func toMergeProfile(p *model.Patient, recordCount, prescriptionCount int64) dto.PatientMergeProfile {
	return dto.PatientMergeProfile{
		ID:                p.ID,
		RecordNo:          p.RecordNo,
		Name:              p.Name,
		Gender:            p.Gender,
		Age:               p.Age,
		IDCard:            p.IDCard,
		Phone:             p.Phone,
		RecordCount:       recordCount,
		PrescriptionCount: prescriptionCount,
		IsMerged:          p.IsMerged,
		MergedIntoID:      p.MergedIntoID,
	}
}
