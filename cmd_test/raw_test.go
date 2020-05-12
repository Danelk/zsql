package main

import (
	"myorm"
	"myorm/session"
)

type Users struct {
		Name string `myorm:"NOT NULL"`
		Age  int8
}

func NewSession() *session.Session {
	engine, _ := myorm.NewEngine("mysql", "root:123456@tcp(localhost:3306)/test?charset=utf8")
	return engine.NewSession()
}