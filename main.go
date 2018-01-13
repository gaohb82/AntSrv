package main

import (
	"AntSrv/service"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {

	dbstr := "default"

	for {
		if len(beego.AppConfig.String(dbstr+"::ip")) <= 0 {
			break
		}

		dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", beego.AppConfig.String(dbstr+"::username"),
			beego.AppConfig.String(dbstr+"::password"), beego.AppConfig.String(dbstr+"::ip"),
			beego.AppConfig.String(dbstr+"::port"), beego.AppConfig.String(dbstr+"::dataname"))

		err := orm.RegisterDataBase(dbstr, "mysql", dns, 30, 30)
		if err != nil {
			fmt.Println("orm 注册失败，继续重试！")
			continue
		} else {
			service.G_DbsName = append(service.G_DbsName, dbstr)
			fmt.Println("orm registerDatabase:", dbstr, dns)
		}
		dbstr = fmt.Sprintf("db%d", len(service.G_DbsName)+1)
	}

	orm.RegisterModel(new(service.Ant), new(service.AntStep), new(service.DeviceInfo), new(service.DoctorInfo),
		new(service.Program), new(service.ProgramList), new(service.OperatorInfo), new(service.Repair), new(service.RepairReason),
		new(service.RepairFinish), new(service.TimePlan))

	fmt.Println("环境中共有洗消数据主机:", len(service.G_DbsName))
}

func main() {

	beego.BConfig.AppName = "AntService"
	beego.BConfig.WebConfig.AutoRender = false
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.BConfig.CopyRequestBody = true
	beego.BConfig.RunMode = "test"
	beego.BConfig.Listen.HTTPAddr = "127.0.0.1"
	beego.BConfig.Listen.HTTPPort = 8866

	//Service 第三方数据获取 及 管理端
	beego.Router("*", &service.Data_Controller{})
	beego.Router("/data/*", &service.Data_Controller{})
	beego.Run()

}
