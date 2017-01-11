package utils

import (
	"fmt"
	"reflect"
)

func ConvertParamType(v interface{}, targetType reflect.Type) (
	targetValue reflect.Value, ok bool) {
	defer func() {
		if re := recover(); re != nil {
			ok = false
			fmt.Println(re)
		}
	}()

	ok = true

	if targetType.Kind() == reflect.Interface ||
		targetType.Kind() == reflect.TypeOf(v).Kind() {

		targetValue = reflect.ValueOf(v)

	} else if reflect.TypeOf(v).Kind() == reflect.Float64 {
		f := v.(float64)
		switch targetType.Kind() {
		case reflect.Int:
			targetValue = reflect.ValueOf(int(f))
		case reflect.Uint8:
			targetValue = reflect.ValueOf(uint8(f))
		case reflect.Uint16:
			targetValue = reflect.ValueOf(uint16(f))
		case reflect.Uint32:
			targetValue = reflect.ValueOf(uint32(f))
		case reflect.Uint64:
			targetValue = reflect.ValueOf(uint64(f))
		case reflect.Int8:
			targetValue = reflect.ValueOf(int8(f))
		case reflect.Int16:
			targetValue = reflect.ValueOf(int16(f))
		case reflect.Int32:
			targetValue = reflect.ValueOf(int32(f))
		case reflect.Int64:
			targetValue = reflect.ValueOf(int64(f))
		case reflect.Float32:
			targetValue = reflect.ValueOf(float32(f))
		default:
			ok = false
		}
	} else {
		ok = false
	}

	return
}
