package dynamicModule

import (
	"fmt"
	lacia "github.com/jialanli/lacia/utils"
	"io/ioutil"
	"os"
	"strings"
	"unicode"
)

func DetectSc(ReportPath string, chaincode_path string) {
	f, err := os.Open(chaincode_path)
	if err != nil {
		return
	}
	content, _ := ioutil.ReadAll(f)
	contents := string(content)
	funcs := strings.Split(contents, "\nfunc")
	funcnames := make([]string, len(funcs)-1) //动态数组
	for i := 1; i < len(funcs); i++ {
		var funcname string
		var ref string
		strs := strings.Split(funcs[i], "{")
		str := strings.ReplaceAll(strs[0], " ", "")
		nstrs := lacia.SplitByManyStrWith(str, []string{"(", ")"})
		if strings.Contains(strs[0], "ChainCode") {
			funcname = nstrs[1]
			ref = nstrs[2]
		} else {
			funcname = nstrs[0]
			if len(nstrs) > 1 {
				ref = nstrs[1]
			}
		}
		if unicode.IsLower([]rune(funcname)[0]) {
			Funcname := Capitalize(funcname)
			contents = strings.ReplaceAll(contents, funcname, Funcname)
		}
		if strings.Contains(ref, "args") && !strings.Contains(funcs[i], "len(args)") {
			WriteToReport(ReportPath, funcname+"函数未检查输入参数")
		}
		funcnames[i-1] = funcname
	}
	writeToFile(contents, chaincode_path)
	for i := 0; i < len(funcnames); i++ {
		IsComment(funcnames[i], ReportPath)
	}

}

//借助godoc，检测注释漏洞
func IsComment(funcname string, ReportPath string) {
	if funcname == "main" {
		return
	}
	opBytes := ExeSysCommand("go doc ./models." + funcname)
	contents := lacia.SplitByManyStrWith(opBytes, []string{"\n"})
	if strings.Contains(contents[len(contents)-1], Capitalize(funcname)) {
		WriteToReport(ReportPath, "注释标题不足以检查实现和使用情况:"+funcname+"函数未添加注释说明")
	}
}
func writeToFile(msg string, URL string) {
	if err := ioutil.WriteFile(URL, []byte(msg), 777); err != nil {
		//os.Exit(111)
		fmt.Println(err.Error())
	}
}

// Capitalize 字符首字母大写
func Capitalize(str string) string {
	var upperStr string
	vv := []rune(str) // 后文有介绍
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			if vv[i] >= 97 && vv[i] <= 122 {
				vv[i] -= 32 // string的码表相差32位
				upperStr += string(vv[i])
			} else {
				return str
			}
		} else {
			upperStr += string(vv[i])
		}
	}
	return upperStr
}
