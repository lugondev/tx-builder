package utils

import (
	"reflect"
)

func ToPtr(v interface{}) interface{} {
	p := reflect.New(reflect.TypeOf(v))
	p.Elem().Set(reflect.ValueOf(v))
	return p.Interface()
}

func CopyPtr(source, destin interface{}) {
	x := reflect.ValueOf(source)
	if x.Kind() != reflect.Ptr {
		return
	}

	starX := x.Elem()
	y := reflect.New(starX.Type())
	starY := y.Elem()
	starY.Set(starX)
	reflect.ValueOf(destin).Elem().Set(y.Elem())
}
