package service

import (
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/blueship581/gbemr/internal/dto"
	"github.com/blueship581/gbemr/internal/model"
	"gorm.io/gorm"
)

type patientRepoFake struct {
	created       *model.Patient
	patients      map[uint]*model.Patient
	records       map[uint]int64
	prescriptions map[uint]int64
	mergeError    error
	mergeCalls    int
}

func (f *patientRepoFake) Create(p *model.Patient) error {
	p.ID = 1
	f.created = p
	return nil
}
func (f *patientRepoFake) FindByID(id uint) (*model.Patient, error) {
	if f.patients != nil {
		p, ok := f.patients[id]
		if !ok || p.DeletedAt.Valid {
			return nil, errors.New("not found")
		}
		return p, nil
	}
	return f.created, nil
}
func (f *patientRepoFake) FindIncludingMergedByID(id uint) (*model.Patient, error) {
	if f.patients != nil {
		p, ok := f.patients[id]
		if !ok {
			return nil, errors.New("not found")
		}
		return p, nil
	}
	return f.created, nil
}
func (f *patientRepoFake) CountsByID(id uint) (int64, int64, error) {
	if f.records == nil {
		return 0, 0, nil
	}
	return f.records[id], f.prescriptions[id], nil
}
func (f *patientRepoFake) Search(string, int, int) ([]model.PatientWithCounts, int64, error) {
	return nil, 0, nil
}
func (f *patientRepoFake) Update(*model.Patient) error { return nil }
func (f *patientRepoFake) Merge(uint, uint, *model.PatientMergeRecord) error {
	f.mergeCalls++
	return f.mergeError
}

type patientMergeRepoFake struct{}

func (f *patientMergeRepoFake) List() ([]model.PatientMergeRecord, error) { return nil, nil }

func TestPatientServiceCreate(t *testing.T) {
	f := &patientRepoFake{}
	s := NewPatientService(f, &patientMergeRepoFake{}, slog.Default())
	cases := []struct {
		name string
		in   dto.PatientInput
	}{{"creates unique archive number", dto.PatientInput{Name: "陈医生", Gender: "男", Age: 40, IDCard: "110101198001010001", Phone: "13800000000", Allergies: "青霉素"}}}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			p, e := s.Create(tt.in)
			if e != nil {
				t.Fatal(e)
			}
			if p.RecordNo == "" || p.Name != tt.in.Name || f.created == nil {
				t.Fatalf("unexpected patient: %#v", p)
			}
		})
	}
}

func TestPatientServiceMergeRejectsAlreadyMerged(t *testing.T) {
	deleted := model.Patient{ID: 2, RecordNo: "EMR-GONE", Name: "重复档案", DeletedAt: gorm.DeletedAt{Time: time.Now(), Valid: true}}
	f := &patientRepoFake{
		patients: map[uint]*model.Patient{
			1: {ID: 1, RecordNo: "EMR-KEEP", Name: "保留档案"},
			2: &deleted,
		},
	}
	s := NewPatientService(f, &patientMergeRepoFake{}, slog.Default())
	_, e := s.Merge(dto.PatientMergeRequest{
		RetainedPatientID: 1,
		MergedPatientID:   2,
		Reason:            "身份证录入错误",
	}, 9, "admin")
	if !errors.Is(e, ErrPatientAlreadyMerged) {
		t.Fatalf("error = %v, want ErrPatientAlreadyMerged", e)
	}
	if f.mergeCalls != 0 {
		t.Fatalf("merge calls = %d, want 0", f.mergeCalls)
	}
}
