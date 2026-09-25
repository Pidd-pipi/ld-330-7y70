-- 患者重复档案合并
-- 开发环境也会由 GORM AutoMigrate 自动执行；生产环境建议使用受控迁移工具执行并备份。

ALTER TABLE patients
    ADD COLUMN IF NOT EXISTS is_merged boolean NOT NULL DEFAULT false,
    ADD COLUMN IF NOT EXISTS merged_into_id bigint,
    ADD COLUMN IF NOT EXISTS merged_at timestamp with time zone;

CREATE INDEX IF NOT EXISTS idx_patients_is_merged ON patients (is_merged);
CREATE INDEX IF NOT EXISTS idx_patients_merged_into_id ON patients (merged_into_id);

CREATE TABLE IF NOT EXISTS patient_merge_logs (
    id bigserial PRIMARY KEY,
    keep_patient_id bigint NOT NULL,
    keep_record_no varchar(32) NOT NULL,
    keep_patient_name varchar(64) NOT NULL,
    merged_patient_id bigint NOT NULL,
    merged_record_no varchar(32) NOT NULL,
    merged_patient_name varchar(64) NOT NULL,
    moved_record_count bigint NOT NULL,
    moved_prescription_count bigint NOT NULL,
    reason text NOT NULL,
    operator_id bigint NOT NULL,
    operator_name varchar(64) NOT NULL,
    merged_at timestamp with time zone
);

CREATE INDEX IF NOT EXISTS idx_patient_merge_logs_keep_patient_id ON patient_merge_logs (keep_patient_id);
CREATE INDEX IF NOT EXISTS idx_patient_merge_logs_merged_patient_id ON patient_merge_logs (merged_patient_id);
CREATE INDEX IF NOT EXISTS idx_patient_merge_logs_operator_id ON patient_merge_logs (operator_id);
CREATE INDEX IF NOT EXISTS idx_patient_merge_logs_merged_at ON patient_merge_logs (merged_at);
