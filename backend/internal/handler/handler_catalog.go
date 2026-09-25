package handler

import (
	"github.com/blueship581/gbemr/internal/dto"
	"github.com/blueship581/gbemr/internal/service"
	"github.com/gin-gonic/gin"
)

type CatalogHandler struct{ s *service.CatalogService }

func NewCatalogHandler(s *service.CatalogService) *CatalogHandler { return &CatalogHandler{s} }
func (h *CatalogHandler) CreateDrug(c *gin.Context) {
	var in dto.CatalogInput
	if !bindJSON(c, &in) {
		return
	}
	v, e := h.s.Drug(in)
	if e != nil {
		Fail(c, e)
		return
	}
	OK(c, v)
}
func (h *CatalogHandler) Drugs(c *gin.Context) {
	v, e := h.s.Drugs()
	if e != nil {
		Fail(c, e)
		return
	}
	OK(c, v)
}
func (h *CatalogHandler) CreateDiagnosis(c *gin.Context) {
	var in dto.CatalogInput
	if !bindJSON(c, &in) {
		return
	}
	v, e := h.s.Diagnosis(in)
	if e != nil {
		Fail(c, e)
		return
	}
	OK(c, v)
}
func (h *CatalogHandler) Diagnoses(c *gin.Context) {
	v, e := h.s.Diagnoses()
	if e != nil {
		Fail(c, e)
		return
	}
	OK(c, v)
}
func (h *CatalogHandler) CreateTemplate(c *gin.Context) {
	var in dto.CatalogInput
	if !bindJSON(c, &in) {
		return
	}
	v, e := h.s.Template(in)
	if e != nil {
		Fail(c, e)
		return
	}
	OK(c, v)
}
func (h *CatalogHandler) Templates(c *gin.Context) {
	v, e := h.s.Templates()
	if e != nil {
		Fail(c, e)
		return
	}
	OK(c, v)
}
func (h *CatalogHandler) Audits(c *gin.Context) {
	v, e := h.s.Audits()
	if e != nil {
		Fail(c, e)
		return
	}
	OK(c, v)
}
