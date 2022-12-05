package main

import (
	"context"
)

type MyContext struct {
	cur   context.Context
	child []context.Context
}

func NewMyContext(k, v any) *MyContext {
	return &MyContext{
		cur: context.WithValue(context.Background(), k, v),
	}
}

func WithValue(parent *MyContext, key, val interface{}) *MyContext {
	nmc := &MyContext{
		cur: context.WithValue(parent.cur, key, val),
	}
	parent.child = append(parent.child, nmc.cur)
	return nmc
}

func (m *MyContext) Value(key any) any {
	res := m.cur.Value(key)
	if res == nil {
		for _, child := range m.child {
			v := child.Value(key)
			if v != nil {
				return v
			}
		}
	}
	return res

}

func main() {
	c := NewMyContext(1, 1)
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
