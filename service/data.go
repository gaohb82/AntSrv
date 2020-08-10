package service

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

var (
	G_DbsName []string //记录系统中 洗消数据库 别名
)

type Data_Controller struct {
	beego.Controller
}

func (this *Data_Controller) Get() {
	log.Println("URL:data/ METHOD:GET PARAM:" + this.Ctx.Input.Param("0"))
	if this.Ctx.Input.Param("0") == "operator" {
		this.Data["json"] = this.operator_GetAll()
		this.ServeJSON()
	} else if this.Ctx.Input.Param("0") == "doctor" {
		this.Data["json"] = this.doctor_GetAll()
		this.ServeJSON()
	} else if this.Ctx.Input.Param("0") == "device" {
		this.Data["json"] = this.device_GetAll()
		this.ServeJSON()
	} else if strings.ToLower(this.Ctx.Input.Param("0")) == "lastrecordbyeid" { //返回指定内镜编号 的最后一次洗消记录
		if this.Ctx.Input.Param("1") != "" {
			this.Data["json"] = this.Getlastrecordbyeid(this.Ctx.Input.Param("1"), this.Ctx.Input.Param("2"))
			this.ServeJSON()
		}
	} else if strings.ToLower(this.Ctx.Input.Param("0")) == "lastrecordbyrid" { //返回指定洗消记录编号的记录
		if this.Ctx.Input.Param("1") != "" {
			this.Data["json"] = this.Getlastrecordbyrid(this.Ctx.Input.Param("1"), this.Ctx.Input.Param("2"))
			this.ServeJSON()
		}
	}
}

func (this *Data_Controller) Post() {

	// 上传的 Ip
	ip := strings.Split(this.Ctx.Request.RemoteAddr, ":")[0]

	log.Println("URL:data/ METHOD:Post PARAM0:" + this.Ctx.Input.Param("0"))

	//刷卡
	if this.Ctx.Input.Param("0") == "brushcards" {
		var dat map[string]interface{}
		if err := json.Unmarshal(this.Ctx.Input.RequestBody, &dat); err == nil {

			if dat["ip"] != nil {
				ip = dat["ip"].(string)
			}

			log.Println("URL 刷卡操作:", ip, dat["number"].(string))
		}
	}
}

func (this *Data_Controller) operator_GetAll() (res map[string]STRUCT_OPERATOR_INFO) {

	_operator := make(map[string]int)
	_operator_Queue := make(map[string]STRUCT_OPERATOR_INFO)
	n := 0

	for _, dbstr := range G_DbsName {

		operator := []orm.Params{}
		o := orm.NewOrm()
		o.Using(dbstr)
		qs := o.QueryTable("operator_info")
		qs.OrderBy("id").Values(&operator)

		for _, val := range operator {
			if _operator[val["Number"].(string)] == 0 {
				_operator[val["Number"].(string)] = 1
				_operator_Queue[strconv.Itoa(n+1)] = STRUCT_OPERATOR_INFO{val["Number"].(string), val["Name"].(string)}
				n++
			}
		}
	}
	return _operator_Queue
}

func (this *Data_Controller) doctor_GetAll() (res map[string]STRUCT_DOCTOR_INFO) {

	_doctor := make(map[string]int)
	_doctor_Queue := make(map[string]STRUCT_DOCTOR_INFO)
	n := 0

	for _, dbstr := range G_DbsName {

		doctor := []orm.Params{}
		o := orm.NewOrm()
		o.Using(dbstr)
		qs := o.QueryTable("doctor_info")
		qs.OrderBy("id").Values(&doctor)

		for _, val := range doctor {
			if _doctor[val["Name"].(string)] == 0 {
				_doctor[val["Name"].(string)] = 1
				_doctor_Queue[strconv.Itoa(n+1)] = STRUCT_DOCTOR_INFO{val["Name"].(string)}
				n++
			}
		}
	}
	return _doctor_Queue
}

func (this *Data_Controller) device_GetAll() (res map[string]STRUCT_DEVICE_INFO) {

	_device := make(map[string]int)
	_device_Queue := make(map[string]STRUCT_DEVICE_INFO)
	n := 0

	for _, dbstr := range G_DbsName {

		device := []orm.Params{}
		o := orm.NewOrm()
		o.Using(dbstr)
		qs := o.QueryTable("device_info")
		qs.OrderBy("id").Values(&device)

		for _, val := range device {
			if _device[val["Endoscope_number"].(string)] == 0 {
				_device[val["Endoscope_number"].(string)] = 1
				_device_Queue[strconv.Itoa(n+1)] = STRUCT_DEVICE_INFO{val["Endoscope_number"].(string),
					val["Endoscope_type"].(string),
					val["Endoscope_info"].(string),
					strconv.FormatInt(val["Status"].(int64), 10)}
				n++
			}
		}
	}
	return _device_Queue
}

func (this *Data_Controller) Getlastrecordbyeid(_number string, _patient string) (res map[string]AntDb) {

	sqlstrpart := " begin_time <> end_time and patient_name == '' and "
	if beego.AppConfig.String("includeprocessing") == "true" {
		sqlstrpart = " patient_name == '' and "
	}

	sqlstr := "select begin_time from ant where " + sqlstrpart + " endoscope_number='" + _number +
		"'  ORDER BY ID DESC  limit 0,1 "
	ants := []orm.Params{}
	log.Print(sqlstr)
	dbname := "defautl"
	var lasttime int64
	lasttime = 0
	for _, dbstr := range G_DbsName {
		o := orm.NewOrm()
		o.Using(dbstr)
		o.Raw(sqlstr).Values(&ants)
		if len(ants) > 0 {
			loc, _ := time.LoadLocation("Local")
			t, _ := time.ParseInLocation("2006-01-02 15:04:05", ants[0]["begin_time"].(string), loc)
			if t.Unix() > lasttime {
				dbname = dbstr
				lasttime = t.Unix()
			}
		}
	}
	if lasttime == 0 {
		return make(map[string]AntDb)
	}
	log.Println(_number, "内窥镜,最后一次清洗在", dbname, "库中。")

	o := orm.NewOrm()
	o.Using(dbname)
	sqlstr = "select * from ant where " + sqlstrpart + " endoscope_number='" + _number +
		"'  ORDER BY ID DESC  limit 0,1 "
	o.Raw(sqlstr).Values(&ants)
	res = make(map[string]AntDb)

	for i, val := range ants {
		steps := make(map[string]AntStepDb)
		_db := AntDb{val["id"].(string), val["number"].(string), val["endoscope_number"].(string),
			val["endoscope_type"].(string), val["operator"].(string), val["patient_name"].(string),
			val["doc_name"].(string), val["diseases"].(string), val["begin_time"].(string),
			val["end_time"].(string), val["total_cost_time"].(string), val["endoscope_info"].(string), steps}

		var ant_steps []orm.Params
		o.Raw("select * from ant_step where number='" + val["number"].(string) + "' ORDER BY id").Values(&ant_steps)
		for x, val_step := range ant_steps {
			_step := AntStepDb{val_step["id"].(string), "", "", "", ""}
			if val_step["number"] != nil {
				_step.S_number = val_step["number"].(string)
			}
			if val_step["step"] != nil {
				_step.S_step = val_step["step"].(string)
			}
			if val_step["cost_time"] != nil {
				_step.S_cost_time = val_step["cost_time"].(string)
			}
			if val_step["washing_machine"] != nil {
				_step.S_washing_machine = val_step["washing_machine"].(string)
			}
			steps[strconv.Itoa(x+1)] = _step
		}
		res[strconv.Itoa(i+1)] = _db
	}
	log.Println(_number, "内镜镜,最后一次洗消记录返回。")

	if len(_patient) > 0 {
		sqlstr = "update ant set patient_name='" + _patient + "' where number='" + ants[0]["number"].(string) + "' and patient_name=''"
		o := orm.NewOrm()
		o.Using(dbname)
		_, err := o.Raw(sqlstr).Exec()
		if err != nil {
			log.Println(err)
		}
		log.Println("更新洗消记录", ants[0]["number"].(string), "病人信息为:"+_patient)
	}
	return res
}

func (this *Data_Controller) Getlastrecordbyrid(_number string, _patient string) (res map[string]AntDb) {

	res = make(map[string]AntDb)

	sqlstr := "select * from ant where number='" + this.Ctx.Input.Param("1") + "'  ORDER BY ID DESC  limit 0,1 "

	for _, dbstr := range G_DbsName {
		ants := []orm.Params{}
		o := orm.NewOrm()
		o.Using(dbstr)
		o.Raw(sqlstr).Values(&ants)
		if len(ants) > 0 {

			log.Println(_number, "洗消编号的记录在:", dbstr, "库中。")

			for i, val := range ants {
				steps := make(map[string]AntStepDb)
				_db := AntDb{val["id"].(string), val["number"].(string), val["endoscope_number"].(string),
					val["endoscope_type"].(string), val["operator"].(string), val["patient_name"].(string),
					val["doc_name"].(string), val["diseases"].(string), val["begin_time"].(string),
					val["end_time"].(string), val["total_cost_time"].(string), val["endoscope_info"].(string), steps}

				var ant_steps []orm.Params
				o.Raw("select * from ant_step where number='" + val["number"].(string) + "' ORDER BY id").Values(&ant_steps)
				for x, val_step := range ant_steps {
					_step := AntStepDb{val_step["id"].(string), "", "", "", ""}
					if val_step["number"] != nil {
						_step.S_number = val_step["number"].(string)
					}
					if val_step["step"] != nil {
						_step.S_step = val_step["step"].(string)
					}
					if val_step["cost_time"] != nil {
						_step.S_cost_time = val_step["cost_time"].(string)
					}
					if val_step["washing_machine"] != nil {
						_step.S_washing_machine = val_step["washing_machine"].(string)
					}
					steps[strconv.Itoa(x+1)] = _step
				}
				res[strconv.Itoa(i+1)] = _db
			}

			if len(_patient) > 0 {
				sqlstr = "update ant set patient_name='" + _patient + "' where number='" + ants[0]["number"].(string) + "'"
				_, err := o.Raw(sqlstr).Exec()
				if err != nil {
					log.Println(err)
				}
				log.Println("更新洗消记录", ants[0]["number"].(string), "病人信息为:"+_patient)
			}
			break
		}
	}

	return res
}
