package datasource

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func Exec(result *[]Result, sql string, start string, end string, count int) {
	db, err := sqlx.Open("mysql", "root:bnsbdlmysql@tcp(172.16.77.154:33306)/analy?charset=utf8&parseTime=True")
	if err != nil {
		fmt.Println("数据库链接错误", err)
	}
	err = db.Select(result, sql, start, end, count)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

}
