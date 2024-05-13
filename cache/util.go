package cache

import (
	"encoding/json"
	"log"
	"reflect"
)

func getValue(value any) any {
	if value == nil {
		value = ""
		return value
	}
	switch reflect.TypeOf(value).Kind() {
	case reflect.Struct, reflect.Map, reflect.Slice, reflect.Array, reflect.Ptr, reflect.Interface:
		res, err := json.Marshal(value)
		if err != nil {
			//Channel, complex, and function values 这三个值
			log.Fatal(" redisvalue marshal value failed %v,it should not occur", err)
		} else {
			value = res
		}
	case reflect.Func, reflect.Chan:
		value = ""
	default:
	}
	return value
}
