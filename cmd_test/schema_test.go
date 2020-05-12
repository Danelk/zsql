package main

import (
	"myorm/dialect"
	schema2 "myorm/schema"
	"testing"
)

//type User struct {
//	Name string `myorm:"PRIMARY KEY"`
//	Age  int
//}

var TestDial, _ = dialect.GetDialect("mysql")

func TestParse(t *testing.T) {
	schema := schema2.Parse(&Users{}, TestDial)
	if schema.Name != "User" || len(schema.Fields) != 2 {
		t.Fatal("failed to parse User struct")
	}
	if schema.GetField("Name").Tag != "PRIMARY KEY" {
		t.Fatal("failed to parse primary key")
	}
}