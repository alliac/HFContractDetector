package dynamicModule

import (
	"HFCD-beego/models"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"time"
)

/*
污点分析
*/
func Taint(data [][][]byte, ReportPath string) int {
	timeout := make(chan bool, 1)
	ch := make(chan int, 5)
	go timer(timeout)
	go mockSC(data, ch)
	//无限循环漏洞
	for {
		select {
		case _ = <-ch:
			//WriteToReport(ReportPath, "不存在无限循环")
			return 0
		case _ = <-timeout:
			WriteToReport(ReportPath, "存在无限循环")
			return 1
		}
	}
}

//调用智能合约
func mockSC(data [][][]byte, ch chan int) {
	//fmt.Println("mockSC called")
	cc := new(models.ChainCode)                    // 创建Chaincode对象
	stub := shim.NewMockStub("./models/chaincode", cc) // 创建MockStub对象
	for i := 0; i < len(data); i++ {
		if string(data[i][0]) == "Init" {
			var tmp [][]byte
			for j := 1; j < len(data[i]); j++ {
				tmp = append(tmp, data[i][j])
			}
			//fmt.Println(tmp)
			res := stub.MockInit("1", tmp)
			if res.Status != 200 {
				ch <- 1
				return
			}
		} else {
			//fmt.Println(data[i])
			res := stub.MockInvoke("1", data[i])
			if res.Status != 200 {
				ch <- 1
				return
			}
		}
	}
	ch <- 1
}

//计时器
func timer(timeout chan bool) {
	time.Sleep(time.Duration(5 * time.Second))
	timeout <- true
}
