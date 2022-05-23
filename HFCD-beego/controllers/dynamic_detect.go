package controllers

import (
	"HFCD-beego/models/dynamicModule"
	"fmt"
	"github.com/astaxie/beego"
)

type DynamicDetectController struct {
	beego.Controller
}

func (c *DynamicDetectController) Get() {
	ReportPath := "./models/TestReport/test_report.txt"
	chaincode := c.GetString("newcode")
	version := c.GetString("version")
	Init(ReportPath, chaincode, version)
	chaincode_path := "./models/chaincode.go"
	writeToFile(chaincode, chaincode_path)
	_ = ExeSysCommand("python3 ./models/dynamicModule/Symbolic.py")
	//time.Sleep(30 * time.Second) //read corpus
	cost := dynamicModule.Main_2(ReportPath,chaincode_path)
	str := "*Dynamic cost: " + fmt.Sprint(cost)
	Statistics(ReportPath, str)
	data := &JSONS{
		Getfile(ReportPath),
	}
	c.Data["json"] = data
	c.ServeJSON()
	fmt.Println("检测执行完成！")
}
