package main

import (
	"myorm"
	"reflect"
	"testing"
)

func OpenDb(t *testing.T) *myorm.Engine {
	t.Helper()
	e, err := myorm.NewEngine("mysql", "root:123456@tcp(localhost:3306)/test?charset=utf8")
	if err != nil{
		t.Fatal("mysql connect faild")
	}
	return e
}

type Users1 struct {
	Name  	string	`myorm: "PRIMARY KEY"`
	Age 	int
}

//func TestEngine_Transaction(t *testing.T)  {
//	t.Run("rollback", func(t *testing.T) {
//		transactionRollback(t)
//	})
//	t.Run("commit", func(t *testing.T) {
//		transactionCommit(t)
//	})
//}
//
//func transactionCommit(t *testing.T)  {
//	engine := OpenDb(t)
//	defer engine.Close()
//	s := engine.NewSession()
//	_ = s.Model(&Users1{}).DropTable()
//	_, err := engine.Transaction(func(s *session.Session) (rs interface{}, err error) {
//		_ = s.Model(&Users1{}).CreateTable()
//		_, err = s.Insert(&Users1{"Tran", 12})
//		return
//	})
//	u := &Users1{}
//	_ = s.First(u)
//	if err != nil || u.Name != "Tran" {
//		t.Fatal("failed to commit")
//	}
//}
//
//func transactionRollback(t *testing.T)  {
//	engine := OpenDb(t)
//	defer engine.Close()
//	s := engine.NewSession()
//	_ = s.Model(&Users1{}).DropTable()
//	_, err := engine.Transaction(func(s *session.Session) (rs interface{}, err error) {
//		_ = s.Model(&Users1{}).CreateTable()
//		_, err = s.Insert(&Users1{"Tran", 12})
//		return
//	})
//	if err != nil || s.HasTable(){
//		t.Fatal("failed to rollback")
//	}
//}
func TestEngine_Migrate(t *testing.T) {
	engine := OpenDb(t)
	defer engine.Close()
	s := engine.NewSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS Users1;").Exec()
	_, _ = s.Raw("CREATE TABLE Users1(Name text PRIMARY KEY, XXX integer);").Exec()
	_, _ = s.Raw("INSERT INTO Users1(`Name`) values (?), (?)", "Tom", "Sam").Exec()
	engine.Migrate(&Users1{})

	rows, _ := s.Raw("SELECT * FROM Users1").QueryRows()
	columns, _ := rows.Columns()
	if !reflect.DeepEqual(columns, []string{"Name", "Age"}) {
		t.Fatal("Failed to migrate table Users1, got columns", columns)
	}
}