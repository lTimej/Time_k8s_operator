package utils

import (
	"reflect"
	"strconv"
)

func IntToString(n interface{}) string {
	t := reflect.TypeOf(n)
	switch t.Kind() {
	case reflect.Uint:
		val := n.(uint)
		return strconv.FormatInt(int64(val), 10)
	case reflect.Uint8:
		val := n.(uint8)
		return strconv.FormatInt(int64(val), 10)
	case reflect.Uint16:
		val := n.(uint16)
		return strconv.FormatInt(int64(val), 10)
	case reflect.Uint32:
		val := n.(uint32)
		return strconv.FormatInt(int64(val), 10)
	case reflect.Uint64:
		val := n.(uint64)
		return strconv.FormatInt(int64(val), 10)
	case reflect.Int:
		val := n.(int)
		return strconv.FormatInt(int64(val), 10)
	case reflect.Int8:
		val := n.(int8)
		return strconv.FormatInt(int64(val), 10)
	case reflect.Int16:
		val := n.(int16)
		return strconv.FormatInt(int64(val), 10)
	case reflect.Int32:
		val := n.(int32)
		return strconv.FormatInt(int64(val), 10)
	default:
		return ""
	}
}
