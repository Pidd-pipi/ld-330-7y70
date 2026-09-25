package service

import (
	"errors"
	"fmt"
	"github.com/blueship581/gbemr/internal/constants"
	"github.com/blueship581/gbemr/internal/model"
	"github.com/blueship581/gbemr/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

type AuthService struct {
	users  repository.UserRepository
	secret []byte
	logger *slog.Logger
}

func NewAuthService(u repository.UserRepository, secret string, l *slog.Logger) *AuthService {
	return &AuthService{u, []byte(secret), l}
}
func (s *AuthService) Login(username, password string) (string, *model.User, error) {
	u, e := s.users.FindByUsername(username)
	if e != nil {
		return "", nil, fmt.Errorf("find login user: %w", e)
	}
	if !u.Active || bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) != nil {
		return "", nil, errors.New("invalid credentials")
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"sub": u.ID, "username": u.Username, "role": u.Role, "exp": time.Now().Add(8 * time.Hour).Unix()})
	token, e := t.SignedString(s.secret)
	if e != nil {
		return "", nil, fmt.Errorf("sign token: %w", e)
	}
	return token, u, nil
}
func (s *AuthService) CreateUser(u *model.User, password string) error {
	if _, e := s.users.FindByUsername(u.Username); e == nil {
		return errors.New("username already exists")
	} else if !errors.Is(e, repository.ErrNotFound) {
		return e
	}
	hash, e := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if e != nil {
		return fmt.Errorf("hash password: %w", e)
	}
	u.PasswordHash = string(hash)
	return s.users.Create(u)
}
func (s *AuthService) ListUsers() ([]model.User, error) { return s.users.List() }
func (s *AuthService) Parse(tokenString string) (jwt.MapClaims, error) {
	t, e := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secret, nil
	})
	if e != nil {
		return nil, e
	}
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok || !t.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

var _ = constants.RoleAdmin
