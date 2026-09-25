package handler

import (
	"github.com/blueship581/gbemr/internal/dto"
	"github.com/blueship581/gbemr/internal/model"
	"github.com/blueship581/gbemr/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct{ s *service.AuthService }

func NewAuthHandler(s *service.AuthService) *AuthHandler { return &AuthHandler{s} }
func (h *AuthHandler) Login(c *gin.Context) {
	var in dto.LoginRequest
	if !bindJSON(c, &in) {
		return
	}
	token, u, e := h.s.Login(in.Username, in.Password)
	if e != nil {
		Fail(c, e)
		return
	}
	OK(c, gin.H{"token": token, "user": u})
}
func (h *AuthHandler) Me(c *gin.Context) {
	id, n, r := currentUser(c)
	OK(c, gin.H{"id": id, "username": n, "role": r})
}
func (h *AuthHandler) CreateUser(c *gin.Context) {
	var in dto.UserInput
	if !bindJSON(c, &in) {
		return
	}
	u := &model.User{Username: in.Username, Name: in.Name, Role: in.Role, DepartmentID: in.DepartmentID, Active: true}
	if e := h.s.CreateUser(u, in.Password); e != nil {
		Fail(c, e)
		return
	}
	OK(c, u)
}
func (h *AuthHandler) ListUsers(c *gin.Context) {
	v, e := h.s.ListUsers()
	if e != nil {
		Fail(c, e)
		return
	}
	OK(c, v)
}
