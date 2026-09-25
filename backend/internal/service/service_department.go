package service

import (
	"github.com/blueship581/gbemr/internal/dto"
	"github.com/blueship581/gbemr/internal/model"
	"github.com/blueship581/gbemr/internal/repository"
	"log/slog"
)

type DepartmentService struct {
	repo   repository.DepartmentRepository
	logger *slog.Logger
}

func NewDepartmentService(r repository.DepartmentRepository, l *slog.Logger) *DepartmentService {
	return &DepartmentService{r, l}
}
func (s *DepartmentService) Create(in dto.DepartmentInput) (*model.Department, error) {
	v := &model.Department{Code: in.Code, Name: in.Name, Description: in.Description}
	return v, s.repo.Create(v)
}
func (s *DepartmentService) List() ([]model.Department, error) { return s.repo.List() }
