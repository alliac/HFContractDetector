package bugs

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

const (
	db_name = "hfdata"
	db_host = "127.0.0.1"
	db_user = "root"
	db_pass = "root"
	db_port = 3306
)

var DB *sql.DB

func Init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", db_user, db_pass, db_host, db_port, db_name)
	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("sql.Open error", err)
		return
	}

	//超时时间
	DB.SetConnMaxLifetime(100 * time.Second)
	// 最大连接数
	DB.SetMaxOpenConns(100)
	// 设置闲置的连接数
	DB.SetMaxIdleConns(16)
	if err := DB.Ping(); err != nil {
		log.Fatal("DB.Ping = ", err)
		return
	}

	//fmt.Println("connnect success")
}

func RecordPos(name string, pos int, end int) {
	Init()
	//开启事务
	tx, err := DB.Begin()
	if err != nil {
		fmt.Println("tx fail")
		return
	}
	//准备sql语句
	stmt, err := tx.Prepare("INSERT INTO func(name, pos,end) VALUES (?,?,?)")
	if err != nil {
		fmt.Println("Prepare fail")
		return
	}
	//将参数传递到sql语句中并且执行
	_, err = stmt.Exec(name, pos, end)
	if err != nil {
		//fmt.Println("Exec fail",err)
		return
	}
	//将事务提交
	tx.Commit()
}
func isExist(pos int, end int) bool {
	Init()
	defer DB.Close()
	rows, err := DB.Query("SELECT * FROM func")
	if err != nil {
		fmt.Printf("query faied, error:[%v]", err.Error())
		return false
	}
	flag := false
	for rows.Next() {
		var (
			name  string
			dbPos int
			dbEnd int
		)
		rows.Scan(&name, &dbPos, &dbEnd)
		//fmt.Println(pos)
		//fmt.Println(end)
		if pos >= dbPos && end <= dbEnd {
			flag = true
		}
	}
	rows.Close()
	if flag{
		return false
	}
	return true
}
func ClearTable(){
	Init()
	defer DB.Close()
	_, err := DB.Exec("delete from func where name like '%'")
	if err != nil {
		fmt.Printf("query faied, error:[%v]", err.Error())
	}
}