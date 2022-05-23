package bugs

import (
	"fmt"
	mapset "github.com/golang-set"
	lacia "github.com/jialanli/lacia/utils"
	"go/ast"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

var MapFlag = false
var PrivateFlag = false

/*
FindBugs comment
*/
func FindBugs(n ast.Node, name string) {
	switch m := n.(type) {
	//记录函数起止位置，区分全局变量和局部变量
	case *ast.FuncDecl:
		RecordPos(fmt.Sprint(m.Name), int(m.Pos()), int(m.End()))
		//注释漏洞-dynamic
		//IsComment(fmt.Sprint(m.Name),ReportPath)
		//使用继承的函数或变量
		if m.Recv != nil {
			for _, f := range m.Recv.List {
				args := lacia.SplitByManyStrWith(fmt.Sprint(f.Type), []string{" ", "}"})
				StructFields[args[1]].Add(fmt.Sprint(m.Name))
			}
		}
	//内置函数漏洞
	case *ast.ExprStmt:
		var infunc string
		infunc = lacia.SplitByManyStrWith(fmt.Sprint(m.X), []string{" ", "{"})[1]
		//fmt.Println(infunc)
		if infuncs[infunc] {
			WriteToReport(ReportPath, "建议不使用内置函数"+infunc)
		}
	//跨通道链码调用
	case *ast.SelectorExpr:
		//fmt.Println(m.Sel)
		//if strings.Contains(fmt.Sprint(m.Sel), "InvokeChaincode") {
		//	fmt.Println("使用InvokeChaincode跨通道调用链码,可能存在隐私数据安全风险")
		//}
	//字段声明漏洞（与全局变量类似，默认链码中变量声明即使用）,使用继承的函数和变量
	case *ast.StructType:
		var InStype [10]string
		id := 0
		//查找结构体名字
		var Sname string
		for k, v := range StructToType {
			if v == fmt.Sprint(m) {
				Sname = k
				break
			}
		}
		//查找结构体中字段名称及类型
		for i := 0; i < len(m.Fields.List); i++ {
			var field []string
			field = lacia.SplitByManyStrWith(fmt.Sprint(m.Fields.List[i]), []string{" ", "[", "]"})
			WriteToReport(ReportPath, "结构体中不建议声明变量（字段声明漏洞）"+field[1])
			//记录字段中结构体类型的变量个数及变量类型
			if len(field) > 2 && StructToType[field[2]] != "" {
				InStype[id] = field[2]
				id++
			}
			StructFields[Sname].Add(field[1])
		}
		if id > 1 {
			StructInherit[Sname] = InStype
		}
	//使用继承的函数和变量
	case *ast.TypeSpec:
		Sname := fmt.Sprint(m.Name)
		Stype := fmt.Sprint(m.Type)
		if StructToType[Sname] == "" {
			StructToType[Sname] = Stype
			StructFields[Sname] = mapset.NewSet()
		}
	//程序并发行漏洞
	case *ast.GoStmt:
		WriteToReport(ReportPath, "不建议使用goroutine，易造成程序并发性漏洞")
	case *ast.ChanType:
		WriteToReport(ReportPath, "不建议使用channel，易造成程序并发性漏洞")
	//映射结构迭代
	case *ast.MapType:
		MapFlag = true
	//未初始化存储指针
	case *ast.AssignStmt:
		args := lacia.SplitByManyStrWith(fmt.Sprint(m.Lhs), []string{"[", "]"})
		AssignValue[args[0]] = true
	//全局漏洞（按位置标识，通道变量会被认为是全局变量，所以排除掉通道变量
	case *ast.ValueSpec:
		//通道变量m.type的值以&开头，所以正则匹配为空表示非通道变量，之后再判断是否是全局变量
		//re := regexp.MustCompile("^&.+")
		//if re.FindString(fmt.Sprint(m.Type)) == "" && m.Type != nil && !MapFlag {
		//非通道变量判断更新
		if fmt.Sprintf("%T", m.Type) != "*ast.ChanType" && !MapFlag {
			if isExist(int(m.Pos()), int(m.End())) {
				WriteToReport(ReportPath, "存在全局变量漏洞："+fmt.Sprint(m.Names))
			}
		}
		//映射结构迭代
		if MapFlag {
			WriteToReport(ReportPath, "存在映射结构迭代："+fmt.Sprint(m.Names))
			MapFlag = false
			return
		}
		//未初始化存储指针
		switch p := m.Type.(type) {
		case *ast.ArrayType:
			args := lacia.SplitByManyStrWith(fmt.Sprint(m.Names), []string{"[", "]"})
			//声明时未初始化
			if p.Len == nil {
				AssignValue[args[0]] = false
			} else {
				//声明时初始化
				AssignValue[args[0]] = true
			}
		case *ast.MapType:
			args := lacia.SplitByManyStrWith(fmt.Sprint(m.Names), []string{"[", "]"})
			AssignValue[args[0]] = false
		case *ast.Ident:
			args := lacia.SplitByManyStrWith(fmt.Sprint(m.Names), []string{"[", "]"})
			AssignValue[args[0]] = false
		}
	case *ast.CallExpr:
		//fmt.Println(m.Fun)
		if infuncs[fmt.Sprint(m.Fun)] {
			WriteToReport(ReportPath, "建议不使用内置函数"+fmt.Sprint(m.Fun))
		}
		//写后读漏洞
		if strings.Contains(fmt.Sprint(m.Fun), "PutState") {
			for _, x := range m.Args {
				arg := fmt.Sprintf("%T", x)
				if arg == "*ast.BasicLit" {
					args := lacia.SplitByManyStrWith(fmt.Sprintf("%s", x), []string{"{", " ", "\""})
					putstate[args[len(args)-2]] = int64(x.Pos())
				}
			}
			//未使用的隐私数据机制
			if !PrivateFlag {
				WriteToReport(ReportPath, "未使用的隐私数据机制:建议使用EMP-toolkit对上链数据进行加密")
			}
		}
		if strings.Contains(fmt.Sprint(m.Fun), "GetState") {
			for _, x := range m.Args {
				arg := fmt.Sprintf("%T", x)
				if arg == "*ast.BasicLit" {
					args := lacia.SplitByManyStrWith(fmt.Sprintf("%s", x), []string{"{", " ", "\""})
					if putstate[args[len(args)-2]] != 0 && putstate[args[len(args)-2]] < int64(x.Pos()) {
						WriteToReport(ReportPath, "存在写后读漏洞: "+fmt.Sprint(args[len(args)-2]))
					}
				}
			}
		}
		//范围查询风险
		for i := 0; i < len(RangeQueryFunc); i++ {
			res := strings.Contains(fmt.Sprint(m.Fun), RangeQueryFunc[i])
			if res {
				WriteToReport(ReportPath, "存在范围查询风险："+fmt.Sprint(RangeQueryFunc[i]))
			}
		}
		//跨通道调用链码
		if strings.Contains(fmt.Sprint(m.Fun), "InvokeChaincode") {
			flag := 0
			//取通道名称
			for _, x := range m.Args {
				if flag < 2 {
					flag++
					continue
				}
				arg := fmt.Sprintf("%T", x)
				if arg == "*ast.BasicLit" {
					args := lacia.SplitByManyStrWith(fmt.Sprintf("%s", x), []string{"\""})
					WriteToReport(ReportPath, "使用InvokeChaincode调用"+args[1]+"通道链码,可能存在隐私数据安全风险")
				}
			}
		}
	}
	switch name {
	//获取代码导入包，对标漏洞：系统命令执行漏洞、外部库调用、Web服务漏洞、外部文件访问、随机数生成漏洞、系统时间戳漏洞、
	case "ImportSpec.Path":
		var name string
		name = lacia.SplitByManyStrWith(fmt.Sprint(n), []string{"\""})[1]
		if errPkg[name] != "" {
			WriteToReport(ReportPath, "错误使用"+name+","+errPkg[name])
		} else if strings.Contains(name, "crypto") {
			cryptoFlag = true
		}
		if !packageExist[name] {
			WriteToReport(ReportPath, "调用了外部库"+name+",外部库调用漏洞")
		}
		if strings.Contains(name, "emp-toolkit") {
			PrivateFlag = true
		}
	case "CallExpr.Args":
		args := lacia.SplitByManyStrWith(fmt.Sprint(n), []string{"{", " ", "\""})
		if len(args) > 2 {
			pos, _ := strconv.ParseInt(args[1], 10, 64)
			PosToKey[pos] = args[len(args)-2]
		}
	}
}

/*
借助errcheck，检测未处理的错误漏洞
*/
func NoAST(path string, ReportPath string) {
	fmt.Println("----------")
	//cmd := exec.Command("sh", "-c", "cd models;errcheck chaincode.go")
	//opBytes, _ := cmd.CombinedOutput()
	//fmt.Println(string(opBytes))
	//strs := lacia.SplitByManyStrWith(string(opBytes), []string{"\n"})
	//i := 0
	//for ; i < len(strs); i++ {
	//	str := lacia.SplitByManyStrWith(strs[i], []string{":"})
	//	if len(str) < 3 {
	//		continue
	//	}
	//	WriteToReport(ReportPath, "未处理的错误："+str[len(str)-1])
	//}
	//str := "python3 ./models/cmd.py \"errcheck ./models/chaincode.go\""
	cmd := exec.Command("python3", "./models/cmd.py", "errcheck ./models/chaincode.go")
	opBytes, _ := cmd.Output()
	fmt.Println(string(opBytes))
}

/*
借助godoc，检测注释漏洞
*/
func IsComment(funcname string, ReportPath string) {
	opBytes := exeSysCommand("go doc ast_graph." + funcname)
	if len((lacia.SplitByManyStrWith(string(opBytes), []string{"\n"}))) < 2 {
		WriteToReport(ReportPath, "注释标题不足以检查实现和使用情况:"+funcname+"函数未添加注释说明")
	}
}
func exeSysCommand(cmdStr string) string {
	cmd := exec.Command("sh", "-c", cmdStr)
	opBytes, _ := cmd.Output()
	return string(opBytes)
}

//打印结果
func Print(s string) {
	f, err := os.Open("1.txt")
	if err != nil {
		log.Fatal(err)
		return
	}
	content, _ := ioutil.ReadAll(f)
	content_s := string(content)
	if !strings.Contains(content_s, s) {
		f, _ := os.OpenFile("1.txt", syscall.O_APPEND, 0666)
		_, err := f.WriteString(s)
		if err != nil {
			log.Fatal(err)
		}
		if err = f.Close(); err != nil {
			log.Fatal(err)
		}
	}
}
