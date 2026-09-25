package repository

import (
	"errors"
	"fmt"
	"github.com/blueship581/gbemr/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByUsername(string) (*model.User, error)
	Create(*model.User) error
	List() ([]model.User, error)
}
type userRepository struct{ db *gorm.DB }

func NewUserRepository(db *gorm.DB) UserRepository { return &userRepository{db} }
func (r *userRepository) FindByUsername(username string) (*model.User, error) {
	var u model.User
	if err := r.db.Preload("Department").Where("username = ?", username).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("find user: %w", err)
	}
	return &u, nil
}
func (r *userRepository) Create(u *model.User) error {
	if err := r.db.Create(u).Error; err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}
func (r *userRepository) List() ([]model.User, error) {
	var us []model.User
	if err := r.db.Preload("Department").Find(&us).Error; err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	return us, nil
}
