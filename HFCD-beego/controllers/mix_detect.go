package controllers

import (
	"HFCD-beego/models/ast_graph"
	"HFCD-beego/models/dynamicModule"
	"fmt"
	"github.com/astaxie/beego"
)

type MixDetectController struct {
	beego.Controller
}

func (c *MixDetectController) Get() {
	//执行构建初始语料库
	ReportPath := "./models/TestReport/test_report.txt"
	chaincode := c.GetString("newcode")
	version := c.GetString("version")
	Init(ReportPath, chaincode, version)
	chaincode_path := "./models/chaincode.go"
	writeToFile(chaincode, chaincode_path)
	cost1 := ast_graph.Main_1(ReportPath)
	fmt.Println("静态检测执行完成！")
	_ = dynamicModule.ExeSysCommand("python3 ./models/dynamicModule/Symbolic.py")
	//time.Sleep(30 * time.Second) //read corpus
	cost2 := dynamicModule.Main_2(ReportPath,chaincode_path)
	fmt.Println("动态检测执行完成！")
	str := "*Static cost: " + fmt.Sprint(cost1) + ",Dynamic cost: " + fmt.Sprint(cost2) + ",Total cost: " + fmt.Sprint(cost1+cost2)
	Statistics(ReportPath, str)
	data := &JSONS{
		Getfile(ReportPath),
	}
	c.Data["json"] = data
	c.ServeJSON()
	fmt.Println("检测执行完成！")
}
