package myorm

import (
	"database/sql"
	"fmt"
	"myorm/dialect"
	"myorm/log"
	"myorm/session"
	"strings"
)

type Engine struct {
	db *sql.DB
	dialect dialect.Dialect
}

//定义事务接口
type TxFunc func(*session.Session) (interface{}, error)

//创建引擎实例
func NewEngine(driver, source string) (e *Engine, err error) {
	//连接数据库，返回 *sql.DB
	db, err := sql.Open(driver, source)
	if err != nil{
		log.Error(err)
		return
	}
	//调用 db.Ping()，检查数据库是否能够正常连接
	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}
	dial, ok := dialect.GetDialect(driver)
	if !ok {
		log.Errorf("dialect %s Not Found", driver)
	}
	e = &Engine{db: db, dialect: dial}
	log.Info("Connect database success!")
	return
}

//关闭数据库连接引擎
func (e *Engine) Close()  {
	if err := e.db.Close(); err!= nil {
		log.Error(err)
	}
	log.Info("Close database success!")
}

//根据引擎实例创建会话
func (e *Engine) NewSession() *session.Session {
	return session.New(e.db, e.dialect)
}

//实现事务接口的方法
func (e *Engine) Transaction(txFunc TxFunc) (result interface{}, err error) {
	s := e.NewSession()
	if err = s.Begin(); err != nil{
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil{
			_ = s.Rollback()
			panic(p)
		}else if err != nil{
			_ = s.Rollback()
		}else{
			err = s.Commit()
		}
	}()
	return txFunc(s)
}

func difference(a []string, b []string) (diff []string) {
	mapB := make(map[string]bool)
	for _, v := range b{
		mapB[v] = true
	}
	for _, v := range a{
		if _, ok := mapB[v]; !ok {
			diff = append(diff, v)
		}
	}
	return 
}

func (e *Engine) Migrate(value interface{}) error {
	_, err := e.Transaction(func(s *session.Session) (result interface{}, err error) {
		if !s.Model(value).HasTable(){
			log.Info("table %s is not exist", s.RefTable().Name)
			return nil,s.CreateTable()
		}
		table := s.RefTable()
		rows, _ := s.Raw(fmt.Sprintf("SELECT * FROM %s LIMIT 1", table.Name)).QueryRows()
		columns, _ := rows.Columns()
		addcols := difference(table.FieldNames, columns)
		delcols := difference(columns, table.FieldNames)
		log.Info("add %v delete %v", addcols, delcols)
		for _, col := range addcols {
			f := table.GetField(col)
			sqlStr := fmt.Sprintf("ALTER TEBLE %s ADD COLUMN %s %s", table.Name, f.Name, f.Type)
			if _, err = s.Raw(sqlStr).Exec(); err != nil{
				return
			}
		}
		if len(delcols) == 0 {
			return
		}
		tmp := "tmp" + table.Name
		fieldStr := strings.Join(table.FieldNames, ",")
		s.Raw(fmt.Sprintf("CREATE TABLE %s AS SELECT %s from %s;", tmp, fieldStr, table.Name))
		s.Raw(fmt.Sprintf("DROP TABLE %s;", table.Name))
		s.Raw(fmt.Sprintf("ALTER TABLE %s RENAME TO %s;", tmp, table.Name))
		_, err = s.Exec()
		return
	})
	return err
}