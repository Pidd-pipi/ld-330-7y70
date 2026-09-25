package repository

import (
	"fmt"
	"github.com/blueship581/gbemr/internal/model"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(*model.MedicalOrder) error
	ListByRecord(uint) ([]model.MedicalOrder, error)
	FindByID(uint) (*model.MedicalOrder, error)
	Update(*model.MedicalOrder) error
}
type orderRepository struct{ db *gorm.DB }

func NewOrderRepository(db *gorm.DB) OrderRepository { return &orderRepository{db} }
func (r *orderRepository) Create(v *model.MedicalOrder) error {
	if err := r.db.Create(v).Error; err != nil {
		return fmt.Errorf("create order: %w", err)
	}
	return nil
}
func (r *orderRepository) ListByRecord(id uint) ([]model.MedicalOrder, error) {
	var vs []model.MedicalOrder
	if err := r.db.Where("medical_record_id = ?", id).Order("created_at desc").Find(&vs).Error; err != nil {
		return nil, fmt.Errorf("list orders: %w", err)
	}
	return vs, nil
}
func (r *orderRepository) FindByID(id uint) (*model.MedicalOrder, error) {
	var v model.MedicalOrder
	if err := r.db.First(&v, id).Error; err != nil {
		return nil, fmt.Errorf("find order: %w", err)
	}
	return &v, nil
}
func (r *orderRepository) Update(v *model.MedicalOrder) error {
	if err := r.db.Save(v).Error; err != nil {
		return fmt.Errorf("update order: %w", err)
	}
	return nil
}
