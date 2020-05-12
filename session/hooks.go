package session

import (
	"myorm/log"
	"reflect"
)

const (
	BeforQuery 	= "BeforQuery"
	AfterQuery 	= "AfterQuery"
	BeforUpdate = "BeforUpdate"
	AfterUpdate = "AfterUpdate"
	BeforDelete = "BeforDelete"
	AfterDelete = "AfterDelete"
	BeforInsert = "BeforInsert"
	AfterInsert = "AfterInsert"
)

func (s *Session) CallMethod(method string, value interface{})  {
	//s.RefTable.Model 或 value 即当前会话正在操作的对象,MethodByName 方法反射得到该对象的方法
	fm := reflect.ValueOf(s.refTable.Model).MethodByName(method)
	if value != nil{
		fm = reflect.ValueOf(value).MethodByName(method)
	}
	param := []reflect.Value{reflect.ValueOf(s)}
	if fm.IsValid(){
		if v:=fm.Call(param); len(v) > 0{
			if err, ok := v[0].Interface().(error); ok{
				log.Error(err)
			}
		}
	}
	return
}
