package handler

import (
	"github.com/blueship581/gbemr/internal/service"
	"github.com/gin-gonic/gin"
)

type ReportHandler struct{ s *service.ReportService }

func NewReportHandler(s *service.ReportService) *ReportHandler { return &ReportHandler{s} }
func (h *ReportHandler) Workload(c *gin.Context) {
	v, e := h.s.DepartmentWorkload()
	if e != nil {
		Fail(c, e)
		return
	}
	OK(c, v)
}
func (h *ReportHandler) Spectrum(c *gin.Context) {
	v, e := h.s.DiseaseSpectrum()
	if e != nil {
		Fail(c, e)
		return
	}
	OK(c, v)
}
