package main

import (
	"context"
	"reflect"
	"sync"
)

type MyContext struct {
	c context.Context
	m *sync.Map
}

func NewMyContext() *MyContext {
	return &MyContext{
		c: context.Background(),
		m: &sync.Map{},
	}
}

func WithValue(parent *MyContext, key, val any) *MyContext {
	if parent == nil {
		panic("cannot create context from nil parent")
	}
	if key == nil {
		panic("nil key")
	}
	if !reflect.TypeOf(key).Comparable() {
		panic("key is not comparable")
	}
	parent.m.Store(key, val)
	return &MyContext{c: context.WithValue(parent.c, key, val), m: parent.m}
}

func (mc *MyContext) Value(key any) any {
	val, _ := mc.m.Load(key)
	return val

}

func main() {
	c := NewMyContext()
	c1 := WithValue(c, 2, 2)
	c2 := WithValue(c1, 3, 2)
	c3 := WithValue(c2, 4, 2)
	c4 := WithValue(c3, 5, 2)

	println(c.Value(1))
	println(c.Value(2))
	println(c.Value(3))
	println(c.Value(4))
	println(c.Value(5))

	println("--------------------")

	println(c1.Value(1))
	println(c1.Value(2))
	println(c1.Value(3))
	println(c1.Value(4))
	println(c1.Value(5))

	println("--------------------")

	println(c2.Value(1))
	println(c2.Value(2))
	println(c2.Value(3))
	println(c2.Value(4))
	println(c2.Value(5))

	println("--------------------")

	println(c3.Value(1))
	println(c3.Value(2))
	println(c3.Value(3))
	println(c3.Value(4))
	println(c3.Value(5))

	println("--------------------")

	println(c4.Value(1))
	println(c4.Value(2))
	println(c4.Value(3))
	println(c4.Value(4))
	println(c4.Value(5))

	println("--------------------")
}
