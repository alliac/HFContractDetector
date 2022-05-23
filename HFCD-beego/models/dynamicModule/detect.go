package dynamicModule

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var data []byte

func Main_2(ReportPath string,chaincode_path string) time.Duration {
	start := time.Now()
	WriteToReport(ReportPath, "---------------------------Dynamic Analysis---------------------------")
	src := "./models/dynamicModule/Corpus/"
	dst := "./models/dynamicModule/OffCorpus/" //已经读过的文件放到这里面
	//detect SC
	DetectSc(ReportPath,chaincode_path)
	count := 0
	for {
		fileCount := 0 //如果中途有新文件加入，并不能获取到，所以循环遍历直到没有新文件
		err := filepath.Walk(src, func(filePath string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}
			if f.IsDir() {
				return nil
			}
			fileByte2, _ := ioutil.ReadFile(filePath)
			str := string(fileByte2)
			data := []byte(str)
			res := Slice(data)
			test := Taint(res, ReportPath)
			if test == 1 {
				Fuzz(res)
			}

			err = os.Rename(filePath, dst+f.Name()) // 读过的文件就移动过去
			if err != nil {
				return err
			}
			//time.Sleep(2 * time.Second)
			fileCount += 1
			count++
			if count > 10 {
				return io.EOF
			}
			return nil
		})
		if err != nil {
			fmt.Printf("filepath.Walk() returned %v\n", err)
		}
		if fileCount == 0 || count > 10 {
			break
		}
	}
	//WriteToReport(ReportPath, "---------------------------Dynamic Detection finished!---------------------------")
	var cost time.Duration
	cost = time.Since(start)
	return cost
}
func ExeSysCommand(cmdStr string) string {
	pwdBytes := exec.Command("sh", "-c", "pwd")
	pwd, _ := pwdBytes.Output()
	npwd := strings.Replace(string(pwd), "\n", "", -1)
	cmdStr = "cd " + npwd + ";" + cmdStr
	fmt.Println(cmdStr)
	cmd := exec.Command("sh", "-c", cmdStr)
	opBytes, _ := cmd.Output()
	return string(opBytes)
}

/*
str to [][][]byte
*/
func Slice(data []byte) [][][]byte {

	var ans [][][]byte
	var tempMap map[string]interface{}
	err := json.Unmarshal(data, &tempMap)
	if err != nil {
		panic(err)
	}
	for k, v := range tempMap {
		var temp_ans [][]byte
		temp_ans = [][]byte{[]byte(k)}
		tempMap2 := v.(map[string]interface{})
		for _, v2 := range tempMap2 {
			temp_ans = append(temp_ans, []byte(v2.(string)))
		}

		ans = append(ans, temp_ans)
	}
	//fmt.Println(ans)

	return ans
}

//打印三维数组
func print(data [][][]byte) {
	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data[i]); j++ {
			s := string(data[i][j])
			fmt.Print(s, " ")
		}
		fmt.Println()
	}
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
	defer f.Close()
}
