package handler

import (
	"github.com/blueship581/gbemr/internal/dto"
	"github.com/blueship581/gbemr/internal/service"
	"github.com/gin-gonic/gin"
)

type PatientHandler struct {
	s       *service.PatientService
	catalog *service.CatalogService
}

func NewPatientHandler(s *service.PatientService, c *service.CatalogService) *PatientHandler {
	return &PatientHandler{s, c}
}

func (h *PatientHandler) Create(c *gin.Context) {
	var in dto.PatientInput
	if !bindJSON(c, &in) {
		return
	}
	v, e := h.s.Create(in)
	if e != nil {
		Fail(c, e)
		return
	}
	id, n, _ := currentUser(c)
	h.catalog.Audit(id, n, "create", "patient", v.RecordNo)
	OK(c, v)
}

func (h *PatientHandler) List(c *gin.Context) {
	p, s := page(c)
	v, total, e := h.s.List(c.Query("keyword"), p, s)
	if e != nil {
		Fail(c, e)
		return
	}
	OK(c, gin.H{"items": v, "total": total, "page": p, "page_size": s})
}

func (h *PatientHandler) Get(c *gin.Context) {
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

func (h *PatientHandler) Update(c *gin.Context) {
	id, e := idParam(c)
	if e != nil {
		Fail(c, e)
		return
	}
	var in dto.PatientInput
	if !bindJSON(c, &in) {
		return
	}
	v, e := h.s.Update(id, in)
	if e != nil {
		Fail(c, e)
		return
	}
	OK(c, v)
}

func (h *PatientHandler) MergePreview(c *gin.Context) {
	var in dto.PatientMergeRequest
	if !bindJSON(c, &in) {
		return
	}
	v, e := h.s.MergePreview(in)
	if e != nil {
		Fail(c, e)
		return
	}
	OK(c, v)
}

func (h *PatientHandler) Merge(c *gin.Context) {
	var in dto.PatientMergeRequest
	if !bindJSON(c, &in) {
		return
	}
	operatorID, operatorName, _ := currentUser(c)
	v, e := h.s.Merge(in, operatorID, operatorName)
	if e != nil {
		Fail(c, e)
		return
	}
	h.catalog.Audit(operatorID, operatorName, "merge", "patient", "保留档案:"+v.KeepRecordNo+",合并档案:"+v.MergedRecordNo+",原因:"+in.Reason)
	OK(c, v)
}

func (h *PatientHandler) MergeLogs(c *gin.Context) {
	v, e := h.s.MergeLogs()
	if e != nil {
		Fail(c, e)
		return
	}
	OK(c, v)
}
