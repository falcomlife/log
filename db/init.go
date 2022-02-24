package db

import "github.com/astaxie/beego/orm"
import _ "github.com/go-sql-driver/mysql"

var Dborm orm.Ormer = nil

func Orm() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	//orm.RegisterDataBase("default", "mysql", "root:bnsbdlmysql@tcp(172.16.77.158:33306)/analy?charset=utf8")
	orm.RegisterDataBase("default", "mysql", "root:1234+ABcd@tcp(172.18.57.42:3306)/analy?charset=utf8")
	orm.RegisterModel(new(Result))
	Dborm = orm.NewOrm()
}
