# Golang 结构体私有字段修改

为了实现数据的封装和信息隐藏,提高代码的健壮性和安全性,在 `Go` 语言中,结构体(`struct`)中的字段如果是私有的,只能在定义该结构体的同一个包内访问. 但是在某些情况下,可能需要在外部包中访问或修改结构体的私有字段. 一般情况下可以使用 `Go` 语言提供的反射(`reflect`)机制来实现这一功能. 

但是这只能实现访问，如果进行修改，试图通过反射设置这些私有字段的值会 `panic`. 更别说通过反射设置一些变量或者字段的值的时候，也会 `panic` 导致报错 `panic: reflect: reflect.Value.Set using unaddressable value`. 

那么就主要由下面三个问题构成：

1. 如何通过 `hack` 的方式访问外部结构体的私有字段?
    
2. 如何通过 `hack` 的方式设置外部结构体的私有字段?
    
3. 如何通过 `hack` 的方式设置 `unaddressable` 的值?
    
这里先介绍通过反射设置值遇到的 `unaddressable` 的困境. 

## 通过反射设置一个变量的值
------------

如果你使用过反射设置值的变量，你可能熟悉下面的代码，而且这个代码工作正常：

```go
package main

import (
	"fmt"
	"reflect"
)

func main() {
	var x = 47

	v := reflect.ValueOf(&x).Elem()
	fmt.Printf("原始值: %d, CanSet: %v\n", v.Int(), v.CanSet()) // 47, false
	v.Set(reflect.ValueOf(50))
}
```

注意这里传入给 `reflect.ValueOf` 的是 `x` 的指针 `&x`, 所以这个 Value 值是 `addresable` 的，我们可以进行赋值. 

如果把 `&x` 替换成 `x`, 我们再尝试运行就会 `panic`:

```bash
Original value: 47, CanSet: false
panic: reflect: reflect.Value.Set using unaddressable value

goroutine 1 [running]:
reflect.flag.mustBeAssignableSlow(0xc0001001a0?)
        /usr/local/go/src/reflect/value.go:272 +0x74
reflect.flag.mustBeAssignable(...)
        /usr/local/go/src/reflect/value.go:259
reflect.Value.Set({0x48b5c0?, 0x52b438?, 0x4a2494?}, {0x48b5c0?, 0x4bfb78?, 0x2?})
        /usr/local/go/src/reflect/value.go:2319 +0x65
main.main()
        /home/fangyuan/code/go/src/github.com/767829413/advanced-go/main.go:14 +0x1d4
exit status 2
```

这里暂时不做处理,后面会有一个 `hack` 的方式来解决.

接下来再看看访问私有字段的问题. 

## 访问外部包的结构体的私有字段
------------

这里先准备一个 `model` 包，定义了两个结构体来测试：

```go
package model

type Person struct {
	Name string
	age  int
}

func NewPerson(name string, age int) Person {
	return Person{
		Name: name,
		age:  age, // unexported field
	}
}

type Teacher struct {
	Name string
	Age  int // exported field
}

func NewTeacher(name string, age int) Teacher {
	return Teacher{
		Name: name,
		Age:  age,
	}
}
```

`Person` 的 `age` 字段是私有的，`Teacher` 的 `Age` 字段是公开的. 

在 `main` 函数中，不能访问 `Person` 的 `age` 字段：

```go
package main

import (
	"fmt"
	// "reflect"
	// "unsafe"

	"github.com/767829413/advanced-go/model"
)

func main() {
	p := model.NewPerson("Alice", 30)
	fmt.Printf("Person: %+v\n", p)

	// fmt.Println(p.age) // error: p.age undefined (cannot refer to unexported field or method age)

	t := model.NewTeacher("smallnest", 18)
	fmt.Printf("Teacher: %+v\n", t) // Teacher: {Name:Alice Age:30}
}
```

打印一下执行结果

```bash
Person: {Name:Alice age:30}
Teacher: {Name:smallnest Age:18}
```

可以通过反射的方式访问私有字段:

```go
package main

import (
	"fmt"
	"reflect"

	"github.com/767829413/advanced-go/model"
)

func main() {
	p := model.NewPerson("Alice", 30)
	age := reflect.ValueOf(p).FieldByName("age")
	fmt.Printf("原始值: %d, CanSet: %v\n", age.Int(), age.CanSet()) // 30, false
}
```

运行这个程序，可以看到我们获得了这个私有字段`age`的值:

```bash
原始值: 30, CanSet: false
```

这样我们就绕过了 Go 语言的访问限制，访问了私有字段. 

## 设置结构体的私有字段
------------

上面做到了获取私有属性的值, 但是如果尝试修改这个私有字段的值，会 panic:

```go
age.SetInt(100)
```

或者

```go
age.Set(reflect.ValueOf(100))
```

报错信息：

```bash
原始值: 30, CanSet: false
panic: reflect: reflect.Value.Set using value obtained using unexported field

goroutine 1 [running]:
reflect.flag.mustBeAssignableSlow(0xc0001261a0?)
        /usr/local/go/src/reflect/value.go:269 +0xb4
reflect.flag.mustBeAssignable(...)
        /usr/local/go/src/reflect/value.go:259
reflect.Value.Set({0x48d860?, 0xc000128010?, 0x4a4077?}, {0x48d860?, 0x4c2bf8?, 0x2?})
        /usr/local/go/src/reflect/value.go:2319 +0x65
```

实际上，`reflect.Value` 的 `Set` 方法会做一系列的检查，包括检查是否是`可寻址的（addressable）` 的，以及是否是`导出的字段（exported）`的字段:

```go
// Set assigns x to the value v.
// It panics if [Value.CanSet] returns false.
// As in Go, x's value must be assignable to v's type and
// must not be derived from an unexported field.
func (v Value) Set(x Value) {
	v.mustBeAssignable()
	x.mustBeExported() // do not let unexported x leak
......
```

`v.mustBeAssignable()` 检查是否是 `可寻址的（addressable）`的，而且是 `导出的字段（exported）` 的字段:

```go
// mustBeAssignable panics if f records that the value is not assignable,
// which is to say that either it was obtained using an unexported field
// or it is not addressable.
func (f flag) mustBeAssignable() {
	if f&flagRO != 0 || f&flagAddr == 0 {
		f.mustBeAssignableSlow()
	}
}

func (f flag) mustBeAssignableSlow() {
	if f == 0 {
		panic(&ValueError{valueMethodName(), Invalid})
	}
	// Assignable if addressable and not read-only.
	if f&flagRO != 0 {
		panic("reflect: " + valueMethodName() + " using value obtained using unexported field")
	}
	if f&flagAddr == 0 {
		panic("reflect: " + valueMethodName() + " using unaddressable value")
	}
}
```

`f&flagRO == 0` 代表是`导出的字段（exported）`，`f&flagAddr != 0` 代表是`可寻址的（addressable）`的,当这两个条件任意一个不满足时，就会报错. 

既然我们明白了它检查的原理，我们就可以通过 hack 的方式绕过这个检查，设置私有字段的值. 我们还是要使用`unsafe`代码. 

这里我们以标准库的`sync.Mutex`结构体为例， `sync.Mutex` 包含两个字段，这两个字段都是私有的：

```go
type Mutex struct {
	state int32
	sema  uint32
}
```

正常情况下只能通过 `Mutex.Lock` 和 `Mutex.Unlock` 来间接的修改这两个字段. 

现在通过 `hack` 的方式修改 `Mutex` 的 `state` 字段的值：

```go
package main

import (
	"fmt"
	"reflect"
	"sync"
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
	var mu sync.Mutex
	mu.Lock()

	field := reflect.ValueOf(&mu).Elem().FieldByName("state")
	fmt.Println("state", field.Int()) // 1

	flagField := reflect.ValueOf(&field).Elem().FieldByName("flag")
	flagPtr := (*uintptr)(unsafe.Pointer(flagField.UnsafeAddr())) // 2

	// 修改flag字段的值
	*flagPtr &= ^uintptr(flagRO) // 3

	// 修改 sync.Mutex 字段 state 的值
	field.Set(reflect.ValueOf(int32(0))) // 4

	mu.Lock() // 5

}

```

1. 通过反射获取了sync.Mutex结构体中的state字段的值，并打印出来. state字段用于表示互斥锁的状态，通常为0（未锁定）或1（锁定）. 

2. 通过反射获取了flag字段的地址，并将其转换为uintptr指针. flag字段用于存储一些标志位信息. 

3. 通过位运算清除了flag字段中的只读标志位（flagRO），使得该字段可以被修改. 这是通过对flag字段的指针进行位操作实现的. 

4. 通过反射将sync.Mutex的state字段的值设置为0. 这相当于手动解锁了互斥锁. 

5. 尝试再次锁定互斥锁. 由于之前通过反射和不安全操作解锁了互斥锁，这里不会导致死锁. 

这里使用的 `reflect.ValueOf(&mu).Elem().FieldByName("state")` 返回的 `Value`

```go

type Value struct {
	// typ_ holds the type of the value represented by a Value.
	// Access using the typ method to avoid escape of v.
	typ_ *abi.Type

	// Pointer-valued data or, if flagIndir is set, pointer to data.
	// Valid when either flagIndir is set or typ.pointers() is true.
	ptr unsafe.Pointer

	// flag holds metadata about the value.
	//
	// The lowest five bits give the Kind of the value, mirroring typ.Kind().
	//
	// The next set of bits are flag bits:
	//	- flagStickyRO: obtained via unexported not embedded field, so read-only
	//	- flagEmbedRO: obtained via unexported embedded field, so read-only
	//	- flagIndir: val holds a pointer to the data
	//	- flagAddr: v.CanAddr is true (implies flagIndir and ptr is non-nil)
	//	- flagMethod: v is a method value.
	// If ifaceIndir(typ), code can assume that flagIndir is set.
	//
	// The remaining 22+ bits give a method number for method values.
	// If flag.kind() != Func, code can assume that flagMethod is unset.
	flag

	// A method value represents a curried method invocation
	// like r.Read for some receiver r. The typ+val+flag bits describe
	// the receiver r, but the flag's Kind bits say Func (methods are
	// functions), and the top bits of the flag give the method number
	// in r's type's method table.
}
```

通过对 `flag` 进行操作来避开检查 `v.mustBeAssignable()` 的检查. 

这样就可以实现了修改私有字段的值了. 

## 使用 unexported 字段的 Value 设置公开字段
------------

看`reflect.Value.Set`的源码，可以看到它会检查参数的值是否 `导出的字段（exported）`，如果是，就会报错,下面就是一个例子：

```go
package main

import (
	"reflect"

	"github.com/767829413/advanced-go/model"
)

func main() {
	alice := model.NewPerson("Alice", 30)
	bob := model.NewTeacher("Bob", 40)

	bobAgent := reflect.ValueOf(&bob).Elem().FieldByName("Age")

	aliceAge := reflect.ValueOf(&alice).Elem().FieldByName("age")

	bobAgent.Set(aliceAge) // 4

}

```

注意 4 处，尝试把 `alice` 的私有字段 `age` 的值赋值给 `bob` 的公开字段 `Age`，这里会报错：

```bash
panic: reflect: reflect.Value.Set using value obtained using unexported field

goroutine 1 [running]:
reflect.flag.mustBeExportedSlow(0x46e101?)
        /usr/local/go/src/reflect/value.go:250 +0x65
reflect.flag.mustBeExported(...)
        /usr/local/go/src/reflect/value.go:241
reflect.Value.Set({0x474aa0?, 0xc000010040?, 0xc0000061c0?}, {0x474aa0?, 0xc000010028?, 0x0?})
        /usr/local/go/src/reflect/value.go:2320 +0x9a
```

原因 `alice` 的 `age` 值被识别为私有字段，它是不能用来赋值给公开字段的. 

有了上一节的经验，我们同样可以绕过这个检查，实现这个赋值：

```go
package main

import (
	"reflect"
	"unsafe"

	"github.com/767829413/advanced-go/model"
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
	alice := model.NewPerson("Alice", 30)
	bob := model.NewTeacher("Bob", 40)

	bobAgent := reflect.ValueOf(&bob).Elem().FieldByName("Age")

	aliceAge := reflect.ValueOf(&alice).Elem().FieldByName("age")
	// 修改flag字段的值
	flagField := reflect.ValueOf(&aliceAge).Elem().FieldByName("flag")
	flagPtr := (*uintptr)(unsafe.Pointer(flagField.UnsafeAddr()))
	*flagPtr &= ^uintptr(flagRO) // 5

	bobAgent.Set(reflect.ValueOf(50))
	bobAgent.Set(aliceAge) // 6
}
```

5 处修改了 `aliceAge` 的 `flag` 字段，去掉了 `flagRO` 标志位，这样就不会报错了,6 处我们成功的把 `alice` 的私有字段 `age` 的值赋值给 `bob` 的公开字段 `Age`. 

这样就可以实现了使用私有字段的值给其他 Value 值进行赋值了. 

## 给 unaddressable 的值设置值
------------

回到最初的问题，我们尝试给一个 `不可寻址的（unaddressable）` 的值设置值，会报错. 

结合上面的 hack 手段，我们也可以绕过限制，给 `不可寻址的（unaddressable）` 的值设置值：

```go
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

```

执行后会得到:

```bash
原始值: 47, CanSet: false
CanSet: true
修改后的值: 50
```

运行这个程序，不会报错，可以看到我们成功的给 unaddressable 的值设置了新的值. 

## 回顾
------------

我们通过修改`Value`值的 flag 标志位，可以绕过`reflect`的检查，实现了访问私有字段、设置私有字段的值、用私有字段设置值，以及给 unaddressable 的值设置值. 

这些都是`unsafe`的方式，一般情况下不鼓励进行这样的 hack 操作，但是这种技术也不是完全没有用户，如果你正在写一个 debugger，用户在断点出可能想修改某些值，或者你在写深拷贝的库，或者编写某种 ORM 库，或者你就像突破限制，访问第三方不愿意公开的字段，你有可能会采用这种非常规的技术. 
