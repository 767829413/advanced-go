# Golang程序是怎么跑起来的

`宏观认识Go的启动和执行流程`

## 理解可执⾏⽂件

`Go程序main.go的编译过程`

![Go程序main.go的编译过程.jpg](https://s2.loli.net/2023/05/24/fXphaZHLKe9PcCT.png)

```bash
[root@~]# go build -x main.go 
WORK=/tmp/go-build379780741
mkdir -p $WORK/b001/
cat >$WORK/b001/importcfg.link << 'EOF' # internal
packagefile command-line-arguments=/root/.cache/go-build/1d/1d5ef4eb9cec52ed9f5e784fce37b8d2bc956780a67e5bb01a8c5e0e7bd084ea-d
packagefile fmt=/usr/local/go/pkg/linux_amd64/fmt.a
packagefile runtime=/usr/local/go/pkg/linux_amd64/runtime.a
packagefile errors=/usr/local/go/pkg/linux_amd64/errors.a
packagefile internal/fmtsort=/usr/local/go/pkg/linux_amd64/internal/fmtsort.a
packagefile io=/usr/local/go/pkg/linux_amd64/io.a
packagefile math=/usr/local/go/pkg/linux_amd64/math.a
packagefile os=/usr/local/go/pkg/linux_amd64/os.a
packagefile reflect=/usr/local/go/pkg/linux_amd64/reflect.a
packagefile strconv=/usr/local/go/pkg/linux_amd64/strconv.a
packagefile sync=/usr/local/go/pkg/linux_amd64/sync.a
packagefile unicode/utf8=/usr/local/go/pkg/linux_amd64/unicode/utf8.a
packagefile internal/bytealg=/usr/local/go/pkg/linux_amd64/internal/bytealg.a
packagefile internal/cpu=/usr/local/go/pkg/linux_amd64/internal/cpu.a
packagefile runtime/internal/atomic=/usr/local/go/pkg/linux_amd64/runtime/internal/atomic.a
packagefile runtime/internal/math=/usr/local/go/pkg/linux_amd64/runtime/internal/math.a
packagefile runtime/internal/sys=/usr/local/go/pkg/linux_amd64/runtime/internal/sys.a
packagefile internal/reflectlite=/usr/local/go/pkg/linux_amd64/internal/reflectlite.a
packagefile sort=/usr/local/go/pkg/linux_amd64/sort.a
packagefile math/bits=/usr/local/go/pkg/linux_amd64/math/bits.a
packagefile internal/oserror=/usr/local/go/pkg/linux_amd64/internal/oserror.a
packagefile internal/poll=/usr/local/go/pkg/linux_amd64/internal/poll.a
packagefile internal/syscall/execenv=/usr/local/go/pkg/linux_amd64/internal/syscall/execenv.a
packagefile internal/syscall/unix=/usr/local/go/pkg/linux_amd64/internal/syscall/unix.a
packagefile internal/testlog=/usr/local/go/pkg/linux_amd64/internal/testlog.a
packagefile sync/atomic=/usr/local/go/pkg/linux_amd64/sync/atomic.a
packagefile syscall=/usr/local/go/pkg/linux_amd64/syscall.a
packagefile time=/usr/local/go/pkg/linux_amd64/time.a
packagefile unicode=/usr/local/go/pkg/linux_amd64/unicode.a
packagefile internal/race=/usr/local/go/pkg/linux_amd64/internal/race.a
EOF
mkdir -p $WORK/b001/exe/
cd .
/usr/local/go/pkg/tool/linux_amd64/link -o $WORK/b001/exe/a.out -importcfg $WORK/b001/importcfg.link -buildmode=exe -buildid=UuhzpN2gMByCeX217N_I/lfzY_Yru8Wk-v-Z75vDW/TUGfC5VLv-Y9aNaN--gY/UuhzpN2gMByCeX217N_I -extld=gcc /root/.cache/go-build/1d/1d5ef4eb9cec52ed9f5e784fce37b8d2bc956780a67e5bb01a8c5e0e7bd084ea-d
/usr/local/go/pkg/tool/linux_amd64/buildid -w $WORK/b001/exe/a.out # internal
cp $WORK/b001/exe/a.out main
rm -r $WORK/b001/
```

`可执行文件在不同的操作系统上的规范不一样`

| Linux | Windows | Mac OS X |
| --- | ---| --- |
| ELF | PE | Mach-O |

Linux 的可执⾏⽂件 ELF(Executable and Linkable Format) 为例，
ELF 由⼏部分构成：

* ELF header
* Section header
* Sections

参考资料: <https://github.com/corkami/pics/blob/28cb0226093ed57b348723bc473cea0162dad366/binary/elf101/elf101.pdf>

`操作系统执⾏可执⾏⽂件的步骤(以 linux 为例)：`

![Go程序main.go的编译过程 _1_.jpg](https://s2.loli.net/2023/05/24/63vOn8zZJiaSFjo.png)

`使⽤ readelf 找到 Go 进程的执⾏⼊⼝`

```bash
[root@~]# readelf -h ./main
ELF Header:
  Magic:   7f 45 4c 46 02 01 01 00 00 00 00 00 00 00 00 00 
  Class:                             ELF64
  Data:                              2's complement, little endian
  Version:                           1 (current)
  OS/ABI:                            UNIX - System V
  ABI Version:                       0
  Type:                              EXEC (Executable file)
  Machine:                           Advanced Micro Devices X86-64
  Version:                           0x1
  Entry point address:               0x45d990
  Start of program headers:          64 (bytes into file)
  Start of section headers:          456 (bytes into file)
  Flags:                             0x0
  Size of this header:               64 (bytes)
  Size of program headers:           56 (bytes)
  Number of program headers:         7
  Size of section headers:           64 (bytes)
  Number of section headers:         25
  Section header string table index: 3
[root@~]# dlv exec ./main
Type 'help' for list of commands.
(dlv) b *0x45d990
Breakpoint 1 (enabled) set at 0x45d990 for _rt0_amd64_linux() /usr/local/go/src/runtime/rt0_linux_amd64.s:8
(dlv) 
```

## Go 进程的启动与初始化

`计算机执⾏程序`

![Go程序main.go的编译过程 _2_.jpg](https://s2.loli.net/2023/05/24/bPVZE2NJs7D4Q95.png)

```text
CPU ⽆法理解⽂本,只能执⾏⼀条⼀条的⼆进制机器码指令每次执⾏完⼀条指令,pc 寄存器就指向下⼀条继续执⾏
在 64 位平台上 pc 寄存器 = rip
```

`Go 语⾔是⼀⻔有 runtime 的语⾔`

参考资料: <https://www.techtarget.com/searchsoftwarequality/definition/runtime>

```text
What is runtime?
Runtime is a piece of code that implements portions of a programming language's execution model. In doing this, it allows the program to interact with the computing resources it needs to work. Runtimes are often integral parts of the programming language and don't need to be installed separately.

Runtime is also when a program is running. That is, when you start a program running in a computer, it is runtime for that program. In some programming languages, certain reusable programs or "routines" are built and packaged as a "runtime library." These routines can be linked to and used by any program when it is running.

Programmers sometimes distinguish between what gets embedded in a program when it is compiled and what gets embedded or used at runtime. The former is sometimes called compile time.
```

**可以认为 runtime 是为了实现额外的功能，⽽在程序运⾏时⾃动加载/运⾏的⼀些模块**

`Go 语⾔的 runtime 模块`

| 名称 | 介绍 |
| --- | --- |
| Scheduler | 调度器管理所有的 G，M，P，在后台执⾏调度循环 |
| Netpoll | ⽹络轮询负责管理⽹络 FD 相关的读写、就绪事件 |
| Memory Management | 当代码需要内存时，负责内存分配⼯作 |
| Garbage Collector | 当内存不再需要时，负责回收内存 |

**最核⼼是 Scheduler，负责串联所有的 runtime 流程**

`通过 entry point 找到 Go 进程的执⾏⼊⼝`

![Go程序main.go的编译过程 _3_.jpg](https://s2.loli.net/2023/05/24/OkQnUP6Ez3uRjMp.png)

## 调度组件与调度循环

`任务简述`

```text
每次写下：

	go func() {
		fmt.Println("hello world in goroutine")
	}()

都是向 runtime 提交了⼀个计算任务。
func() { xxxxx } ⾥包裹的代码就是这个计算任务的内容
```

`Go 的调度流程本质上是⼀个⽣产-消费流程`

![Go程序main.go的编译过程 _4_.jpg](https://s2.loli.net/2023/05/24/2Wr8YKiJSBHamU3.png)

`调度组件`

![组件调度.png](https://s2.loli.net/2023/05/24/tB1qyDG8eJ2NjZV.png)

`goroutine 的⽣产端`

![Go程序main.go的编译过程 _5_.jpg](https://s2.loli.net/2023/05/24/rlbmnkvTUjJLZFW.png)

演示动画: <https://www.figma.com/proto/gByIPDf4nRr6No4dNYjn3e/bootstrap?node-id=5004-6&starting-point-node-id=5004%3A6>

`goroutine 的消费端`

![Go程序main.go的编译过程 _6_.jpg](https://s2.loli.net/2023/05/24/A1Nmsbxh25QBcv7.png)

**Work stealing 就是说的 runqsteal -> runqgrab 这个流程**

**P.schedtick就是用来计算调度循环次数**

演示动画: <https://www.figma.com/proto/gByIPDf4nRr6No4dNYjn3e/bootstrap?node-id=5004-1299&starting-point-node-id=5004%3A1299>

`字段定义`

* G：goroutine，⼀个计算任务。由需要执⾏的代码和其上下⽂组成，上下⽂包括：当前代码位置，栈顶、栈底地址，状态等。

* M：machine，系统线程，执⾏实体，想要在 CPU 上执⾏代码，必须有线程，与 C 语⾔中的线程相同，通过系统调⽤ clone 来创建。

* P：processor，虚拟处理器，M 必须获得 P 才能执⾏代码，否则必须陷⼊休眠(后台监控线程除外)，你也可以将其理解为⼀种 token，有这个 token，才有在物理 CPU 核⼼上执⾏的权⼒

## 处理阻塞

`不会创建线程的阻塞(调度循环不会被阻塞,runtime 可以处理)`

```go
// 1.无缓存通道的接和收
var chSend = make(chan int)
ch <- 1

var chRecv = make(chan int)
<- ch

// 2. 用户态执行停止
time.Sleep(time.Hour)

// 3. 网络读写
var c net.Conn
var buf = make([]byte, 1024)

n,err = c.Read(buf)
n,err = c.Write(buf)

// 4. 执行 select 语句
var (
    ch1 = make(chan int)
    ch2 = make(chan int)
)

select {
    case <- ch1:
        println("ch1 is ready")
    case <- ch2:
        println("ch2 is ready")
}

// 5. 执行锁操作
var l sync.RWMutex
l.lock()
```

上述这些情况不会阻塞调度循环，⽽是会把 goroutine 挂起,所谓的挂起，其实让 g 先进某个数据结构，待 ready 后再继续执⾏.

**不会占⽤线程**

线程会进⼊ schedule，继续消费队列，执⾏其它的 g

1~4对应的等待结构:

![359.png](https://s2.loli.net/2023/05/24/W2SMPBTcbUxtvrR.png)

5对应的等待结构:

![lock3.png](https://s2.loli.net/2023/05/24/EpS7jYa2LWdIOB3.png)

`sudog 和 G的关系`

```text
// sudog represents a g in a wait list, such as for sending/receiving
// on a channel.
//
// sudog is necessary because the g ↔ synchronization object relation
// is many-to-many. A g can be on many wait lists, so there may be
// many sudogs for one g; and many gs may be waiting on the same
// synchronization object, so there may be many sudogs for one object.
//
// sudogs are allocated from a special pool. Use acquireSudog and
// releaseSudog to allocate and free them.
```

**⼀个 g 可能对应多个 sudog，⽐如⼀个 g 会同时 select 多个channel**

`runtime无法处理的阻塞(必定占用系统线程)`

```go
// 1. CGO
package main
/*
#include <stdio.h>
#include <stdlib.h>
#include <ubistd.h>
void output(char *str) { 
    usleep(100000);
    printf("%s\n", str);
}
*/
import "C"
import "unsafe"

// 2. 系统调用

//sysnb: syscall nonblocking
//sys: syscall blocking
```

**在执⾏ c 代码，或者阻塞在 syscall 上时，必须占⽤⼀个线程**

`sysmon: system monitor 处理阻塞`

*sysmon具有高优先级,在专有线程中执行,不需要绑定P*

| 作用 | 说明 |
| --- | --- |
| checkdead | 常见误解: 这个可以检查死锁,其实并不能 |
| netpoll | inject g list to global runqueue |
| retake | 1. syscall 卡了很久,那就把 P 剥离(handoffp) 2. 用户 G 运行很久,那么发信号 SIGURG 抢占 |

**刚开始是等待20us,后面超时重试,直到g 运行超过10ms或者syscall阻塞超过10ms后sysmon执行 retake, 执行成功就重置为20us**

## 调度器的发展历史

官方资料: <https://github.com/golang-design/history#scheduler>

## 与调度有关的常⻅问题

`Goroutine VS Thread`

|  | Goroutine | Thread |
| --- | --- | --- |
| 内存占用 | 2kb ~ 1GB | 从 8k 开始，服务端程序上限很多是 8M(⽤ulimit -a 可看)，调⽤多会 stack overflow |
| 上下文切换 | ⼏⼗ NS 级 | 1-2 us |
| 被谁管理 | Go runtime | 操作系统 |
| 通信⽅式  | CSP/传统共享内存  | 传统共享内存 |
| ID | 有，⽤户⽆法访问   | 有 |
| 抢占  | 1.13 以前需主动让出, 1.14 开始可由信号中断   | 内核抢占 |

`goroutine 的切换成本`

```go
// go/src/runtime/runtime2.go
type gobuf struct {
	// The offsets of sp, pc, and g are known to (hard-coded in) libmach.
	//
	// ctxt is unusual with respect to GC: it may be a
	// heap-allocated funcval, so GC needs to track it, but it
	// needs to be set and cleared from assembly, where it's
	// difficult to have write barriers. However, ctxt is really a
	// saved, live register, and we only ever exchange it between
	// the real register and the gobuf. Hence, we treat it as a
	// root during stack scanning, which means assembly that saves
	// and restores it doesn't need write barriers. It's still
	// typed as a pointer so that any other writes from Go get
	// write barriers.
	sp   uintptr
	pc   uintptr
	g    guintptr
	ctxt unsafe.Pointer
	ret  uintptr
	lr   uintptr
	bp   uintptr // for framepointer-enabled architectures
}
```

gobuf 描述⼀个 goroutine 所有现场，从⼀个 g 切换到另⼀个 g，只要把这⼏个现场字段保存下来，再把 g 往队列⾥⼀扔，m 就可以执⾏其它 g 了,**⽆需进⼊内核态**

`输出问题`

```go
package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	f2()
}

func f1() {
	runtime.GOMAXPROCS(1)
	for i := 0; i < 10; i++ {
		i := i
		go func() {
			fmt.Println("A: ", i)
		}()

	}

	var ch = make(chan int)
	<-ch
}

func f2() {
	runtime.GOMAXPROCS(1)
	for i := 0; i < 10; i++ {
		i := i
		go func() {
			fmt.Println("A: ", i)
		}()

	}

	time.Sleep(time.Hour)
}
```

```bash
[root@~]# go run main.go 
A:  9
A:  0
A:  1
A:  2
A:  3
A:  4
A:  5
A:  6
A:  7
A:  8
```

`死循环导致进程 hang 死问题`

```go
package main

import "time"

func main() {
    var i = 1;
    go func() {
        // 这个 goroutine 会导致进程在 gc 的时候 hang 死
        for {
            i++
        }
    }()
}
```

GC 时需要停⽌所有 goroutine,⽽⽼版本的 Go 的 g 停⽌需要主动让
出

1.14 增加基于信号的抢占之后，该问题被解决

`与 GMP 有关的⼀些缺陷`

1. **runtime 中有⼀个 allgs 数组所有创建过的 g 都会进该数组⼤⼩与 g 瞬时最⾼值相关**

2. **创建的 M 正常情况下是⽆法被回收**

```go
package main

/*
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
void output(char *str) {
    usleep(1000000);
    printf("%s\n", str);
}
*/
import "C"
import "unsafe"

import "net/http"
import _ "net/http/pprof"

func init() {
	go http.ListenAndServe(":9999", nil)
}

func main() {
    for i := 0;i < 1000;i++ {
        go func(){
            str := "hello cgo"
            //change to char*
            cstr := C.CString(str)
            C.output(cstr)
            C.free(unsafe.Pointer(cstr))

        }()
    }
    select{}
}
```

```go
package main

/*
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
void output(char *str) {
    usleep(1000000);
    printf("%s\n", str);
}
*/
import "C"

import (
	"net/http"
	"unsafe"

	"log"
	_ "net/http/pprof"
	"runtime"
	"sync"
)

func init() {
	go http.ListenAndServe(":9999", nil)
}

func main() {
	for i := 0; i < 1000; i++ {
		go func() {
			str := "hello cgo"
			//change to char*
			cstr := C.CString(str)
			C.output(cstr)
			C.free(unsafe.Pointer(cstr))

		}()
	}
	killThreadService()
	select {}
}

func sayhello(wr http.ResponseWriter, r *http.Request) {
	KillOne()
}

func killThreadService() {
	http.HandleFunc("/", sayhello)
	err := http.ListenAndServe(":10003", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

// KillOne kills a thread
func KillOne() {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		runtime.LockOSThread()
		return
	}()

	wg.Wait()
}
```

参考：<https://xargin.com/shrink-go-threads/>

## 总结

1. 可执⾏⽂件 ELF:

* 使⽤ go build -x 观察编译和链接过程
* 通过 readelf -H 中的 entry 找到程序⼊⼝
* 在 dlv 调试器中 b *entry_addr 找到代码位置

2. 启动流程：

* 处理参数 -> 初始化内部数据结构 -> 主线程 -> 启动调度循环

3. Runtime 构成：

* Scheduler、Netpoll、内存管理、垃圾回收

4. GMP：

* M，任务消费者；G，计算任务；P，可以使⽤ CPU 的 token

5. 队列：

* P 的本地 runnext 字段 -> P 的 local run queue -> global run queue，多级队列减少锁竞争

6. 调度循环：

* 线程 M 在持有 P 的情况下不断消费运⾏队列中的 G 的过程。

7. 处理阻塞：

* 可以接管的阻塞：channel 收发，加锁，⽹络连接读/写，select
* 不可接管的阻塞：syscall，cgo，⻓时间运⾏需要剥离 P 执⾏

8. sysmon:

* ⼀个后台⾼优先级循环，执⾏时不需要绑定任何的 P
* 负责：
  * 检查是否已经没有活动线程，如果是，则崩溃
  * 轮询 netpoll
  * 剥离在 syscall 上阻塞的 M 的 P
  * 发信号，抢占已经执⾏时间过⻓的 G

## References

ELF ⽂件解析：
<https://github.com/corkami/pics/blob/28cb0226093ed57b348723bc473cea0162dad366/binarelf101/elf101.pdf>

Go Scheduler 变更历史：
<https://github.com/golang-design/history#scheduler>

Goroutine vs Thread: 
<https://www.geeksforgeeks.org/golang-goroutine-vs-thread/>

Measuring context switching and memory overheads for Linux threads: 
<https://eli.thegreenplace.net/2018/measuring-context-switching-and-memory-overheads-for-linux-threads/>