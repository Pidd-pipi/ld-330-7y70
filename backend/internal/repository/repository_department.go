package repository

import (
	"fmt"
	"github.com/blueship581/gbemr/internal/model"
	"gorm.io/gorm"
)

type DepartmentRepository interface {
	Create(*model.Department) error
	List() ([]model.Department, error)
	FindByID(uint) (*model.Department, error)
}
type departmentRepository struct{ db *gorm.DB }

func NewDepartmentRepository(db *gorm.DB) DepartmentRepository { return &departmentRepository{db} }
func (r *departmentRepository) Create(v *model.Department) error {
	if err := r.db.Create(v).Error; err != nil {
		return fmt.Errorf("create department: %w", err)
	}
	return nil
}
func (r *departmentRepository) List() ([]model.Department, error) {
	var vs []model.Department
	if err := r.db.Order("id").Find(&vs).Error; err != nil {
		return nil, fmt.Errorf("list departments: %w", err)
	}
	return vs, nil
}
func (r *departmentRepository) FindByID(id uint) (*model.Department, error) {
	var v model.Department
	if err := r.db.First(&v, id).Error; err != nil {
		return nil, fmt.Errorf("find department: %w", err)
	}
	return &v, nil
}
