package service

import (
	"fmt"
	"github.com/blueship581/gbemr/internal/dto"
	"github.com/blueship581/gbemr/internal/model"
	"github.com/blueship581/gbemr/internal/repository"
	"log/slog"
	"time"
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
