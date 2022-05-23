package bugs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var gopath = os.Getenv("GOROOT")
var rootPath = gopath + "/src"
var packageExist = make(map[string]bool)

func getDirlist() {
	rootPath = "/usr/local/go/src"
	err := filepath.Walk(rootPath, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if !f.IsDir() {
			return nil
		}
		if len(path) <= len(rootPath) {
			return nil
		}
		path = path[len(rootPath)+1:]
		path = strings.Replace(path, "\\", "/", -1) //-1全部替换
		//fmt.Println(path)
		packageExist[path] = true
		return nil
	})
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}
}
