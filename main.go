package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type flag uintptr

const (
	flagKindWidth        = 5 // there are 27 kinds
	flagKindMask    flag = 1<<flagKindWidth - 1
	flagStickyRO    flag = 1 << 5
	flagEmbedRO     flag = 1 << 6
	flagIndir       flag = 1 << 7
	flagAddr        flag = 1 << 8
	flagMethod      flag = 1 << 9
	flagMethodShift      = 10
	flagRO          flag = flagStickyRO | flagEmbedRO
)

func main() {
	var x = 47

	v := reflect.ValueOf(x)
	fmt.Printf("原始值: %d, CanSet: %v\n", v.Int(), v.CanSet()) // 47, false
	// v.Set(reflect.ValueOf(50))

	flagField := reflect.ValueOf(&v).Elem().FieldByName("flag")
	flagPtr := (*uintptr)(unsafe.Pointer(flagField.UnsafeAddr()))

	// 修改flag字段的值
	*flagPtr |= uintptr(flagAddr)          // 设置可寻址标志位
	fmt.Printf("CanSet: %v\n", v.CanSet()) // true
	v.SetInt(50)
	fmt.Printf("修改后的值: %d\n", v.Int()) // 50
}
