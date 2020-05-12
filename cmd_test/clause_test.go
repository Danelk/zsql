package main

import (
	clause2 "myorm/clause"
	"reflect"
	"testing"
)

func testSelect(t *testing.T) {
	var clause clause2.Clause
	clause.Set(clause2.LIMIT, 3)
	clause.Set(clause2.SELECT, "Users", []string{"*"})
	clause.Set(clause2.WHERE, "Name = ?", "Tom")
	clause.Set(clause2.ORDERBY, "Age ASC")
	sql, vars := clause.Build(clause2.SELECT, clause2.WHERE, clause2.ORDERBY, clause2.LIMIT)
	t.Log(sql, vars)
	if sql != "SELECT * FROM Users WHERE Name = ? ORDER BY Age ASC LIMIT ?" {
		t.Fatal("failed to build SQL")
	}
	if !reflect.DeepEqual(vars, []interface{}{"Tom", 3}) {
		t.Fatal("failed to build SQLVars")
	}
}

func TestClause_Build(t *testing.T) {
	t.Run("select", func(t *testing.T) {
		testSelect(t)
	})
}