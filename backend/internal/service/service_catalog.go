package service

import (
	"github.com/blueship581/gbemr/internal/dto"
	"github.com/blueship581/gbemr/internal/model"
	"github.com/blueship581/gbemr/internal/repository"
	"log/slog"
)

type CatalogService struct {
	repo   repository.CatalogRepository
	logger *slog.Logger
}

func NewCatalogService(r repository.CatalogRepository, l *slog.Logger) *CatalogService {
	return &CatalogService{r, l}
}
func (s *CatalogService) Drug(in dto.CatalogInput) (*model.Drug, error) {
	v := &model.Drug{Name: in.Name, Specification: in.Specification, Unit: in.Unit}
	return v, s.repo.CreateDrug(v)
}
func (s *CatalogService) Drugs() ([]model.Drug, error) { return s.repo.ListDrugs() }
func (s *CatalogService) Diagnosis(in dto.CatalogInput) (*model.DiagnosisCode, error) {
	v := &model.DiagnosisCode{Code: in.Code, Name: in.Name}
	return v, s.repo.CreateDiagnosis(v)
}
func (s *CatalogService) Diagnoses() ([]model.DiagnosisCode, error) { return s.repo.ListDiagnoses() }
func (s *CatalogService) Template(in dto.CatalogInput) (*model.RecordTemplate, error) {
	v := &model.RecordTemplate{Name: in.Name, RecordType: in.RecordType, Content: in.Content}
	return v, s.repo.CreateTemplate(v)
}
func (s *CatalogService) Templates() ([]model.RecordTemplate, error) { return s.repo.ListTemplates() }
func (s *CatalogService) Audit(user uint, username, action, resource, detail string) {
	if e := s.repo.CreateAudit(&model.AuditLog{UserID: user, Username: username, Action: action, Resource: resource, Detail: detail}); e != nil {
		s.logger.Error("audit failure", "error", e)
	}
}
func (s *CatalogService) Audits() ([]model.AuditLog, error) { return s.repo.ListAudit() }
