package main

import (
	"myorm"
	"testing"
)

//type Users struct {
//	Id	 int	`myorm:"PRIMARY KEY"`
//	Name string `myorm:"NOT NULL"`
//	Age  int8
//}

func TestSession_CreateTable(t *testing.T) {
	engine, _ := myorm.NewEngine("mysql", "root:123456@tcp(localhost:3306)/test?charset=utf8")
	s := engine.NewSession().Model(&Users{})
	_ = s.DropTable()
	_ = s.CreateTable()
	if !s.HasTable() {
		t.Fatal("Failed to create table User")
	}

}
