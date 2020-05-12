package dialect

import "reflect"

var dialectsMap = map[string]Dialect{}

type Dialect interface {
	DataTypeOf(typ reflect.Value) string //用于将 Go 语言的类型转换为数据库的数据类型
	TableExistSQL(tableName string) (string, []interface{}) //返回某个表是否存在的 SQL 语句，参数是表名(table)
}

//注册 dialect 实例
func RegistDialect(name string, dialect Dialect)  {
	dialectsMap[name] = dialect
}

//获取 dialect 实例
func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectsMap[name]
	return
}
