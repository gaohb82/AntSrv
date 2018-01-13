package service

import "time"

type Ant struct {
	Id               int
	Number           string
	Endoscope_number string
	Endoscope_type   string
	Operator         string
	Patient_name     string
	Doc_name         string
	Diseases         int
	Begin_time       time.Time
	End_time         time.Time
	Total_cost_time  int
	Endoscope_info   string
}

type AntStep struct {
	Id              int
	Number          string
	Step            string
	Cost_time       int
	Washing_machine string
}

type AntDb struct {
	A_id               string `json:"id"`
	A_number           string `json:"number"`
	A_endoscope_number string `json:"endoscope_number"`
	A_endoscope_type   string `json:"endoscope_type"`
	A_operator         string `json:"operator"`
	A_patient_name     string `json:"patient"`
	A_doc_name         string `json:"doctor"`
	A_diseases         string `json:"diseases"`
	A_begin_time       string `json:"begin_time"`
	A_end_time         string `json:"end_time"`
	A_total_cost_time  string `json:"total_cost_time"`
	A_endoscope_info   string `json:"endoscope_info"`

	A_steps map[string]AntStepDb `json:"steps"`
}

type AntStepDb struct {
	S_id              string `json:"step_id"`
	S_number          string `json:"step_number"`
	S_step            string `json:"step"`
	S_cost_time       string `json:"step_cost_time"`
	S_washing_machine string `json:"step_washing_machine"`
}

type DeviceInfo struct {
	Id               int
	Endoscope_number string
	Endoscope_type   string
	Endoscope_info   string
	Status           int
}

type STRUCT_DEVICE_INFO struct {
	Number string `json:"number"`
	Type   string `json:"type"`
	Info   string `json:"info"`
	Status string `json:"status"`
}

type DoctorInfo struct {
	Id   int
	Name string
}

type STRUCT_DOCTOR_INFO struct {
	Name string `json:"name"`
}

type Program struct {
	Id              int
	Name            string
	Total_cost_time int
}

type ProgramList struct {
	Id        int
	Name      string
	Step      string
	Cost_time int
}

type STRUCT_PROGRAM_INFO struct {
	Name          string                              `json:"name"`
	TotalCostTime int64                               `json:"TotalCostTime"`
	StepList      map[string]STRUCT_PROGRAM_LIST_INFO `json:"StepList"`
}

type STRUCT_PROGRAM_LIST_INFO struct {
	Step      string `json:"Step"`
	Cost_time int64  `json:"Step_time"`
}

type OperatorInfo struct {
	Id     int
	Number string
	Name   string
}

type STRUCT_OPERATOR_INFO struct {
	Number string `json:"number"`
	Name   string `json:"name"`
}

type Repair struct {
	Id          int
	Device_id   string
	Finish_id   int64
	Repair_time time.Time
	Repair_name string
	Comment     string
	Createtime  time.Time
}

type RepairReason struct {
	Id        int
	Repair_id int64
	Reason    string
}

type RepairFinish struct {
	Id          int
	Device_id   string
	Finish_time time.Time
	Cost_amount float64
	Comment     string
	Createtime  time.Time
}

type TimePlan struct {
	Id                 int
	Name               string
	Rinse_time         int
	Rinse_solution     string
	Disinfect_time     int
	Disinfect_solution string
}