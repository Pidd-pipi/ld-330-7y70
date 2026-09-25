package handler

import (
	"github.com/blueship581/gbemr/internal/dto"
	"github.com/blueship581/gbemr/internal/service"
	"github.com/gin-gonic/gin"
)

type PrescriptionHandler struct{ s *service.PrescriptionService }

func NewPrescriptionHandler(s *service.PrescriptionService) *PrescriptionHandler {
	return &PrescriptionHandler{s}
}
func (h *PrescriptionHandler) Create(c *gin.Context) {
	var in dto.PrescriptionInput
	if !bindJSON(c, &in) {
		return
	}
	id, _, _ := currentUser(c)
	v, e := h.s.Create(in, id)
	if e != nil {
		Fail(c, e)
		return
	}
	OK(c, v)
}
func (h *PrescriptionHandler) List(c *gin.Context) {
	id, e := idParam(c)
	if e != nil {
		Fail(c, e)
		return
	}
	v, e := h.s.List(id)
	if e != nil {
		Fail(c, e)
		return
	}
	OK(c, v)
}
func (h *PrescriptionHandler) Get(c *gin.Context) {
	id, e := idParam(c)
	if e != nil {
		Fail(c, e)
		return
	}
	v, e := h.s.Get(id)
	if e != nil {
		Fail(c, e)
		return
	}
	OK(c, v)
}
func (h *PrescriptionHandler) Status(c *gin.Context) {
	id, e := idParam(c)
	if e != nil {
		Fail(c, e)
		return
	}
	var in struct {
		Status string `json:"status" binding:"required,oneof=pending_review reviewed executed"`
	}
	if e = c.ShouldBindJSON(&in); e != nil {
		Fail(c, e)
		return
	}
	v, e := h.s.Status(id, in.Status)
	if e != nil {
		Fail(c, e)
		return
	}
	OK(c, v)
}
