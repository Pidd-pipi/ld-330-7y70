package repository

import (
	"github.com/blueship581/gbemr/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func TestPatientRepositorySearch(t *testing.T) {
	db, e := gorm.Open(sqlite.Open("file:patientrepo?mode=memory&cache=shared"), &gorm.Config{})
	if e != nil {
		t.Fatal(e)
	}
	if e = db.AutoMigrate(&model.Patient{}); e != nil {
		t.Fatal(e)
	}
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
		})
	}
}
