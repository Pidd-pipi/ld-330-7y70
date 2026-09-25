package service

import (
	"fmt"
	"github.com/blueship581/gbemr/internal/dto"
	"github.com/blueship581/gbemr/internal/model"
	"github.com/blueship581/gbemr/internal/repository"
	"log/slog"
)

type OrderService struct {
	repo    repository.OrderRepository
	records repository.RecordRepository
	logger  *slog.Logger
}

func NewOrderService(r repository.OrderRepository, rr repository.RecordRepository, l *slog.Logger) *OrderService {
	return &OrderService{r, rr, l}
}
func (s *OrderService) Create(in dto.OrderInput, user uint) (*model.MedicalOrder, error) {
	if _, e := s.records.FindByID(in.MedicalRecordID); e != nil {
		return nil, fmt.Errorf("validate record: %w", e)
	}
	v := &model.MedicalOrder{MedicalRecordID: in.MedicalRecordID, Type: in.Type, Content: in.Content, Status: "pending", CreatorID: user}
	if e := s.repo.Create(v); e != nil {
		return nil, e
	}
	return v, nil
}
func (s *OrderService) List(record uint) ([]model.MedicalOrder, error) {
	return s.repo.ListByRecord(record)
}
func (s *OrderService) Status(id uint, status string) (*model.MedicalOrder, error) {
	v, e := s.repo.FindByID(id)
	if e != nil {
		return nil, e
	}
	v.Status = status
	if e = s.repo.Update(v); e != nil {
		return nil, e
	}
	return v, nil
}
