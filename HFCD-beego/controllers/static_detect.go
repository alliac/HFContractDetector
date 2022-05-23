package controllers

import (
	"HFCD-beego/models/ast_graph"
	"fmt"
	"github.com/astaxie/beego"
	"time"
)

type StaticDetectContronller struct {
	beego.Controller
}
type JSONS struct {
	//必须的大写开头
	Result string
}

func (c *StaticDetectContronller) Get() {
	ReportPath := "./models/TestReport/test_report.txt"
	chaincode := c.GetString("newcode")
	version := c.GetString("version")
	Init(ReportPath, chaincode, version)
	chaincode_path := "./models/chaincode.go"
	writeToFile(chaincode, chaincode_path)
	var cost time.Duration
	cost = ast_graph.Main_1(ReportPath)
	str := "*Static cost: " + fmt.Sprint(cost)
	Statistics(ReportPath, str)
	data := &JSONS{
		Getfile(ReportPath),
	}
	c.Data["json"] = data
	c.ServeJSON()
	fmt.Println("检测执行完成！")
}
