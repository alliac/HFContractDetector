package ast_graph

import (
	_ "crypto/rand"
	_ "encoding/json"
	"fmt"
	"io/ioutil"
	_ "io/ioutil"
	_ "math/rand"
	_ "net/http"
	_ "os"
	_ "os/exec"
	_ "time"
)

//全局变量
var Qint int
var Qbool bool

type Test struct {
	sint int
}

func (t *Test) Uses() {
	fmt.Println("hello word")
}
func main1() {
	var Jint int
	Jint = 2
	_ = 2                      //未处理错误
	ioutil.ReadFile("main.go") // UNCHECKED
	print("1")                 //内置函数
	fmt.Printf(string(Jint))
	//结构体全局变量
	var test Test
	test.Uses()
	var Jmap = make(map[int32]string)
	Jmap[0]="test"
	//程序并发协程+channel
	var messages chan string = make(chan string)
	go func(message string) {
		messages <- message // 存消息
	}("Ping!")

	fmt.Println(<-messages) // 取消息
}
