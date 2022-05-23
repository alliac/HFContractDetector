package ast_graph

import (
	"encoding/json"
	"fmt"
	_ "github.com/emp-toolkit/jialanli/lacia/utils"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"strconv"
)

type ChainCode struct {
}
type FoodInfo struct {
	FoodID      string  `json:FoodID`
	FoodProInfo ProInfo `json:FoodProInfo`
	//FoodIngInfo []IngInfo `json:FoodIngInfo`
	//FoodLogInfo LogInfo   `json:FoodLogInfo`
}
type ProInfo struct {
	FoodName     string `json:FoodName`
	FoodSpec     string `json:FoodSpec`
	FoodMFGDate  string `json:FoodMFGDate`
	FoodEXPDate  string `json:FoodEXPDate`
	FoodLOT      string `json:FoodLOT`
	FoodQSID     string `json:FoodQSID`
	FoodMFRSName string `json:FoodMFRSName`
	FoodProPrice string `json:FoodProPrice`
	FoodProPlace string `json:FoodProPlace`
}

//使用继承的函数和变量
type TestA struct {
	test int ``
}
type TestB struct {
	test int ``
}
type TestC struct {
	testA TestA ``
	testB TestB ``
}

func (t *TestA) FuncTest() {

}
func (t *TestB) FuncTest() {

}

// @title    Init
// @description   Init
// @auth       -             -
// @param     stub        shim.ChaincodeStubInterface         " "
// @return    Response        pb.Response         " "
func (a *ChainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("init is called")
	//写后读漏洞
	stub.PutState("key", []byte("1"))
	stub.GetState("key")
	return shim.Success(nil)
}
func (a *ChainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("invoke is called")
	fn, args := stub.GetFunctionAndParameters()
	if fn == "AddProInfo" {
		return a.AddProInfo(stub, args)
	}
	if fn == "GetProInfo" {
		return a.GetProInfo(stub, args)
	}

	return shim.Error("Recevied unkown function invocation")
}

//@title addProInfo
func (a *ChainCode) AddProInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error
	var FoodInfos FoodInfo

	if len(args) != 10 {
		return shim.Error("Incorrect number of arguments.")
	}
	FoodInfos.FoodID = args[0]

	if FoodInfos.FoodID == "" {
		return shim.Error("FoodID can not be empty.")
	}

	FoodInfos.FoodProInfo.FoodName = args[1]
	FoodInfos.FoodProInfo.FoodSpec = args[2]
	FoodInfos.FoodProInfo.FoodMFGDate = args[3]
	FoodInfos.FoodProInfo.FoodEXPDate = args[4]
	FoodInfos.FoodProInfo.FoodLOT = args[5]
	FoodInfos.FoodProInfo.FoodQSID = args[6]
	FoodInfos.FoodProInfo.FoodMFRSName = args[7]
	FoodInfos.FoodProInfo.FoodProPrice = args[8]
	FoodInfos.FoodProInfo.FoodProPlace = args[9]
	Dividend, _ := strconv.Atoi(FoodInfos.FoodProInfo.FoodProPrice)
	Divisor, _ := strconv.Atoi(FoodInfos.FoodProInfo.FoodSpec)
	var channel chan int = make(chan int)
	single_price := MutiClose(Dividend, Divisor, channel)
	fmt.Println(single_price)

	ProInfosJSONasBytes, err := json.Marshal(FoodInfos)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(FoodInfos.FoodID, ProInfosJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(ProInfosJSONasBytes)
}
func (a *ChainCode) GetProInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//跨通道链码调用
	stub.InvokeChaincode("", nil, "channel")
	//未初始化存储指针
	var UnInitArr []string
	var UnInitStr TestA
	var InitStr1 = TestA{}
	var InitStr2 TestA
	InitStr2=TestA{}
	var UnInitMap map[string]bool
	fmt.Println(UnInitArr)
	fmt.Println(UnInitStr)
	fmt.Println(UnInitMap)
	fmt.Println(InitStr1)
	fmt.Println(InitStr2)
	return shim.Success(nil)
}

//@title MutiClose
func MutiClose(Dividend int, Divisor int, channel chan int) int {
	defer close(channel)
	if Divisor == 0 {
		close(channel)
		return 0
	}
	return Dividend / Divisor
}