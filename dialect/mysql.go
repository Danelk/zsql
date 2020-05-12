package dialect

import (
	"fmt"
	"reflect"
	"time"
)

type mysql struct {}

var _ Dialect = (*mysql)(nil)

func init()  {
	RegistDialect("mysql", &mysql{})
}

func (ms *mysql) DataTypeOf(typ reflect.Value) string {
	switch typ.Kind() {
	case reflect.Bool:
		return "bool"
	case reflect.Int, reflect.Int32:
		return "integer"
	case reflect.Int8:
		return "tinyint"
	case reflect.Int16:
		return "smallint"
	case reflect.Int64:
		return "bigint"
	case reflect.Uint,  reflect.Uint32:
		return "integer unsigned"
	case reflect.Uint8:
		return "tinyint unsigned"
	case reflect.Uint16:
		return "smallint unsigned"
	case reflect.Uint64:
		return "bigint unsigned"
	case reflect.Float32, reflect.Float64:
		return "float"
	case reflect.String:
		return "varchar(256)"
	case reflect.Struct:
		if _, ok := typ.Interface().(time.Time); ok {
			return "datetime"
		}
	}
	panic(fmt.Sprintf("invalid sql type %s (%s)", typ.Type().Name(), typ.Kind()))
}

func (ms *mysql) TableExistSQL(tableName string) (string, []interface{}) {
	args := []interface{}{tableName}
	return "SHOW TABLES ?",args
}
