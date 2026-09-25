export interface ApiResult<T>{code:number;message:string;data:T}
export interface User{id:number;username:string;name:string;role:'admin'|'doctor'|'nurse';department_id?:number;active:boolean}
export interface Department{id:number;code:string;name:string;description:string}
export interface Patient{id:number;record_no:string;name:string;gender:string;age:number;id_card:string;phone:string;allergies:string;medical_history:string;created_at:string}
export interface MedicalRecord{id:number;patient_id:number;patient?:Patient;doctor?:User;department?:Department;department_id:number;record_type:string;chief_complaint:string;present_illness:string;past_history:string;physical_exam:string;auxiliary_exam:string;diagnosis:string;treatment_plan:string;rich_content:string;status:string;created_at:string}
export interface Prescription{id:number;medical_record_id:number;patient_id:number;status:string;items:PrescriptionItem[];created_at:string}
export interface PrescriptionItem{drug_name:string;specification:string;dosage:string;frequency:string;duration:string}
