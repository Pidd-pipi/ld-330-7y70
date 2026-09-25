package service

import (
	"fmt"
	"gorm.io/gorm"
)

type ReportService struct{ db *gorm.DB }

func NewReportService(db *gorm.DB) *ReportService { return &ReportService{db} }
func (s *ReportService) DepartmentWorkload() ([]map[string]interface{}, error) {
	var rows []map[string]interface{}
	e := s.db.Table("medical_records").Select("departments.name AS department, COUNT(medical_records.id) AS count").Joins("LEFT JOIN departments ON departments.id=medical_records.department_id").Group("departments.name").Order("count desc").Scan(&rows).Error
	if e != nil {
		return nil, fmt.Errorf("department workload: %w", e)
	}
	return rows, nil
}
func (s *ReportService) DiseaseSpectrum() ([]map[string]interface{}, error) {
	var rows []map[string]interface{}
	e := s.db.Table("medical_records").Select("diagnosis, COUNT(id) AS count").Where("diagnosis <> ''").Group("diagnosis").Order("count desc").Limit(10).Scan(&rows).Error
	if e != nil {
		return nil, fmt.Errorf("disease spectrum: %w", e)
	}
	return rows, nil
}
