package controllers

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func Init(path string, chaincode string, version string) {
	//create test_report
	fr, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer fr.Close()
}
func writeToFile(msg string, URL string) {
	if err := ioutil.WriteFile(URL, []byte(msg), 777); err != nil {
		//os.Exit(111)
		fmt.Println(err.Error())
	}
}
func Statistics(ReportPath string, cost string) {
	res := ExeSysCommand("python3 ./models/dynamicModule/Statistics.py")
	fmt.Println(res)
	WriteToReport(ReportPath, "---------------------------Test Result---------------------------")
	WriteToReport(ReportPath, res)
	WriteToReport(ReportPath, cost)
}
func ExeSysCommand(cmdStr string) string {
	pwdBytes := exec.Command("sh", "-c", "pwd")
	pwd, _ := pwdBytes.Output()
	npwd := strings.Replace(string(pwd), "\n", "", -1)
	cmdStr = "cd " + npwd + ";" + cmdStr
	cmd := exec.Command("sh", "-c", cmdStr)
	opBytes, _ := cmd.Output()
	return string(opBytes)
}

//print test report
func WriteToReport(ReportPath string, s string) {
	content_s := Getfile(ReportPath)
	if !strings.Contains(content_s, s) {
		//f, _ := os.OpenFile(ReportPath, syscall.O_APPEND, 0666)
		f, _ := os.OpenFile(ReportPath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		_, err := f.WriteString(s + "\n")
		defer f.Close()
		if err != nil {
			fmt.Println("090")
			log.Fatal(err)
		}
		if err = f.Close(); err != nil {
			log.Fatal(err)
		}
	}
}
func Getfile(ReportPath string) string {
	data, err := ioutil.ReadFile(ReportPath)
	if err != nil {
		return fmt.Sprint(err)
	}
	return string(data)
}
