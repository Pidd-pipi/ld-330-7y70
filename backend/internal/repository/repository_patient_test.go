package repository

import (
	"errors"
	"testing"

	"github.com/blueship581/gbemr/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newPatientTestDB(t *testing.T, name string) *gorm.DB {
	t.Helper()
	dsn := "file:" + name + "?mode=memory&cache=shared"
	db, e := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if e != nil {
		t.Fatal(e)
	}
	if e = db.AutoMigrate(&model.Patient{}, &model.PatientMergeLog{}, &model.MedicalRecord{}, &model.Prescription{}, &model.PrescriptionItem{}); e != nil {
		t.Fatal(e)
	}
	return db
}

func TestPatientRepositorySearch(t *testing.T) {
	db := newPatientTestDB(t, "patientsearch")
	repo := NewPatientRepository(db)
	for _, p := range []model.Patient{{RecordNo: "EMR001", Name: "王小明", Gender: "男", Age: 30, IDCard: "110101199001010001", Phone: "13800000001"}, {RecordNo: "EMR002", Name: "李华", Gender: "女", Age: 28, IDCard: "110101199201010002", Phone: "13900000002"}} {
		if e := repo.Create(&p); e != nil {
			t.Fatal(e)
		}
	}
	tests := []struct {
		q    string
		want int64
	}{{"王", 1}, {"13900000002", 1}, {"", 2}}
	for _, tt := range tests {
		t.Run(tt.q, func(t *testing.T) {
			got, total, e := repo.Search(tt.q, 1, 10)
			if e != nil || total != tt.want || len(got) != int(tt.want) {
				t.Fatalf("got len=%d total=%d err=%v", len(got), total, e)
			}
			if tt.q == "" {
				for _, p := range got {
					if p.RecordCount != 0 {
						t.Fatalf("unexpected record count: %d", p.RecordCount)
					}
				}
			}
		})
	}
}

func TestPatientRepositoryMergePatients(t *testing.T) {
	db := newPatientTestDB(t, "patientmerge")
	repo := NewPatientRepository(db)
	keep := model.Patient{RecordNo: "EMR-KEEP", Name: "同一患者", Gender: "女", Age: 32, IDCard: "110101199101010001", Phone: "13800000001"}
	merged := model.Patient{RecordNo: "EMR-MERGED", Name: "同一患者", Gender: "女", Age: 32, IDCard: "110101199101010002", Phone: "13800000001"}
	if e := repo.Create(&keep); e != nil {
		t.Fatal(e)
	}
	if e := repo.Create(&merged); e != nil {
		t.Fatal(e)
	}
	records := []model.MedicalRecord{
		{PatientID: keep.ID, DoctorID: 1, DepartmentID: 1, RecordType: "outpatient", Status: "draft"},
		{PatientID: merged.ID, DoctorID: 1, DepartmentID: 1, RecordType: "outpatient", Status: "draft"},
		{PatientID: merged.ID, DoctorID: 1, DepartmentID: 1, RecordType: "inpatient", Status: "archived"},
	}
	for i := range records {
		if e := db.Create(&records[i]).Error; e != nil {
			t.Fatal(e)
		}
	}
	prescriptions := []model.Prescription{
		{PatientID: merged.ID, DoctorID: 1, MedicalRecordID: records[1].ID, Status: "pending_review"},
		{PatientID: merged.ID, DoctorID: 1, MedicalRecordID: records[2].ID, Status: "reviewed"},
	}
	for i := range prescriptions {
		if e := db.Create(&prescriptions[i]).Error; e != nil {
			t.Fatal(e)
		}
	}

	logEntry, e := repo.MergePatients(keep.ID, merged.ID, 9, "管理员", "身份证录错，合并重复档案")
	if e != nil {
		t.Fatal(e)
	}
	if logEntry.MovedRecordCount != 2 || logEntry.MovedPrescriptionCount != 2 {
		t.Fatalf("unexpected moved counts: %#v", logEntry)
	}

	var keepRecordCount, keepPrescriptionCount, mergedRecordCount, mergedPrescriptionCount int64
	db.Model(&model.MedicalRecord{}).Where("patient_id = ?", keep.ID).Count(&keepRecordCount)
	db.Model(&model.Prescription{}).Where("patient_id = ?", keep.ID).Count(&keepPrescriptionCount)
	db.Model(&model.MedicalRecord{}).Where("patient_id = ?", merged.ID).Count(&mergedRecordCount)
	db.Model(&model.Prescription{}).Where("patient_id = ?", merged.ID).Count(&mergedPrescriptionCount)
	if keepRecordCount != 3 || keepPrescriptionCount != 2 || mergedRecordCount != 0 || mergedPrescriptionCount != 0 {
		t.Fatalf("unexpected association counts records=(keep %d, merged %d), prescriptions=(keep %d, merged %d)", keepRecordCount, mergedRecordCount, keepPrescriptionCount, mergedPrescriptionCount)
	}

	hidden, e := repo.FindByID(merged.ID)
	if !errors.Is(e, ErrNotFound) || hidden != nil {
		t.Fatalf("merged patient should disappear from normal lookup, got %#v %v", hidden, e)
	}
	remaining, total, e := repo.Search("", 1, 10)
	if e != nil || total != 1 {
		t.Fatalf("merged patient should disappear from search, total=%d err=%v", total, e)
	}
	if remaining[0].ID != keep.ID || remaining[0].RecordCount != 3 {
		t.Fatalf("search should show retained patient with 3 records, got %#v", remaining[0])
	}
	source, e := repo.FindForMergeByID(merged.ID)
	if e != nil || !source.IsMerged || source.MergedIntoID == nil || *source.MergedIntoID != keep.ID {
		t.Fatalf("merged source marker incorrect: %#v %v", source, e)
	}
	if _, e = repo.MergePatients(keep.ID, merged.ID, 9, "管理员", "再次合并"); !errors.Is(e, ErrPatientAlreadyMerged) {
		t.Fatalf("expected already merged error, got %v", e)
	}
}
