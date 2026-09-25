package handler

import (
	"errors"
	"github.com/blueship581/gbemr/internal/dto"
	"github.com/blueship581/gbemr/internal/service"
	"github.com/gin-gonic/gin"
)

type RecordHandler struct {
	s       *service.RecordService
	catalog *service.CatalogService
}

func NewRecordHandler(s *service.RecordService, c *service.CatalogService) *RecordHandler {
	return &RecordHandler{s, c}
}
func (h *RecordHandler) Create(c *gin.Context) {
	var in dto.RecordInput
	if !bindJSON(c, &in) {
		return
	}
	id, n, _ := currentUser(c)
	v, e := h.s.Create(in, id)
	if e != nil {
		Fail(c, e)
		return
	}
	h.catalog.Audit(id, n, "create", "medical_record", v.Diagnosis)
	OK(c, v)
}
func (h *RecordHandler) List(c *gin.Context) {
	p, s := page(c)
	f := map[string]string{}
	for _, k := range []string{"patient_id", "department_id", "doctor_id", "keyword", "start_date", "end_date"} {
		f[k] = c.Query(k)
	}
	v, total, e := h.s.List(f, p, s)
	if e != nil {
		Fail(c, e)
		return
	}
	OK(c, gin.H{"items": v, "total": total, "page": p, "page_size": s})
}
func (h *RecordHandler) Get(c *gin.Context) {
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
func (h *RecordHandler) Review(c *gin.Context) {
	id, e := idParam(c)
	if e != nil {
		Fail(c, e)
		return
	}
	_, _, role := currentUser(c)
	if role != "doctor" && role != "admin" {
		Fail(c, errors.New("forbidden"))
		return
	}
	var in struct {
		Archive bool `json:"archive"`
	}
	if e = c.ShouldBindJSON(&in); e != nil {
		Fail(c, e)
		return
	}
	uid, _, _ := currentUser(c)
	v, e := h.s.Review(id, uid, in.Archive)
	if e != nil {
		Fail(c, e)
		return
	}
	OK(c, v)
}
func (h *RecordHandler) ChangeRequest(c *gin.Context) {
	id, e := idParam(c)
	if e != nil {
		Fail(c, e)
		return
	}
	var in struct {
		Reason string `json:"reason" binding:"required"`
	}
	if e = c.ShouldBindJSON(&in); e != nil {
		Fail(c, e)
		return
	}
	uid, _, _ := currentUser(c)
	if e = h.s.RequestChange(id, uid, in.Reason); e != nil {
		Fail(c, e)
		return
	}
	OK(c, gin.H{"status": "pending"})
}
