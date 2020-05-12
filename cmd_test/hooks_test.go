package main

import (
	"myorm/log"
	"myorm/session"
	"testing"
)

type Account struct {
	ID 			int		`myorm:"PRIMARY KEY"`
	Password 	string
}

func (a *Account)BeforInsert(s *session.Session) error {
	log.Info("befor insert ...")
	a.ID +=1000
	return nil
}

func (a *Account)AfterQuery(s *session.Session) error {
	log.Info("after select ...")
	a.Password = "******"
	return nil
}

func TestSession_CallMethod(t *testing.T) {
	s := NewSession().Model(&Account{})
	_ = s.DropTable()
	_ = s.CreateTable()
	_, _ = s.Insert(&Account{1, "123456"}, &Account{2, "qwerty"})

	u := &Account{}

	err := s.First(u)
	if err != nil || u.ID != 1001 || u.Password != "******" {
		t.Fatal("Failed to call hooks after query, got", u)
	}
}