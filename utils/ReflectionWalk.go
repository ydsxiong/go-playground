package utils

import (
	"fmt"
	"reflect"
)

func walk(x interface{}, fn func(input string)) {

	val := getValue(x)

	numOfFields := 0
	var fieldAt func(int) reflect.Value

	switch val.Kind() {
	case reflect.String:
		fn(val.String())
	case reflect.Struct:
		numOfFields = val.NumField()
		fieldAt = val.Field
	case reflect.Slice, reflect.Array:
		numOfFields = val.Len()
		fieldAt = val.Index
	case reflect.Map:
		for _, key := range val.MapKeys() {
			walk(val.MapIndex(key).Interface(), fn)
		}
		return
	}

	for i := 0; i < numOfFields; i++ {
		walk(fieldAt(i).Interface(), fn)
	}
}

func getValue(x interface{}) reflect.Value {

	val := reflect.ValueOf(x)

	if val.Kind() == reflect.Ptr {
		fmt.Println("val.Kind():", val.Kind())
		fmt.Println("val is:", val)
		fmt.Println("val string is:", val.String())
		fmt.Println("val.Elem().Kind():", val.Elem().Kind())
		fmt.Println("val.Elem() is:", val.Elem())
		fmt.Println("val.Elem() string is:", val.Elem().String())

		return val.Elem()
	}

	return val
}

//field := fieldAt(i)
//var fieldVal interface{}
//switch field.Kind() {
//case reflect.String:
//	fieldVal = field.String()
//default: //case reflect.Struct, reflect.Ptr, reflect.Slice, reflect.Array, reflect.Map:
//	fieldVal = field.Interface()
//}
// slice val.Index(i) returns value that can be casted to interface type, no problem there,
// whereas the struct val.Field(i) returns primitive value that can't be casted to interface, while non-primitive value returned can,
// because the field declared in lower case which means not exported, so can't be used by reflection

// switch val.Kind() {
// case reflect.String:
// 	fn(val.String())
// case reflect.Slice:
// 	for i := 0; i < val.Len(); i++ {
// 		walk(val.Index(i).Interface(), fn)
// 	}
// case reflect.Struct:
// 	for i := 0; i < val.NumField(); i++ {
// 		walk(val.Field(i).Interface(), fn)
// 	}
// }
