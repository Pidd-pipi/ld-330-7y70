package dto

type LoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=64"`
	Password string `json:"password" validate:"required,min=6,max=100"`
}
type PatientInput struct {
	Name           string `json:"name" validate:"required,max=64"`
	Gender         string `json:"gender" validate:"required,oneof=男 女 其他"`
	Age            int    `json:"age" validate:"gte=0,lte=150"`
	IDCard         string `json:"id_card" validate:"required,min=10,max=32"`
	Phone          string `json:"phone" validate:"required,min=6,max=32"`
	Allergies      string `json:"allergies"`
	MedicalHistory string `json:"medical_history"`
}
type PatientMergePreview struct {
	RetainedPatientID uint `json:"retained_patient_id" validate:"required"`
	MergedPatientID   uint `json:"merged_patient_id" validate:"required"`
}
type PatientMergeRequest struct {
	RetainedPatientID uint   `json:"retained_patient_id" validate:"required,nefield=MergedPatientID"`
	MergedPatientID   uint   `json:"merged_patient_id" validate:"required"`
	Reason            string `json:"reason" validate:"required,min=4,max=500"`
}
type PatientMergeSummary struct {
	ID                uint   `json:"id"`
	RecordNo          string `json:"record_no"`
	Name              string `json:"name"`
	Gender            string `json:"gender"`
	Age               int    `json:"age"`
	IDCard            string `json:"id_card"`
	Phone             string `json:"phone"`
	RecordCount       int64  `json:"record_count"`
	PrescriptionCount int64  `json:"prescription_count"`
	MergedIntoID      *uint  `json:"merged_into_id,omitempty"`
}
type PatientMergePreviewData struct {
	Retained              PatientMergeSummary `json:"retained"`
	Merged                PatientMergeSummary `json:"merged"`
	WillMoveRecords       int64               `json:"will_move_records"`
	WillMovePrescriptions int64               `json:"will_move_prescriptions"`
}
type RecordInput struct {
	PatientID      uint   `json:"patient_id" validate:"required"`
	DepartmentID   uint   `json:"department_id" validate:"required"`
	RecordType     string `json:"record_type" validate:"required,oneof=outpatient inpatient"`
	ChiefComplaint string `json:"chief_complaint" validate:"required"`
	PresentIllness string `json:"present_illness"`
	PastHistory    string `json:"past_history"`
	PhysicalExam   string `json:"physical_exam"`
	AuxiliaryExam  string `json:"auxiliary_exam"`
	Diagnosis      string `json:"diagnosis" validate:"required"`
	TreatmentPlan  string `json:"treatment_plan"`
	RichContent    string `json:"rich_content"`
}
type OrderInput struct {
	MedicalRecordID uint   `json:"medical_record_id" validate:"required"`
	Type            string `json:"type" validate:"required,oneof=long temporary"`
	Content         string `json:"content" validate:"required"`
}
type PrescriptionItemInput struct {
	DrugName      string `json:"drug_name" validate:"required"`
	Specification string `json:"specification"`
	Dosage        string `json:"dosage" validate:"required"`
	Frequency     string `json:"frequency" validate:"required"`
	Duration      string `json:"duration" validate:"required"`
}
type PrescriptionInput struct {
	MedicalRecordID uint                    `json:"medical_record_id" validate:"required"`
	Items           []PrescriptionItemInput `json:"items" validate:"required,min=1,dive"`
}
type DepartmentInput struct {
	Code        string `json:"code" validate:"required,max=32"`
	Name        string `json:"name" validate:"required,max=100"`
	Description string `json:"description"`
}
type UserInput struct {
	Username     string `json:"username" validate:"required,min=3"`
	Password     string `json:"password" validate:"required,min=6"`
	Name         string `json:"name" validate:"required"`
	Role         string `json:"role" validate:"required,oneof=admin doctor nurse"`
	DepartmentID *uint  `json:"department_id"`
}
type CatalogInput struct {
	Name          string `json:"name" validate:"required"`
	Specification string `json:"specification"`
	Unit          string `json:"unit"`
	Code          string `json:"code"`
	RecordType    string `json:"record_type"`
	Content       string `json:"content"`
}
