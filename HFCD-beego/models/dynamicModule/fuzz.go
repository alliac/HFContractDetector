package dynamicModule

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"
)

/*
语料变异
*/
func Fuzz(data [][][]byte) {
	var res [][][]byte
	//mutator
	for i := 0; i < len(data); i++ {
		var tmp [][]byte
		if len(data[i]) > 0 {
			tmp = append(tmp, data[i][0])
		}
		for j := 1; j < len(data[i]); j++ {
			tmp = append(tmp, MutatorCorpus(data[i][j]))
		}
		res = append(res, tmp)
	}
	//[][][]byte to str
	str := ByteToStr(res)
	save(str)

}
/*糊器
 */
func MutatorCorpus(data []byte) []byte {
	m:=newMutator()
	return m.Ｍutate(data)
}
/*
三维数组转string
 */
func ByteToStr(data [][][]byte) string {
	ans := map[string]interface{}{}
	for i := 0; i < len(data); i++ {
		temp_ans := map[string]string{}
		for j := 1; j < len(data[i]); j++ {
			val := string(data[i][j])
			key := "arg" + strconv.Itoa(j)
			temp_ans[key] = val
		}
		key := string(data[i][0])
		ans[key] = temp_ans
	}
	dataByte, _ := json.Marshal(ans)
	dataString := string(dataByte)
	return dataString
}
func save(str string) {
	//保存到文件
	timestamp := time.Now().Unix()
	f, err := os.Create("./models/dynamicModule/Corpus/corpus-"+strconv.Itoa(int(timestamp))+".txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	_, err2 := f.WriteString(str)
	if err2 != nil {
		log.Fatal(err2)
	}
}