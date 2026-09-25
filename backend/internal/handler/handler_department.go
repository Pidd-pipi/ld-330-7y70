package handler

import (
	"github.com/blueship581/gbemr/internal/dto"
	"github.com/blueship581/gbemr/internal/service"
	"github.com/gin-gonic/gin"
)

type DepartmentHandler struct{ s *service.DepartmentService }

func NewDepartmentHandler(s *service.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{s}
}
func (h *DepartmentHandler) Create(c *gin.Context) {
	var in dto.DepartmentInput
	if !bindJSON(c, &in) {
		return
	}
	v, e := h.s.Create(in)
	if e != nil {
		Fail(c, e)
		return
	}
	OK(c, v)
}
func (h *DepartmentHandler) List(c *gin.Context) {
	v, e := h.s.List()
	if e != nil {
		Fail(c, e)
		return
	}
	OK(c, v)
}
