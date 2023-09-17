package errors

import (
	"reflect"
)

func IsNil(a any) bool {
	if a == nil {
		return true
	}

	v := reflect.ValueOf(a)
	switch v.Kind() {
	case reflect.Invalid:
		return true
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return v.IsNil()
	}
	return false
}

func NotNil(a any) bool {
	return !IsNil(a)
}
