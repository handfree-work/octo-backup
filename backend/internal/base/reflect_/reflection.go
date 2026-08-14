package reflect_

import (
	"reflect"
)

func HasAttr(target interface{}, fieldName string) bool {
	obj := reflect.ValueOf(target)
	elem := obj.Elem()
	field := elem.FieldByName(fieldName)
	return field.IsValid()
}

func GetAttrValue(target interface{}, fieldName string, kind reflect.Kind) reflect.Value {
	obj := reflect.ValueOf(target)
	elem := obj.Elem()
	field := elem.FieldByName(fieldName)
	if field.Kind() == reflect.Ptr {
		return field.Elem()
	}
	if field.IsValid() {
		// 在使用 Value.Int() 之前，检查 Value 是否为 int 类型
		return field
	}
	//返回0值
	if kind == reflect.Int {
		return reflect.ValueOf(0)
	} else if kind == reflect.String {
		return reflect.ValueOf("")
	} else if kind == reflect.Bool {
		return reflect.ValueOf(false)
	} else if kind == reflect.Float64 {
		return reflect.ValueOf(0.0)
	} else if kind == reflect.Float32 {
		return reflect.ValueOf(0.0)
	} else if kind == reflect.Int64 {
		return reflect.ValueOf(int64(0))
	}
	return reflect.ValueOf(nil)
}

func SetAttrValue(target interface{}, fieldName string, value interface{}) {
	obj := reflect.ValueOf(target)
	elem := obj.Elem().Elem()
	field := elem.FieldByName(fieldName)
	fieldElem := field.Elem()
	valueRef := reflect.ValueOf(value)
	if valueRef.Kind() == reflect.Int64 {
		fieldElem.SetInt(valueRef.Int())
	} else if valueRef.Kind() == reflect.Int {
		fieldElem.SetInt(valueRef.Int())
	} else if valueRef.Kind() == reflect.String {
		fieldElem.SetString(valueRef.String())
	} else if valueRef.Kind() == reflect.Bool {
		fieldElem.SetBool(valueRef.Bool())
	} else if valueRef.Kind() == reflect.Float64 {
		fieldElem.SetFloat(valueRef.Float())
	} else if valueRef.Kind() == reflect.Float32 {
		fieldElem.SetFloat(valueRef.Float())
	} else {
		field.Set(valueRef)
	}
}
