package bugs

import (
	"fmt"
	mapset "github.com/golang-set"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
)

//bugs about pkg
var errPkg = make(map[string]string)

//标志位：标记是否使用了crypto包
var cryptoFlag = false

//范围查询方法
var RangeQueryFunc [3]string

//内置包
var infuncs = make(map[string]bool)

//putstatekey:key[ox form],value=pos;PosToKey:key=pos,value=key[string]
var putstate = make(map[string]int64)
var PosToKey = make(map[int64]string)

//StructToType:key=struct name,value=struct type
var StructToType = make(map[string]string)

//StructFields:key=struct name,value=struct fields
var StructFields = make(map[string]mapset.Set)

//判断变量是否被赋值
var AssignValue = make(map[string]bool)

//存在多重继承的结构体--它包含的超级类
var StructInherit = make(map[string][10]string)

//test report path
var ReportPath string

func InitDtct(reportPath string) {
	ReportPath = reportPath
	errPkg["os/exec"] = "系统命令执行漏洞"
	errPkg["net/http"] = "Web服务漏洞"
	errPkg["io/ioutil"] = "外部文件访问"
	errPkg["os"] = "外部文件访问"
	errPkg["crypto/rand"] = "随机数生成漏洞"
	errPkg["math/rand"] = "随机数生成漏洞"
	errPkg["time"] = "系统时间戳漏洞"
	//内部包列表
	getDirlist()
	//go语言内置包
	infuncs["close"] = true
	infuncs["len"] = true
	infuncs["cap"] = true
	infuncs["new"] = true
	infuncs["make"] = true
	infuncs["copy"] = true
	infuncs["append"] = true
	infuncs["print"] = true
	infuncs["println"] = true
	infuncs["panic"] = true
	infuncs["recover"] = true
	infuncs["delete"] = true
	//范围查询方法
	RangeQueryFunc[0] = "GetQueryResult"
	RangeQueryFunc[1] = "GetHistoryForKey"
	RangeQueryFunc[2] = "GetPRivateDataQueryResult"
}

/*
善后
*/
func FinshDtct(path string, ReportPath string) {
	//未处理的错误漏洞
	//NoAST(path, ReportPath)
	//未加密的敏感数据
	if !cryptoFlag {
		WriteToReport(ReportPath, "未使用crypto,未加密的敏感数据")
	}
	// 未初始化存储指针
	for k, v := range AssignValue {
		if !v {
			WriteToReport(ReportPath, "未初始化存储指针: "+k+" 未初始化")
		}
	}
	//使用继承的函数或变量
	for k, v := range StructInherit {
		id := len(v)
		for i := 0; i < id; i++ {
			for j := i + 1; j < id; j++ {
				if v[i] == "" || v[j] == "" || StructFields[v[j]].Cardinality() < 1 || StructFields[v[i]].Cardinality() < 1 {
					continue
				}
				res := StructFields[v[i]].Intersect(StructFields[v[j]]).Cardinality()
				if res > 0 {
					WriteToReport(ReportPath, "使用继承的函数和变量漏洞:"+k+"结构体多重继承"+v[i]+"和"+v[j]+",且超级类具有相同名称的方法或变量"+fmt.Sprint(StructFields[v[i]].Intersect(StructFields[v[j]])))
				}
			}
		}
	}
	errPkg = make(map[string]string)
	cryptoFlag = false
	RangeQueryFunc = [3]string{"", "", ""}
	infuncs = make(map[string]bool)
	putstate = make(map[string]int64)
	PosToKey = make(map[int64]string)
	StructToType = make(map[string]string)
	StructFields = make(map[string]mapset.Set)
	AssignValue = make(map[string]bool)
	StructInherit = make(map[string][10]string)
	packageExist = make(map[string]bool)
}

/*
adjust if str exist []str
*/
func in(target string, str_array []string) bool {
	sort.Strings(str_array)
	index := sort.SearchStrings(str_array, target)
	//index的取值：[0,len(str_array)]
	if index < len(str_array) && str_array[index] == target { //需要注意此处的判断，先判断 &&左侧的条件，如果不满足则结束此处判断，不会再进行右侧的判断
		return true
	}
	return false
}

//print test report
func WriteToReport(ReportPath string, s string) {
	//f, err := os.Open(ReportPath)
	f, err := os.Open(ReportPath)
	if err != nil {
		log.Fatal(err)
		return
	}
	content, _ := ioutil.ReadAll(f)
	content_s := string(content)
	if !strings.Contains(content_s, s) {
		//f, _ := os.OpenFile(ReportPath, syscall.O_APPEND, 0666)
		f, _ := os.OpenFile(ReportPath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		_, err := f.WriteString(s + "\n")
		if err != nil {
			fmt.Println("090")
			log.Fatal(err)
		}
		if err = f.Close(); err != nil {
			log.Fatal(err)
		}
	}
	_,_ = f.WriteString("++++++++++++++++++++++++++" + "\n")
	defer f.Close()
}
