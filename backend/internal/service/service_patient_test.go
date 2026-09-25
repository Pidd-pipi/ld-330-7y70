package service

import (
	"errors"
	"testing"

	"github.com/blueship581/gbemr/internal/dto"
	"github.com/blueship581/gbemr/internal/model"
	"log/slog"
)

type patientRepoFake struct {
	created *model.Patient
}

func (f *patientRepoFake) Create(p *model.Patient) error            { p.ID = 1; f.created = p; return nil }
func (f *patientRepoFake) FindByID(id uint) (*model.Patient, error) { return f.created, nil }
func (f *patientRepoFake) FindForMergeByID(id uint) (*model.Patient, error) {
	return f.created, nil
}
func (f *patientRepoFake) Search(string, int, int) ([]model.Patient, int64, error) {
	return nil, 0, nil
}
func (f *patientRepoFake) Update(*model.Patient) error { return nil }
func (f *patientRepoFake) CountRelated(uint) (int64, int64, error) {
	return 0, 0, nil
}
func (f *patientRepoFake) MergePatients(uint, uint, uint, string, string) (*model.PatientMergeLog, error) {
	return nil, errors.New("merge patients should use repository test")
}
func (f *patientRepoFake) ListMergeLogs() ([]model.PatientMergeLog, error) {
	return nil, nil
}

func TestPatientServiceCreate(t *testing.T) {
	f := &patientRepoFake{}
	s := NewPatientService(f, slog.Default())
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

func TestPatientServiceMergeRejectsSamePatient(t *testing.T) {
	s := NewPatientService(&patientRepoFake{}, slog.Default())
	in := dto.PatientMergeRequest{KeepPatientID: 7, MergedPatientID: 7, Reason: "重复建档"}
	if _, e := s.MergePreview(in); e == nil {
		t.Fatal("expected same patient merge to be rejected")
	}
}
