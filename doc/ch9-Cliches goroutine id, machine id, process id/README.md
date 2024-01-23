# 老生常谈: goroutine id, machine id, process id

## 主要方向

1. 使用汇编实现getg函数,获取g对象
2. Go语言调用汇编函数getg
3. 伪造一个数据结构,通过指针强转
4. 通过指针偏移获取某个字段的值
5. 通过名称(字符串)获得一个reflect.Type
6. goroutine 具备父子关系,可以获取(Go版本1.21及以上)
7. 遍历所有的 goroutine
8. 获取GPM模型中的M的id, P的id
9. per-cpu (per-p)并发模式,性能远远好于Mutex和atomic

## GPM 模型和背景知识

### GPM 模型基础概念介绍

G, P, M 是 Go 调度器的三个核心组件,分别代表 goroutine、machine、和 processor这就是 Go 运行时的 G-P-M 模型的核心概念,它们分别代表：

1. G (goroutine): 
	- G 代表 Goroutine,是 Go 语言并发的基本单元
	- 每个 Goroutine 都有自己的栈,用于保存局部变量等信息

2. M (machine): 
	- M 代表 Machine,是一个抽象的逻辑执行环境
	- 每个 M 都关联一个 P,并负责调度 Goroutines 执行
	- 在系统调用被阻塞的情况下,可能会创建很多的 M,并且这些 M 不会自动被回收

3. P (processor): 
	- P 代表 Processor,是实际执行 Goroutine 的执行器
	- 每个 P 都关联一个线程（OS 线程或内核线程）
	- P 的数量通常由 GOMAXPROCS 决定,它限制了同时执行用户级 Go 代码的操作系统线程的数量

### GPM 模型的由来与变化

网上已经有很多对 Go 的 GPM 模型的介绍和分析了,最早看到的是 2013 年 Morsing 的 [The Go scheduler](https://morsmachine.dk/go-scheduler).

有些人会把 machine 直接和线程划等号,甚至还会把 Processor 理解成 CPU 处理器,这是不对的.G、P、M 是 Go 调度器的概念,Go 运行时的代码中专门为它们定义了相应的数据结构,Processor 和 CPU 处理器半毛关系都没有,它是 Go 调度器管理 P 和 M 的处理器.甚至都不能把 P 和 CPU 处理器对应上,因为多个 P 运行的线程可能只在 CPU 多个个处理器上的一个核上运行.

这里有以前梳理的一份 GPM 模型和操作系统线程调度相关文章: [调度组件与调度循环](https://github.com/767829413/advanced-go/blob/main/doc/ch1-How%20golang%20works/README.md#%E8%B0%83%E5%BA%A6%E7%BB%84%E4%BB%B6%E4%B8%8E%E8%B0%83%E5%BA%A6%E5%BE%AA%E7%8E%AF)

虽然 GPM 模型已经广而告之了,但是不妨简单的再回顾一下：

在 Golang 1.1 版本之前调度器中还没有 P 组件,当时的 Go 调度器性能有点差,Dmitry Vyukov 大佬针对此调度器中存在的问题进行了重新设计,引入 P 组件来解决当前面临的问题（[Scalable Go Scheduler Design Doc]( https://docs.google.com/document/d/1TTj4T2JO42uD5ID9e89oa0sLKhJYD0Y_kqxDv3I3XMw/)）,开始在 Go 1.1 版本中引入了 P 这个对象.

 G-P-M 沿用至今,有些小变动,主要解决：

1. 全局互斥锁(Sched.Lock)和中心化的状态管理.全局一把大锁管理 goroutine 创建、调度和完成,竞争带来的性能严重下降.

2. goroutine 的切换毫无规律. M 之间经常切换 goroutine 执行,不能利用 M-G 的亲和性.

3. Per-M 内存缓存（M.mcache）：每个 M 都有自己的 mcache, 在一些 M 被阻塞的时候,实际会有大量的 M 产生,每个 M 的 mcache 都占用一部分的内存,存在内存浪费

4. 线程阻塞与唤醒频繁, 增加了很多的开销.

新增加了 P 的角色后,很多的功能交给 P 来完成,而 P 的数量是在 runtime 启动时设置的,默认等于 cpu 的逻辑核数.

可以在程序启动时通过环境变量  GOMAXPROCS  或者  runtime.GOMAXPROCS()  方法来修改默认配置,也可以在程序运行的时候修改 P 的数量,但是一般不会这么做.

### 扩展: 一些奇怪的问题

在 IO 密集型场景下, 可以适当调高 P 的数量, 可以有技巧性的提高程序的性能.[badger](https://github.com/dgraph-io/badger) 的作者 Manish Rai Jain 就遇到了这样一个[问题](https://groups.google.com/forum/#!topic/golang-nuts/jPb_h3TvlKE/discussion).

Manish Rai Jain 写了一段 [Go 代码]( https://github.com/dgraph-io/badger-bench/blob/master/randread/main.go), 用来测试 Go 的读写 SSD 的性能, 看看是否能达到 Linux 的 I/O 测试工具 [fio](https://linux.die.net/man/1/fio)的读写性能. 使用 fio, 可以在 AWS(Amazon i3.large instance with NVMe SSD)达到 100K IOPS,  但是使用这个 Go 程序, 怎么也无法接近这个 IOPS.

如果大家都使用小于 CPU 的核数, Go 和 Fio 的吞吐率接近, 但是如果把 goroutine 设置为大于 CPU 的核数, Fio 使用也设置使用查过CPU核数, 这种情况下 fio 性能提升明显, 可以达到最大的 IOPS, 但是 Go 程序却没有显著变化.通过将GOMAXPROCS设置更大的数(64/128, 数倍 CPU 核数), Go 程序可以取得几乎和 Fio 一样的吞吐率.

这是一个怪异的现象, 在 2017 年的 GopherCon 大会上, Manish Rai Jain 碰到了 Russ Cox, 就咨询了这个问题：

如果 File::Read 阻塞 goroutines, Go 运行时会生成新的线程(M), 因此在一个长时间运行的任务中, Go 应该创建了大量的线程, 因此随机读取吞吐量应该随时间增加并稳定到最大可能的值.但是, 我没有在基准测试中看到预期的情况.

Russ Cox 解释说 GOMAXPROCS 就像一个多路复用器, 所以 GOMAXPROCS 会是 I/O 读取的瓶颈.P 数量少的吞吐就小.

```text
GOMAXPROCS 在某种程度上充当复用器.根据文档的说法, “GOMAXPROCS 变量限制了可以同时执行用户级 Go 代码的操作系统线程的数量.” 这基本上意味着, 在切换到某个操作系统线程之前, 所有读取必须首先仅通过 GOMAXPROCS 数量的 goroutine 运行（实际上并非真正的切换, 但在概念上是这样）

这为吞吐量引入了一个瓶颈
```

后来 Manish Rai Jain 就在他的高性能的 K/V 数据库调大了这个[参数](https://github.com/dgraph-io/dgraph/commit/30237a1429debab73eff38fea2f724914ca38b77).

## goroutine id 的获取和应用

Go 运行时中定义了 goroutine 模型,它的数据结构如下[g](https://github.com/golang/go/blob/master/src/runtime/runtime2.go#L422)

不过可以先看两个字段: goid(goroutine的id) , m(当前关联的 M)

```go
type g struct {
......
	m         *m      // current m; offset known to arm liblink
......
	goid         uint64
......
}
```

g()这个数据结构并没有公开(exported), goid 也没有公开,正常一个 goroutine 在运行的时候,是没有办法拿到它自己的 id 的,这是 Go 开发团队故意为之,在 [Go FAQ](https://go.dev/doc/faq) 中专门介绍了这个[no_goroutine_id](https://go.dev/doc/faq#no_goroutine_id),以及有些讨论：[No.9](https://groups.google.com/forum/#!topic/golang-nuts/Nt0hVV_nqHE) [No.10](https://groups.google.com/forum/#!topic/golang-nuts/0HGyCOrhuuI) [No.11](http://stackoverflow.com/questions/19115273/looking-for-a-call-or-thread-id-to-use-for-logging).有些人会使用 goroutine 实现 goroutine-local storage,比如[jtolio/gls](jtolio/gls: https://github.com/jtolio/gls)、[tylerstillwater/gls](https://github.com/tylerstillwater/gls),也有人会在 debug log 时进行并发分析.

虽然 Go 并不直接提供获取 goroutine 的方法,避免大家误用,但是有一个还算可以的方法,可以解析出 goroutine id,那就是通过分析 goroutine 的调用栈,来得到 goroutine id:

```go
func GoID() int {
    var buf [64]byte
    n := runtime.Stack(buf[:], false)
    idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
    id, err := strconv.Atoi(idField)
    if err != nil {
        panic(fmt.Sprintf("cannot get goroutine id: %v", err))
    }
    return id
}
```

就像我们 panic 常常打印出的 goroutine 信息一样,goroutine 字符串后面紧跟的就是 goroutine id:

```bash
panic: 1233

goroutine 1 [running]:
main.main()
        /root/go/src/github.com/767829413/advanced-go/main.go:22 +0x78
exit status 2
```

这是一个兜底的方法,效率不高,需要读取当前 goroutine 的栈信息,还需要字符串解析.

更高效的获取 goroutine id 的方法,如果是在前几年,我们一般会使用 [petermattis/goid](https://github.com/petermattis/goid),提供了一个巧妙的方法来获取 goroutine id,而且还支持低版本的 Go.

它的巧妙之处在于使用汇编,实现了getg()方法:

```c
//go:build (386 || amd64 || amd64p32 || arm || arm64) && gc && go1.5
// +build 386 amd64 amd64p32 arm arm64
// +build gc
// +build go1.5

#include "textflag.h"

// func getg() *g
// 定义一个名为getg的文本段,表示这是一个函数.
// NOSPLIT表示该函数不需要栈分裂.
// $0-8表示该函数不接受参数,但返回一个 8 字节的值
TEXT ·getg(SB),NOSPLIT,$0-8
// 处理了不同的架构,比如 386 的架构
#ifdef GOARCH_386
    // 将 TLS（线程本地存储）中的值加载到寄存器 AX 中
    MOVL (TLS), AX
	// 将 AX 中的值存储到返回值ret的位置
	// 这个位置是相对于函数帧指针（FP）的偏移量
    MOVL AX, ret+0(FP)
#endif
#ifdef GOARCH_amd64
    MOVQ (TLS), AX
    MOVQ AX, ret+0(FP)
#endif
#ifdef GOARCH_arm
    MOVW g, ret+0(FP)
#endif
#ifdef GOARCH_arm64
    MOVD g, ret+0(FP)
#endif
    RET
```

这样我们通过返回值就可以得到 g 的数据.

`这里就是使用汇编实现getg函数,获取g对象`

我们在一个 go 文件中使用这个汇编定义的函数即可:

```go
//go:build (386 || amd64 || amd64p32 || arm || arm64) && gc && go1.5
// +build 386 amd64 amd64p32 arm arm64
// +build gc
// +build go1.5

package goid

// 在汇编文件中定义的函数
func getg() *g

func Get() int64 {
    return getg().goid
}
```

`Go语言调用汇编函数getg`

在 Go 语言中调用使用汇编编写的函数,它只定义了getg签名,就可以直接使用了.

可能会有疑问:

**g数据结构不是 Go 运行时中没有公开么？**

确实,不过可以自己模仿 go 运行时定义一个自己的g,甚至比它更简单,因为只需要goid,所以我们的g只需要定义到goid即可,而且前面的字段保持 offset 一致即可, `伪造一个数据结构,通过指针强转`:

```go
//go:build gc && go1.9
// +build gc,go1.9

package goid

type stack struct {
    lo uintptr
    hi uintptr
}

type gobuf struct {
    sp   uintptr
    pc   uintptr
    g    uintptr
    ctxt uintptr
    ret  uintptr
    lr   uintptr
    bp   uintptr
}

type g struct {
    stack       stack
    stackguard0 uintptr
    stackguard1 uintptr

    _panic       uintptr
    _defer       uintptr
    m            uintptr
    sched        gobuf
    syscallsp    uintptr
    syscallpc    uintptr
    stktopsp     uintptr
    param        uintptr
    atomicstatus uint32
    stackLock    uint32
    goid         int64 // 定义到这里
}
```

这样,就可以实现快速获取到 goroutine id 了.

以前,如果需要获取 goroutine id 的时候,都会想到[petermattis/goid](https://github.com/petermattis/goid).不过这个库是需要跟随 Go 版本发布的,不然就有可能运行时中g的数据结构有变化,那么导致必须做相应的更新,否则goid字段的偏移量变了,就获取不到正确的 id 了.

但是通过 hacker 方式获取 goroutine id 的方法还是在有很多人在琢磨.2020 年 7 月 Laevus Dexter 提供了另外一种方法： [Go Playground](https://go.dev/play/p/CSOp9wyzydP),非常的巧妙,Dan Kortschak 在 github 上为它建立了一个库：[kortschak/goroutine](https://github.com/kortschak/goroutine)

```go
package goroutine

import (
    "reflect"
    "unsafe"
)

// ID returns the runtime ID of the calling goroutine.
func ID() int64 {
    return idOf(getg(), goidoff)
}

func idOf(g unsafe.Pointer, off uintptr) int64 {
    return *(*int64)(add(g, off))
}

//go:nosplit
func getg() unsafe.Pointer {
    return *(*unsafe.Pointer)(add(getm(), curgoff))
}

//go:linkname add runtime.add
//go:nosplit
func add(p unsafe.Pointer, x uintptr) unsafe.Pointer

//go:linkname getm runtime.getm
//go:nosplit
func getm() unsafe.Pointer

var (
    curgoff = offset("*runtime.m", "curg")
    goidoff = offset("*runtime.g", "goid")
)

// offset returns the offset into typ for the given field.
func offset(typ, field string) uintptr {
    rt := toType(typesByString(typ)[0])
    f, _ := rt.Elem().FieldByName(field)
    return f.Offset
}

//go:linkname typesByString reflect.typesByString
func typesByString(s string) []unsafe.Pointer

//go:linkname toType reflect.toType
func toType(t unsafe.Pointer) reflect.Type
```

通过go:linkname,可以链接其它库中的未公开的方法.后面有机会再说.这里只需知道//go:linkname typesByString reflect.typesByString会把typesByString链接到reflect.typesByString这个 reflect 库未公开的 typesByString 上.

这段代码最下面三行其实是获得一个类型的某个字段在这个类型中的偏移量,这样我们通过移动指针就可以得到一个对象的某个字段的值.这是 `通过指针偏移获取某个字段的值`.

这两行就是获取m.curg字段的偏移量和g.goid的偏移量：

```go
   curgoff = offset("*runtime.m", "curg")
   goidoff = offset("*runtime.g", "goid")
```

接下来就要获取g这个对象的地址了,通过getm得到当前的 M,然后通过它的curg字段得到当前的 goroutine 的地址

```go
//go:nosplit
func getg() unsafe.Pointer {
    return *(*unsafe.Pointer)(add(getm(), curgoff))
}

//go:linkname add runtime.add
//go:nosplit
func add(p unsafe.Pointer, x uintptr) unsafe.Pointer

//go:linkname getm runtime.getm
//go:nosplit
func getm() unsafe.Pointer
```

这里getm链接到 runtime 的getm方法,获取当前的 M.add方法链接到 runtime 的add方法,用来移动指针.

add(getm(), curgoff)其实通过移动偏移量,得到curg的指针.

```go
func ID() int64 {
    return idOf(getg(), goidoff)
}

func idOf(g unsafe.Pointer, off uintptr) int64 {
    return *(*int64)(add(g, off))
}
```

idOf就是获取g对象的某个字段,这也是通过指针挪动偏移量获取的.

而ID就是获取goid这个字段的值.

把整体串起来,就是完整的获取 goroutine id 的方法了.通过一系列的黑科技的操作,我们不需要写汇编代码,不需要伪造g的数据结构,就能高效的获取 goroutine 的 id 了.

如果再发掘一下,也就是通过一个类型的名称,我们就能得到一个对象的类型(reflect.Type),得到这个反射类型我们就可以凭空创建一个对象了.通过名称创建一个对象在 Java 中很容易,在 Go 语言中确实很难得,`通过名称(字符串)获得一个reflect.Type`.

## goroutine 的父子关系

对于Go语言来说, goroutine 创建完子 goroutine 之后, 它们之间就没有什么瓜葛了.

1. 父子两个 goroutine 都可能在不同的 P 中运行, 任何一个的退出都不会影响到另一个的运行.

2. 父 goroutine 要想控制子 goroutine 的退出,一般会使用 context 或者 channel, 通过 context 或者 channel 产生信号告诉子 goroutine. 子 goroutine 得到信号后自己决定要退出或者忽略,或者做点啥.

**goroutine 这个父子关系我们在运行时能够获取到吗？**

在 Go 1.21 之前,是没有好办法的,但是在 Go 1.21 中,g结构体增加了parentGoid这样一个字段,可以通过这个字段得到一个 goroutine 的父 goroutine id.

`goroutine 具备父子关系,可以获取(Go版本1.21及以上)`

```go
func ParentID() int64 {
    return idOf(getg(), parentGoidoff)
}

var parentGoidoff = offset("*runtime.g", "parentGoid")
```

基于先前的介绍, 你已经知道了getg的实现了.既然得到了g,通过指针移动就能够获取到 parentGoid 字段的值了.

我们可以写一个程序测试这个功能

```go
package main

import (
    "fmt"
    "time"

    "github.com/kortschak/goroutine"
)

func foo(depth int) {
    fmt.Printf("goroutine #%d: depth=%d, parent: #%d\n", goroutine.ID(), depth, goroutine.ParentID())

    if depth == 0 {
        return
    }

    go foo(depth - 1)
}
func main() {
    depth := 5
    go foo(depth)

    time.Sleep(1 * time.Second)
}
```

这个程序中通过递归的方式创建 6 个 goroutine, 每个 goroutine 会打印出它自己的 goroutine id, 深度, 以及父 goroutine id.

运行这个程序如下：

```bash
root@GreenCloud:advanced-go# go run "/root/go/src/github.com/767829413/advanced-go/main.go"
goroutine #6: depth=5, parent: #1
goroutine #7: depth=4, parent: #6
goroutine #8: depth=3, parent: #7
goroutine #9: depth=2, parent: #8
goroutine #10: depth=1, parent: #9
goroutine #11: depth=0, parent: #10
```

甚至可以遍历程序运行时的所有的 goroutine 的父子关系, 这是通过All方法实现的：

```go
// Link is a goroutine parent-child relationship.
type Link struct {
    Parent, Child int64
}

// All returns all the known goroutine parent-child relationships.
func All() []Link {
    var s []Link
    forEachG(func(g unsafe.Pointer) {
        s = append(s, Link{Parent: idOf(g, parentGoidoff), Child: idOf(g, goidoff)})
    })
    return s
}

//go:linkname forEachG runtime.forEachG
//go:nosplit
func forEachG(fn func(gp unsafe.Pointer))
```

`可以遍历所有的 goroutine, 通过链接 runtime 的forEachG方法实现.`

写一个程序测试:

```go
func main() {
    links := goroutine.All()
    sort.Slice(links, func(i, j int) bool {
        return links[i].Child < links[j].Child
    })

    for _, link := range links {
        fmt.Printf("%d -> %d\n", link.Parent, link.Child)
    }
}
```

这里没有创建额外的 goroutine,就只有运行时创建的 goroutine.

也可以自己创建一些 goroutine, sleep 一下不退出, 看看结果.

运行上面的程序, 可以看到 g0 创建了 g1, g1 创建了 2,3,4,5 goroutine：

```bash
root@GreenCloud:advanced-go# go run "/root/go/src/github.com/767829413/advanced-go/main.go"
0 -> 1
1 -> 2
1 -> 3
1 -> 4
1 -> 5
```

`获取GPM模型中的M的id`

通过上面的介绍, 也很容易实现获取M的 id 了, 通过getm就能得到当前 M 的指针, 获取它的id字段就好了.

```go
// MID returns the "M" ID of the calling goroutine.
func MID() int {
    return int(idOf(getm(), midoff))
}

var (
    midoff = offset("*runtime.m", "id")
)
```

## 获取 processor id

`获取GPM模型中的P的id`

前面已经介绍了获取 goroutine 和 machine 的 id, 现在到了获取 P 的 id 了. 

获取 P 的 id 我们使用另外一个更简单的方式: 

```go
//go:linkname runtime_procPin runtime.procPin
func runtime_procPin() int

//go:linkname runtime_procUnpin runtime.procUnpin
func runtime_procUnpin() int

// PID returns the "P" ID of the calling goroutine.
func PID() int {
    pid := runtime_procPin()
    runtime_procUnpin()
    return pid
}
```

runtime_procPin功能是将当前的 goroutine 和 P 固定住, 不会被别的 P 抢去, 它返回 P 的 id. 

获取到 pid 之后, 再调用runtime_procUnpin解除固定. 

注意这里的pid不是操作系统的进程的 id,而是 Go 运行时中的 GPM 模型中的 P 的 id. 

我们能够利用这个技巧优化性能. 因为每个 P 同时只能运行一个 goroutine, 我们可以实现per-p类型的数据结构, 不必使用 Mutex 等互斥锁, 这样就减少了数据的并发, 可以提高程序的性能. 

sync.Pool的实现就是利用这个技巧:

[go/src/sync/pool.go at b25f5558c69140deb652337afaab5c1186cd0ff1 · golang/go (github.com):](https://github.com/golang/go/blob/b25f5558c69140deb652337afaab5c1186cd0ff1/src/sync/pool.go#L207)

## per-cpu 模型和性能比较

如果将上面的功能提炼, 就可以实现一个基于 P 的 Shard 数据结构, 可以把它提炼成了一个类型Shard

```go
package util

import (
	"runtime"

	"github.com/767829413/advanced-go/goroutine"
	"golang.org/x/sys/cpu"
)

// Shard类型
type Shard[T any] struct {
	values []value[T]
}

// NewShard创建一个新的Shard.
func NewShard[T any]() *Shard[T] {
	n := runtime.GOMAXPROCS(0)

	return &Shard[T]{
		values: make([]value[T], n),
	}
}

// 避免伪共享
type value[T any] struct {
	_ cpu.CacheLinePad
	v T
	_ cpu.CacheLinePad
}

// 得到当前P的值
func (s *Shard[T]) Get() *T {
	if len(s.values) == 0 {
		panic("sync: Sharded is empty and has not been initialized")
	}

	return &s.values[int(goroutine.PID())%len(s.values)].v
}

// 遍历所有的值
func (s *Shard[T]) Range(f func(*T)) {
	if len(s.values) == 0 {
		panic("sync: Sharded is empty and has not been initialized")
	}

	for i := range s.values {
		f(&s.values[i].v)
	}
}
```

比如实现一个计数器, 可以拿它和 Mutex、atomic 实现的计数器, 在并发情况下进行比较

```go
package util

import (
	"sync"
	"sync/atomic"
	"testing"
)

func BenchmarkShardCounter(b *testing.B) {
	counter := NewShard[atomic.Int64]()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			counter.Get().Add(1)
		}
	})
}

func BenchmarkMutexCounter(b *testing.B) {
	var counter int64
	var mu sync.Mutex
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu.Lock()
			counter += 1
			mu.Unlock()
		}
	})
}

func BenchmarkAtomicCounter(b *testing.B) {
	var counter atomic.Int64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			counter.Add(1)
		}
	})
}
```

执行下benchmark,结果如下：

```bash
go test -bench Counter -benchmem .

goos: linux
goarch: amd64
pkg: github.com/767829413/advanced-go/util
cpu: Intel(R) Xeon(R) Platinum
BenchmarkShardCounter-2         125707922                9.389 ns/op           0 B/op          0 allocs/op
BenchmarkMutexCounter-2         52898214                29.94 ns/op            0 B/op          0 allocs/op
BenchmarkAtomicCounter-2        69747235                15.02 ns/op            0 B/op          0 allocs/op
PASS
ok      github.com/767829413/advanced-go/util   5.810s
```

观察上来看, Shard 的性能好于 Mutex 和 atomic.

Go 的 issue 中有多个提议希望 Go 官方库增加这个功能, 比如[Issue #8281](https://github.com/golang/go/issues/8281)、[Issue #18802](https://github.com/golang/go/issues/18802), 一些 Gopher 尝试实现, 比如 qiulaidongfeng 的 [Sharded](https://go.dev/cl/552515), 上面实现的 Shard 就是模仿这个实现的, 只不过通过黑魔法获取到了 pid,不需要改运行时.

还有一个非常好的库[cespare/percpu](https://github.com/cespare/percpu), 它的实现和当前实现的 Shard 差不多, 只不过当前实现的泛型的方式, 它还没有修改为泛型, 还在使用interface{}代表值.和现在这个测试差不多, 也可以看到它的测试性能比 Mutex 和 atomic要好点, 对于正在挖掘性能的 Gopher 来说值得关注.但是因为不支持泛型, 它的性能理论上来说比泛型的实现要差一些.