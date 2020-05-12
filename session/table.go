package session

import (
	"fmt"
	"myorm/log"
	"myorm/schema"
	"reflect"
	"strings"
)

//给 refTable 赋值。解析操作是比较耗时的，因此将解析的结果保存在成员变量 refTable 中，即使 Model() 被调用多次，如果传入的结构体名称不发生变化，则不会更新 refTable 的值
func (s *Session) Model(value interface{}) *Session {
	if s.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(s.refTable.Model) {
		s.refTable = schema.Parse(value, s.dialect)
	}
	return  s
}

//返回 refTable 的值，如果 refTable 未被赋值，则打印错误日志
func (s *Session) RefTable() *schema.Schema {
	if s.refTable == nil {
		log.Error("Model is not set")
	}
	return s.refTable
}

//创建表
func (s *Session) CreateTable() error {
	table := s.refTable
	var columns []string
	pk  := ""
	sql :=""
	for _, field := range table.Fields{
		if field.Tag == "PRIMARY KEY" {
			columns = append(columns, fmt.Sprintf("%s %s", field.Name, field.Type))
			pk = fmt.Sprintf("%s (%s)", field.Tag, field.Name)
		}else{
			columns = append(columns, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
		}
	}
	desc := strings.Join(columns, ",")
	if pk != ""{
		sql = desc+","+pk
	}else {
		sql = desc
	}

	log.Info(fmt.Sprintf("CREATE TABLE %s (%s);", table.Name, sql))
	_, err := s.Raw(fmt.Sprintf("CREATE TABLE %s (%s);", table.Name, sql)).Exec()
	return  err
}

//删除表
func (s *Session) DropTable() error {
	_, err := s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s;", s.refTable.Name)).Exec()
	return  err
}

//表是否存在
func (s *Session) HasTable() bool {
	sql, values := s.dialect.TableExistSQL(s.RefTable().Name)
	row := s.Raw(sql, values...).QueryRow()
	var tmp string
	_ = row.Scan(&tmp)
	return tmp == s.RefTable().Name
}