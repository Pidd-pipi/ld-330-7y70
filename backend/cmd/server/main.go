package main

import (
	"context"
	"fmt"
	"github.com/blueship581/gbemr/internal/config"
	"github.com/blueship581/gbemr/internal/handler"
	"github.com/blueship581/gbemr/internal/model"
	"github.com/blueship581/gbemr/internal/repository"
	"github.com/blueship581/gbemr/internal/router"
	"github.com/blueship581/gbemr/internal/service"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, e := config.Load()
	if e != nil {
		panic(fmt.Errorf("load config: %w", e))
	}
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	db, e := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{Logger: logger.Default.LogMode(logger.Warn)})
	if e != nil {
		panic(fmt.Errorf("open database: %w", e))
	}
	if e = migrateAndSeed(db); e != nil {
		panic(fmt.Errorf("migrate database: %w", e))
	}
	patientRepo := repository.NewPatientRepository(db)
	patientMergeRepo := repository.NewPatientMergeRepository(db)
	userRepo := repository.NewUserRepository(db)
	recordRepo := repository.NewRecordRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	presRepo := repository.NewPrescriptionRepository(db)
	departmentRepo := repository.NewDepartmentRepository(db)
	catalogRepo := repository.NewCatalogRepository(db)
	authSvc := service.NewAuthService(userRepo, cfg.JWTSecret, log)
	patientSvc := service.NewPatientService(patientRepo, patientMergeRepo, log)
	recordSvc := service.NewRecordService(recordRepo, patientRepo, departmentRepo, log)
	orderSvc := service.NewOrderService(orderRepo, recordRepo, log)
	presSvc := service.NewPrescriptionService(presRepo, recordRepo, log)
	deptSvc := service.NewDepartmentService(departmentRepo, log)
	catalogSvc := service.NewCatalogService(catalogRepo, log)
	reportSvc := service.NewReportService(db)
	r := router.New(authSvc, router.Handlers{Auth: handler.NewAuthHandler(authSvc), Patient: handler.NewPatientHandler(patientSvc, catalogSvc), Record: handler.NewRecordHandler(recordSvc, catalogSvc), Order: handler.NewOrderHandler(orderSvc), Prescription: handler.NewPrescriptionHandler(presSvc), Department: handler.NewDepartmentHandler(deptSvc), Catalog: handler.NewCatalogHandler(catalogSvc), Report: handler.NewReportHandler(reportSvc)})
	srv := &http.Server{Addr: ":" + cfg.Port, Handler: r, ReadHeaderTimeout: 5 * time.Second}
	go func() {
		log.Info("server started", "port", cfg.Port)
		if e := srv.ListenAndServe(); e != nil && e != http.ErrServerClosed {
			log.Error("server failed", "error", e)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if e := srv.Shutdown(ctx); e != nil {
		log.Error("shutdown failed", "error", e)
	}
}
func migrateAndSeed(db *gorm.DB) error {
	if e := db.AutoMigrate(&model.Department{}, &model.User{}, &model.Patient{}, &model.PatientMergeRecord{}, &model.MedicalRecord{}, &model.RecordChangeRequest{}, &model.MedicalOrder{}, &model.Prescription{}, &model.PrescriptionItem{}, &model.Drug{}, &model.DiagnosisCode{}, &model.RecordTemplate{}, &model.AuditLog{}); e != nil {
		return e
	}
	var count int64
	db.Model(&model.Department{}).Count(&count)
	if count > 0 {
		return nil
	}
	dept := model.Department{Code: "IM", Name: "内科", Description: "内科门诊与住院诊疗"}
	if e := db.Create(&dept).Error; e != nil {
		return e
	}
	accountSeeds := []struct {
		user     model.User
		password string
	}{
		{user: model.User{Username: "admin", Name: "系统管理员", Role: "admin", Active: true}, password: "admin123"},
		{user: model.User{Username: "doctor", Name: "张医生", Role: "doctor", DepartmentID: &dept.ID, Active: true}, password: "doctor123"},
		{user: model.User{Username: "nurse", Name: "李护士", Role: "nurse", DepartmentID: &dept.ID, Active: true}, password: "nurse123"},
	}
	users := make([]model.User, 0, len(accountSeeds))
	for _, seed := range accountSeeds {
		hash, e := bcrypt.GenerateFromPassword([]byte(seed.password), bcrypt.DefaultCost)
		if e != nil {
			return fmt.Errorf("hash seed password for %s: %w", seed.user.Username, e)
		}
		seed.user.PasswordHash = string(hash)
		users = append(users, seed.user)
	}
	if e := db.Create(&users).Error; e != nil {
		return e
	}
	return db.Create(&[]model.RecordTemplate{{Name: "门诊初诊模板", RecordType: "outpatient", Content: "主诉：\n现病史：\n诊断：\n治疗方案："}, {Name: "住院病历模板", RecordType: "inpatient", Content: "入院记录：\n体格检查：\n辅助检查：\n诊断："}}).Error
}
