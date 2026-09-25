package service

import (
	"fmt"
	"github.com/blueship581/gbemr/internal/dto"
	"github.com/blueship581/gbemr/internal/model"
	"github.com/blueship581/gbemr/internal/repository"
	"log/slog"
)

type PrescriptionService struct {
	repo    repository.PrescriptionRepository
	records repository.RecordRepository
	logger  *slog.Logger
}

func NewPrescriptionService(r repository.PrescriptionRepository, rr repository.RecordRepository, l *slog.Logger) *PrescriptionService {
	return &PrescriptionService{r, rr, l}
}
func (s *PrescriptionService) Create(in dto.PrescriptionInput, doctor uint) (*model.Prescription, error) {
	rec, e := s.records.FindByID(in.MedicalRecordID)
	if e != nil {
		return nil, fmt.Errorf("validate record: %w", e)
	}
	v := &model.Prescription{MedicalRecordID: rec.ID, PatientID: rec.PatientID, DoctorID: doctor, Status: "pending_review"}
	for _, x := range in.Items {
		v.Items = append(v.Items, model.PrescriptionItem{DrugName: x.DrugName, Specification: x.Specification, Dosage: x.Dosage, Frequency: x.Frequency, Duration: x.Duration})
	}
	if e = s.repo.Create(v); e != nil {
		return nil, e
	}
	return s.repo.FindByID(v.ID)
}
func (s *PrescriptionService) List(record uint) ([]model.Prescription, error) {
	return s.repo.ListByRecord(record)
}
func (s *PrescriptionService) Get(id uint) (*model.Prescription, error) { return s.repo.FindByID(id) }
func (s *PrescriptionService) Status(id uint, status string) (*model.Prescription, error) {
	v, e := s.repo.FindByID(id)
	if e != nil {
		return nil, e
	}
	v.Status = status
	if e = s.repo.Update(v); e != nil {
		return nil, e
	}
	return s.repo.FindByID(id)
}
