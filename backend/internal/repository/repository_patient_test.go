package repository

import (
	"testing"

	"github.com/blueship581/gbemr/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestPatientRepositorySearch(t *testing.T) {
	db, e := gorm.Open(sqlite.Open("file:patientrepo?mode=memory&cache=shared"), &gorm.Config{})
	if e != nil {
		t.Fatal(e)
	}
	if e = db.AutoMigrate(&model.Patient{}, &model.MedicalRecord{}, &model.Prescription{}); e != nil {
		t.Fatal(e)
	}
	repo := NewPatientRepository(db)
	patients := []model.Patient{
		{RecordNo: "EMR001", Name: "王小明", Gender: "男", Age: 30, IDCard: "110101199001010001", Phone: "13800000001"},
		{RecordNo: "EMR002", Name: "李华", Gender: "女", Age: 28, IDCard: "110101199201010002", Phone: "13900000002"},
	}
	for i := range patients {
		if e := repo.Create(&patients[i]); e != nil {
			t.Fatal(e)
		}
	}
	record := model.MedicalRecord{PatientID: patients[0].ID, DoctorID: 1, DepartmentID: 1, RecordType: "outpatient", Status: "draft"}
	if e := db.Create(&record).Error; e != nil {
		t.Fatal(e)
	}
	if e := db.Create(&model.Prescription{MedicalRecordID: record.ID, PatientID: patients[0].ID, DoctorID: 1}).Error; e != nil {
		t.Fatal(e)
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
			if tt.q == "王" && (got[0].RecordCount != 1 || got[0].PrescriptionCount != 1) {
				t.Fatalf("counts = records:%d prescriptions:%d, want 1/1", got[0].RecordCount, got[0].PrescriptionCount)
			}
		})
	}
}

func TestPatientRepositoryMerge(t *testing.T) {
	db, e := gorm.Open(sqlite.Open("file:patientmerge?mode=memory&cache=shared"), &gorm.Config{})
	if e != nil {
		t.Fatal(e)
	}
	if e = db.AutoMigrate(&model.Patient{}, &model.PatientMergeRecord{}, &model.MedicalRecord{}, &model.Prescription{}, &model.PrescriptionItem{}); e != nil {
		t.Fatal(e)
	}
	repo := NewPatientRepository(db)
	retained := model.Patient{RecordNo: "EMR-KEEP", Name: "张保留", Gender: "女", Age: 31, IDCard: "110101199001010011", Phone: "13800000011"}
	merged := model.Patient{RecordNo: "EMR-MERGE", Name: "张重复", Gender: "女", Age: 31, IDCard: "110101199001010012", Phone: "13800000012"}
	if e = db.Create(&retained).Error; e != nil {
		t.Fatal(e)
	}
	if e = db.Create(&merged).Error; e != nil {
		t.Fatal(e)
	}
	records := []model.MedicalRecord{
		{PatientID: merged.ID, DoctorID: 1, DepartmentID: 1, RecordType: "outpatient", Status: "draft"},
		{PatientID: merged.ID, DoctorID: 1, DepartmentID: 1, RecordType: "inpatient", Status: "draft"},
	}
	for i := range records {
		if e = db.Create(&records[i]).Error; e != nil {
			t.Fatal(e)
		}
	}
	prescriptions := []model.Prescription{
		{MedicalRecordID: records[0].ID, PatientID: merged.ID, DoctorID: 1},
		{MedicalRecordID: records[1].ID, PatientID: merged.ID, DoctorID: 1},
	}
	for i := range prescriptions {
		if e = db.Create(&prescriptions[i]).Error; e != nil {
			t.Fatal(e)
		}
	}

	logRow := &model.PatientMergeRecord{MergedByID: 1, MergedByName: "管理员", Reason: "身份证号录入错误导致重复建档"}
	if e = repo.Merge(retained.ID, merged.ID, logRow); e != nil {
		t.Fatalf("merge failed: %v", e)
	}
	var movedRecords int64
	if e = db.Model(&model.MedicalRecord{}).Where("patient_id = ?", retained.ID).Count(&movedRecords).Error; e != nil {
		t.Fatal(e)
	}
	if movedRecords != 2 {
		t.Fatalf("moved records = %d, want 2", movedRecords)
	}
	var movedPrescriptions int64
	if e = db.Model(&model.Prescription{}).Where("patient_id = ?", retained.ID).Count(&movedPrescriptions).Error; e != nil {
		t.Fatal(e)
	}
	if movedPrescriptions != 2 {
		t.Fatalf("moved prescriptions = %d, want 2", movedPrescriptions)
	}
	_, total, e := repo.Search("", 1, 10)
	if e != nil {
		t.Fatal(e)
	}
	if total != 1 {
		t.Fatalf("visible patients = %d, want 1", total)
	}
	if logRow.MergedRecordNo != "EMR-MERGE" || logRow.MedicalRecordCount != 2 || logRow.PrescriptionCount != 2 {
		t.Fatalf("unexpected merge log: %#v", logRow)
	}
	if e = repo.Merge(retained.ID, merged.ID, &model.PatientMergeRecord{MergedByID: 1, MergedByName: "admin", Reason: "重复合并必须拒绝"}); e == nil {
		t.Fatal("second merge of the same merged profile must fail")
	}
}

func TestPatientRepositoryMergeKeepsProfilesOnFailure(t *testing.T) {
	db, e := gorm.Open(sqlite.Open("file:patientmergefail?mode=memory&cache=shared"), &gorm.Config{})
	if e != nil {
		t.Fatal(e)
	}
	if e = db.AutoMigrate(&model.Patient{}, &model.PatientMergeRecord{}, &model.MedicalRecord{}, &model.Prescription{}, &model.PrescriptionItem{}); e != nil {
		t.Fatal(e)
	}
	repo := NewPatientRepository(db)
	source := model.Patient{RecordNo: "EMR-SOURCE", Name: "错误源档案", Gender: "女", Age: 31, IDCard: "110101199001010021", Phone: "13800000021"}
	if e = db.Create(&source).Error; e != nil {
		t.Fatal(e)
	}
	record := model.MedicalRecord{PatientID: source.ID, DoctorID: 1, DepartmentID: 1, RecordType: "outpatient", Status: "draft"}
	if e = db.Create(&record).Error; e != nil {
		t.Fatal(e)
	}

	e = repo.Merge(999999, source.ID, &model.PatientMergeRecord{MergedByID: 1, MergedByName: "admin", Reason: "目标不存在时必须回滚"})
	if e == nil {
		t.Fatal("expected merge failure when retained patient does not exist")
	}
	var sourcePatient model.Patient
	if e = db.First(&sourcePatient, source.ID).Error; e != nil {
		t.Fatalf("source profile must remain visible: %v", e)
	}
	var sourceRecords int64
	if e = db.Model(&model.MedicalRecord{}).Where("patient_id = ?", source.ID).Count(&sourceRecords).Error; e != nil {
		t.Fatal(e)
	}
	if sourceRecords != 1 {
		t.Fatalf("source records = %d, want 1", sourceRecords)
	}
	var mergeLogs int64
	if e = db.Model(&model.PatientMergeRecord{}).Count(&mergeLogs).Error; e != nil {
		t.Fatal(e)
	}
	if mergeLogs != 0 {
		t.Fatalf("merge logs = %d, want 0", mergeLogs)
	}
}
