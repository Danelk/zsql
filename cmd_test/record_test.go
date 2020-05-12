package main

import (
	"myorm/session"
	"testing"
)


var (
	user1 = Users{"Tom", 18}
	user2 = Users{"Sam", 25}
	user3 = Users{"Jack", 25}
)

func testRecordInit(t *testing.T) *session.Session {
	t.Helper()
	//engine, _ := myorm.NewEngine("mysql", "root:123456@tcp(localhost:3306)/test?charset=utf8")
	s := NewSession().Model(Users{})
	//err1 := s.DropTable()
	//err2 := s.CreateTable()
	//_, err3 := s.Insert(user1, user2)
	//if err1 != nil || err2 != nil || err3 != nil {
	//	t.Fatal("failed init test records")
	//}
	return s
}

//func TestSession_Insert(t *testing.T) {
//	s := testRecordInit(t)
//	affected, err := s.Insert(user3)
//	if err != nil || affected != 1 {
//		t.Fatal("failed to create record")
//	}
//}
//
//func TestSession_Find(t *testing.T) {
//	s := testRecordInit(t)
//	var users []Users
//	if err := s.Find(&users); err != nil || len(users) != 2 {
//		t.Fatal("failed to query all")
//	}
//}

func TestSession_Limit(t *testing.T) {
	s := testRecordInit(t)
	var users []Users
	err := s.Limit(1).Find(&users)
	if err != nil || len(users) != 1 {
		t.Fatal("failed to query with limit condition")
	}
}

func TestSession_Update(t *testing.T) {
	s := testRecordInit(t)
	affected, _ := s.Where("Name = ?", "Tom").Update("Age", 30)
	u := &Users{}
	_ = s.OrderBy("Age DESC").First(u)

	if affected != 1 || u.Age != 30 {
		t.Fatal("failed to update")
	}
}

func TestSession_DeleteAndCount(t *testing.T) {
	s := testRecordInit(t)
	affected, _ := s.Where("Name = ?", "Tom").Delete()
	count, _ := s.Count()

	if affected != 1 || count != 1 {
		t.Fatal("failed to delete or count")
	}
}
