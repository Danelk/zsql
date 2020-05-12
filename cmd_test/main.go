package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"myorm"
	_ "myorm/log"
)

func main()  {
	engine, _ := myorm.NewEngine("mysql", "root:123456@tcp(localhost:3306)/test?charset=utf8")
	defer engine.Close()
	s := engine.NewSession()
	//_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	//result, _ := s.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec()
	//fmt.Println(result)
	row:= s.Raw("select name from A limit 1").QueryRow()
	var name string
	err := row.Scan(&name)
	if err==nil {
		fmt.Println(name)
	}else{
		fmt.Println(err)
	}
}
