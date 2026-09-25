package router

import (
	"log/slog"
	"net/http"

	"github.com/blueship581/gbemr/internal/constants"
	"github.com/blueship581/gbemr/internal/handler"
	"github.com/blueship581/gbemr/internal/middleware"
	"github.com/blueship581/gbemr/internal/service"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Auth         *handler.AuthHandler
	Patient      *handler.PatientHandler
	Record       *handler.RecordHandler
	Order        *handler.OrderHandler
	Prescription *handler.PrescriptionHandler
	Department   *handler.DepartmentHandler
	Catalog      *handler.CatalogHandler
	Report       *handler.ReportHandler
}

func New(auth *service.AuthService, h Handlers) *gin.Engine {
	r := gin.New()
	r.Use(middleware.Recovery(), middleware.RequestLogger(slog.Default()))
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"code": constants.CodeOK, "message": "ok", "data": gin.H{"status": "healthy"}})
	})
	registerAPI(r.Group("/api/v1"), auth, h)
	registerAPI(r.Group("/v1"), auth, h)
	return r
}

// /v1 is an internal compatibility alias because Nginx strips /api/ when proxy_pass ends in a slash.
func registerAPI(api *gin.RouterGroup, auth *service.AuthService, h Handlers) {
	api.POST("/auth/login", h.Auth.Login)
	secured := api.Group("")
	secured.Use(middleware.Auth(auth))
	secured.GET("/auth/me", h.Auth.Me)
	secured.GET("/patients", h.Patient.List)
	secured.GET("/patients/:id", h.Patient.Get)
	secured.POST("/patients", middleware.Roles(constants.RoleDoctor, constants.RoleNurse, constants.RoleAdmin), h.Patient.Create)
	secured.PUT("/patients/:id", middleware.Roles(constants.RoleDoctor, constants.RoleNurse, constants.RoleAdmin), h.Patient.Update)
	secured.POST("/patients/merge-preview", middleware.Roles(constants.RoleAdmin), h.Patient.MergePreview)
	secured.POST("/patients/merge", middleware.Roles(constants.RoleAdmin), h.Patient.Merge)
	secured.GET("/patient-merge-logs", middleware.Roles(constants.RoleAdmin), h.Patient.MergeLogs)
	secured.GET("/records", h.Record.List)
	secured.GET("/records/:id", h.Record.Get)
	secured.POST("/records", middleware.Roles(constants.RoleDoctor), h.Record.Create)
	secured.POST("/records/:id/review", middleware.Roles(constants.RoleDoctor, constants.RoleAdmin), h.Record.Review)
	secured.POST("/records/:id/change-requests", middleware.Roles(constants.RoleDoctor), h.Record.ChangeRequest)
	secured.POST("/orders", middleware.Roles(constants.RoleDoctor), h.Order.Create)
	secured.GET("/records/:id/orders", h.Order.List)
	secured.PUT("/orders/:id/status", middleware.Roles(constants.RoleDoctor, constants.RoleNurse), h.Order.Status)
	secured.POST("/prescriptions", middleware.Roles(constants.RoleDoctor), h.Prescription.Create)
	secured.GET("/records/:id/prescriptions", h.Prescription.List)
	secured.GET("/prescriptions/:id", h.Prescription.Get)
	secured.PUT("/prescriptions/:id/status", middleware.Roles(constants.RoleDoctor, constants.RoleNurse), h.Prescription.Status)
	secured.GET("/departments", h.Department.List)
	admin := secured.Group("/admin", middleware.Roles(constants.RoleAdmin))
	admin.POST("/departments", h.Department.Create)
	admin.POST("/users", h.Auth.CreateUser)
	admin.GET("/users", h.Auth.ListUsers)
	admin.POST("/drugs", h.Catalog.CreateDrug)
	admin.GET("/drugs", h.Catalog.Drugs)
	admin.POST("/diagnoses", h.Catalog.CreateDiagnosis)
	admin.GET("/diagnoses", h.Catalog.Diagnoses)
	admin.POST("/templates", h.Catalog.CreateTemplate)
	admin.GET("/templates", h.Catalog.Templates)
	admin.GET("/audit-logs", h.Catalog.Audits)
	secured.GET("/reports/department-workload", middleware.Roles(constants.RoleDoctor, constants.RoleAdmin), h.Report.Workload)
	secured.GET("/reports/disease-spectrum", middleware.Roles(constants.RoleDoctor, constants.RoleAdmin), h.Report.Spectrum)
}
