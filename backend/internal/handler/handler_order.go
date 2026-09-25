package handler

import (
	"github.com/blueship581/gbemr/internal/dto"
	"github.com/blueship581/gbemr/internal/service"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct{ s *service.OrderService }

func NewOrderHandler(s *service.OrderService) *OrderHandler { return &OrderHandler{s} }
func (h *OrderHandler) Create(c *gin.Context) {
	var in dto.OrderInput
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
func (h *OrderHandler) List(c *gin.Context) {
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
func (h *OrderHandler) Status(c *gin.Context) {
	id, e := idParam(c)
	if e != nil {
		Fail(c, e)
		return
	}
	var in struct {
		Status string `json:"status" binding:"required,oneof=pending executed stopped"`
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
