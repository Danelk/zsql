package session

import (
	"database/sql"
	"myorm/clause"
	"myorm/dialect"
	"myorm/log"
	"myorm/schema"
	"strings"
)
//包含三个成员变量，第一个是 db *sql.DB，即使用 sql.Open() 方法连接数据库成功之后返回的指针
//第二个和第三个成员变量用来拼接 SQL 语句和 SQL 语句中占位符的对应值。用户调用 Raw() 方法即可改变这两个变量的值
type Session struct {
	Db  		*sql.DB
	tx 			*sql.Tx	
	dialect 	dialect.Dialect
	refTable 	*schema.Schema
	clause      clause.Clause
	sql 		strings.Builder
	sqlVars 	[]interface{}
}

type CommonDB interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result,error)
}

var _ CommonDB = (*sql.DB)(nil)
var _ CommonDB = (*sql.Tx)(nil)


//新建连接会话
func New(db *sql.DB, d dialect.Dialect) *Session {
	return &Session{
		Db: db,
		dialect: d,
	}
}

//清空连接会话 (s *Session).sql 和 (s *Session).sqlVars 两个变量 使 session可复用
func (s *Session) Clear()  {
	s.sql.Reset()
	s.sqlVars = nil
	s.clause = clause.Clause{}
}

//返回数据库连接指针
//tx *sql.Tx tx不为空时，则使用 tx 执行 SQL 语句，否则使用 db 执行 SQL 语句
func (s *Session) DB() CommonDB {
	if s.tx != nil{
		return s.tx
	}
	return s.Db
}

//拼接sql语句与变量
func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, values...)
	return s
}

//封装 Exec()、Query() 和 QueryRow() 三个原生方法
//执行
func (s *Session) Exec() (result sql.Result, err error)  {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if result, err = s.DB().Exec(s.sql.String(), s.sqlVars...); err != nil{
		log.Error(err)
	}
	return result, err
}
//单条查询
func (s *Session) QueryRow() *sql.Row{
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	return s.DB().QueryRow(s.sql.String(), s.sqlVars...)
}
//返回多行
func (s *Session) QueryRows() (rows *sql.Rows, err error)  {
	defer  s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if rows, err = s.DB().Query(s.sql.String(), s.sqlVars...); err != nil{
		log.Error(err)
	}
	return rows, err
}