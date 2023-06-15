# Go 并发编程最佳实践

## 并发内置数据结构

`sync.Once`

**sync.Once 只有⼀个⽅法，Do(),但 o.Do 需要保证：**

* 初始化⽅法必须且只能被调⽤⼀次
* Do 返回后，初始化⼀定已经执⾏完成

源码: <https://github.com/golang/go/blob/master/src/sync/once.go>

```go
type Once struct {
	done uint32
	m    Mutex
}

func (o *Once) Do(f func()) {
	if atomic.LoadUint32(&o.done) == 0 {
		o.doSlow(f)
	}
}

func (o *Once) doSlow(f func()) {
	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 {
		defer atomic.StoreUint32(&o.done, 1)
		f()
	}
}
```

*官方的注释非常详细,这里就不把注释带出来了*

`sync.Pool`

**主要在两种场景使⽤：**

* 查看 pprof 中进程中的 inuse_objects 数过多，gc mark 消耗⼤量 CPU
* 查看 pprof 中进程中的 inuse_objects 数过多，进程 RSS 占⽤过⾼

请求⽣命周期开始时，pool.Get，请求结束时，pool.Put。在 [fasthttp](https://github.com/valyala/fasthttp/blob/b433ecfcbda586cd6afb80f41ae45082959dfa91/server.go#L402) 中有⼤量应⽤

**sync.Pool结构图解**

 ![sync_pool.png](https://pic.imgdb.cn/item/648928631ddac507cc4d8b28.png)

 ```go
 // A Pool is a set of temporary objects that may be individually saved and
// retrieved.
//
// Any item stored in the Pool may be removed automatically at any time without
// notification. If the Pool holds the only reference when this happens, the
// item might be deallocated.
//
// A Pool is safe for use by multiple goroutines simultaneously.
//
// Pool's purpose is to cache allocated but unused items for later reuse,
// relieving pressure on the garbage collector. That is, it makes it easy to
// build efficient, thread-safe free lists. However, it is not suitable for all
// free lists.
//
// An appropriate use of a Pool is to manage a group of temporary items
// silently shared among and potentially reused by concurrent independent
// clients of a package. Pool provides a way to amortize allocation overhead
// across many clients.
//
// An example of good use of a Pool is in the fmt package, which maintains a
// dynamically-sized store of temporary output buffers. The store scales under
// load (when many goroutines are actively printing) and shrinks when
// quiescent.
//
// On the other hand, a free list maintained as part of a short-lived object is
// not a suitable use for a Pool, since the overhead does not amortize well in
// that scenario. It is more efficient to have such objects implement their own
// free list.
//
// A Pool must not be copied after first use.
//
// In the terminology of the Go memory model, a call to Put(x) “synchronizes before”
// a call to Get returning that same value x.
// Similarly, a call to New returning x “synchronizes before”
// a call to Get returning that same value x.
type Pool struct {
	noCopy noCopy

	local     unsafe.Pointer // local fixed-size per-P pool, actual type is [P]poolLocal
	localSize uintptr        // size of the local array

	victim     unsafe.Pointer // local from previous cycle
	victimSize uintptr        // size of victims array

	// New optionally specifies a function to generate
	// a value when Get would otherwise return nil.
	// It may not be changed concurrently with calls to Get.
	New func() any
}
 ```

* 根据上图,应用层如果频繁生成 sync.Pool 会导致频繁的对 runtime.allPools 做 append 和 locaked,性能问题会很严重
* noCopy表示不要copy Pool 结构体,主要是因为可能 copy 了锁 或者只是浅拷贝,这样会导致各种问题
* local,localSize 是一对,local是pool的数组对象,大小应该和 P 的数量对应
* 在初始化 sync.Pool 的时候需要提供的一个 new 函数
* 获取 sync.Pool 里的对象开始是在当前 P 下面获取对应 [P]poolLocal 里的 poolLocalInternal,都是先取 poolLocalInternal.private,如果为空就去 poolLocalInternal.shared 的链表里取找,如果还找不到那就去其他 P 下面的 shared 里去偷一个,如果没有偷到, 那就去 victim(和local类似)里面找

**sync.Pool GC时候操作**

 ![sync_pool_1.png](https://pic.imgdb.cn/item/64892e5d1ddac507cc613d36.png)

* 发生GC应该就是直接把 local,localSize 平移到 victim,victimSize,然后原来的 victim,victimSize 值就会被丢弃
* 如果GC结束后来取对象,那么因为 local 为空,会到 victim 里面取找,流程类似上述的 local 查找过程
* 1.13之前包括1.13 share 都是有锁的,后面采用无锁的双端链表,以前有锁版本在 GC 结束的时候,如果重度依赖 Pool 且有大量的获取对象操作,会导致阻塞在 share 的锁上

`semaphore(信号量)`

**semaphore(信号量) 是锁的实现基础，所有同步原语的基础设施**

![sync_pool_2.png](https://pic.imgdb.cn/item/648932011ddac507cc6dba41.png)

`sync.Mutex(互斥锁)`

 ![sync_1.png](https://s2.loli.net/2023/06/14/Mymdh9jZXioJxwb.png)

`sync.RWMutex`

 ![sync_2.png](https://s2.loli.net/2023/06/14/56BGivoA7I1HsPa.png)

`sync.Map`

 演示动画: <https://www.figma.com/proto/FMzUIdkjm4BEHSpwWFecew/concurrency?node-id=6-16&starting-point-node-id=6%3A16>

`sync.Waitgroup`

 ![sync_3.png](https://s2.loli.net/2023/06/14/JWfn69PGpRrXbKg.png)

## 并发编程模式举例

`CSP 和传统并发模式`

 ![sync_4.png](https://s2.loli.net/2023/06/14/RGACgD8r7MztmNu.png)

`Fan-in: 合并多个 channel 操作`

```go
package main

import (
	"time"
)

func main() {
	var ch1, ch2, ch3, ch4, ch5, ch6 = make(chan any), make(chan any), make(chan any), make(chan any), make(chan any), make(chan any)

	go func() {
		time.Sleep(1 * time.Second)
		ch1 <- 1
		close(ch1)
	}()
	go func() {
		time.Sleep(1 * time.Second)
		ch2 <- 2
		close(ch2)
	}()
	go func() {
		time.Sleep(1 * time.Second)
		ch3 <- 3
		close(ch3)
	}()
	go func() {
		time.Sleep(1 * time.Second)
		ch4 <- 4
		close(ch4)
	}()
	go func() {
		time.Sleep(1 * time.Second)
		ch5 <- 5
		close(ch5)
	}()
	go func() {
		time.Sleep(1 * time.Second)
		ch6 <- 6
		close(ch6)
	}()
	res := fanInRec(ch1, ch2, ch3, ch4, ch5, ch6)
	for v := range res {
		println(v.(int))
	}

}

func fanInRec(chans ...<-chan any) <-chan any {
	switch len(chans) {
	case 0:
		c := make(chan any)
		close(c)
		return c
	case 1:
		return chans[0]
	case 2:
		return mergeTwo(chans[0], chans[1])
	default:
		m := len(chans) / 2
		return mergeTwo(fanInRec(chans[:m]...), fanInRec(chans[m:]...))
	}

}

func mergeTwo(a, b <-chan any) <-chan any {
	c := make(chan any)
	go func() {
		defer close(c)
		for a != nil || b != nil {
			select {
			case v, ok := <-a:
				if !ok {
					a = nil
					continue
				}
				c <- v
			case v, ok := <-b:
				if !ok {
					b = nil
					continue
				}
				c <- v
			}
		}
	}()
	return c
}
```

`Or channel: 任意 channel 返回全部返回`

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	var ch1 = make(chan any)
	var ch2 = make(chan any)
	var ch3 = make(chan any)
	var ch4 = make(chan any)
	var ch5 = make(chan any)

	go func() {
		defer close(ch1)
		time.Sleep(1 * time.Second)
		ch1 <- 1
	}()

	go func() {
		defer close(ch2)
		time.Sleep(2 * time.Second)
		ch2 <- 2
	}()

	go func() {
		defer close(ch3)
		time.Sleep(2 * time.Second)
		ch3 <- 3
	}()

	go func() {
		defer close(ch4)
		time.Sleep(2 * time.Second)
		ch4 <- 4
	}()

	go func() {
		defer close(ch5)
		time.Sleep(3 * time.Second)
		ch5 <- 5
	}()

	res := or(ch1, ch2, ch3, ch4, ch5)
	v := <-res
	fmt.Println(v)
}

func or(channels ...<-chan any) <-chan any {
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}
	orDone := make(chan any)
	go func() {
		defer close(orDone)

		switch len(channels) {
		case 2:
			select {
			case v := <-channels[0]:
				orDone <- v
			case v := <-channels[1]:
				orDone <- v
			}

		default:
			select {
			case v := <-channels[0]:
				orDone <- v
			case v := <-channels[1]:
				orDone <- v
			case v := <-channels[2]:
				orDone <- v
			case v := <-or(append(channels[3:], orDone)...):
				orDone <- v
			}
		}
	}()
	return orDone
}
```

`Pipeline: 串联在⼀起的 channel`

```go
package main

import (
	"fmt"
)

func main() {
	var c = make(chan int, 3)
	c <- 10
	c <- 20
	c <- 30
	out := sq(c)
	fmt.Println(<-out)
	fmt.Println(<-out)
	fmt.Println(<-out)
}

func sq(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}
```

`并发同时保序`

```go
type task struct {
	sync.WaitGroup
}

func main() {
	for t := range <-task {
		sendToWorker(ch,t)
		sendToKeepOrder(fifo,t)
	}

	for t ;= range fifo { 
		t .wait()
	}
}

func worker(ch) {
	for t := range ch {
		t.Done()
	}
}
```

## 常⻅的并发 bug

`死锁`

**死锁问题可以通过 pprof 进⼊ goroutine ⻚⾯查看**

```go
// 重入锁
func x() {
	a.RLock()
	defer a.RUnlock()

	y()
}

func y() {
	a.RLock()
	defer a.RUnlock()
}

// 循环依赖
func x() {
	a.lock()
	// do something
	b.lock()
}

func y() {
	b.lock()
	// do something
	a.lock()
}

// goroutine 之间循环依赖
func goroutine1() {
	m.Lock()
	// ch <- request // blocks
	select {
		case ch <- request:
		default:
	}
	m.Unlock()
}

func goroutine2() {
	for {
		m.Lock()// blocks
		m.Unlock()
		request <- ch
	}
}
```

`Map concurrent writes/reads`

**崩溃时输出 stderr，请注意重定向你的 stderr 到单独的⽂件中**

`Channel 关闭 panic`

**通道关闭原则**

1. 有多个消费者,一个生产者,需要生产者来主动通知(close(chan))消费者数据的发送已经结束
2. 一个消费者，多个生产者，通过添加一个额外的channel被消费者来通知所有生产者接收数据已经结束,停止生产,生产者在逻辑中要有判断
3. M个接收者，N个发送者，其中任何一个说 "让我们结束游戏"，通知主持人关闭一个额外的信号通道。

参考资料: <https://go101.org/article/channel-closing.html>

`fn() 超时后，ch <- result 阻塞导致 goroutine 永久泄露`

**Figure 1. A blocking bug caused by channel**

```go
func finishReq(timeout time.Duration) r ob {
	// ch := make(chan ob)
	ch := make(chan ob, 1)
	go func() {
		result := fn()
		ch <- result // block
	}()

	select {
		case result := <-ch:
			return result
		case <-time.After(timeout):
			return nil
	}
}
```

`wait group 使⽤不当，永久阻塞`

**Figure 5. A blocking bug caused by WaitGroup**

```go
var group sync.WaitGroup
group.Add(len(pm.plugins))
for _, p := range pm.plugins {
	go func(p *pm.Plugin) {
		defer group.Done()
	}()
	// group.Wait
}
group.Wait()
```

`context.WithCancel 内部启动 goroutine，在 ctx 被覆盖后泄露`

**Figure 6. A blocking bug caused by context**

```go
// hctx, hcancel := context.WithCancel(ctx)
var hctx context.Context
var hcancel context.CancelFunc
if timeout > 0 {
	hctx, hcancel = context.WithTimeout(ctx)
}else {
	hctx, hcancel = context.WithCancel(ctx)
}
```

`闭包捕获本地变量`

**Figure 8. A data race caused by anonymous function.**

```go
for i:=17; i <= 21; i++ {
	// go func() {
	// 	  apiVersion := fmt.Sprintf("v1.%d", i)
	// }()
	go func(i int) {
		apiVersion := fmt.Sprintf("v1.%d", i)
	}(i)
}
```

`启动 goroutine 前要保证 Add 完成`

**Figure 9. A non-blocking bug caused by misusingWaitGroup**

```go
func (p *peer) send() {
	p.mu.Lock()
	defer p.mu.Unlock()
	switch p.status {
		case idle:
			p.wg.Add(1)
			go func() {
				// p.wg.Add(1)
				...
				p.wg.Done()
			}()
		case stopped:
	}
}

func (p *peer) stop() {
	p.mu.Lock()
	p.status = stopped
	p.mu.Unlock()
	p.wg.Wait()
}
```

`并发操作 channel 时，多次关闭同⼀个 ch`

**Figure 10. A bug caused by closing a channel twice**

```go
//select {
//	case <- c.closed:
//	default:
//		close(c.closed)
//}
select {
	case <- c.closed:
	default:
	    sync.Once(func(){
			close(c.closed)
		})
}
```

`Fn 耗时很久但进⼊之前没有判断外部给的 stopCh 中的通知导致浪费算⼒`

**Figure 11. A non-blocking bug caused by select and channel.**

```go
ticker := time.NewTicker()
//for {
//	fn()
//	select {
//		case <- stopCh:
//			return
//		case <- ticker:
//	}
//}
for {
	select {
		case <- stopCh:
			return
		default:
	}
	fn()
	select {
		case <- stopCh:
			return
		case <- ticker:
	}
}
```

## 内存模型

**有并发问题使⽤显式同步(锁 channel)就可以保证正确性**

`现代计算机的多级存储结构`

 *L1D cache ⼜会被划分为多个cache line，每个 cache line = 64 bytes*

 ![mm-1.png](https://s2.loli.net/2023/06/15/iL3fpskXxOeGN5P.png)

 参考资料: <http://15418.courses.cs.cmu.edu/spring2015/lecture/basicarch/slide_042>

 *L1 cache ⼜被划分为更细粒度的 cacheline，下⾯是某个服务器上获取 L1 cache line size 的命令*

 ```bash
 ➜  ~  getconf LEVEL1_DCACHE_LINESIZE
 64
 ```

 Runtime 中的 cacheline pad : <https://github.com/golang/go/blob/master/src/runtime/mgc.go>

 cpu.CacheLinePad

 ```go
var work workType

type workType struct {
	full  lfstack          // lock-free list of full blocks workbuf
	_     cpu.CacheLinePad // prevents false-sharing between full and empty
	empty lfstack          // lock-free list of empty blocks workbuf
	_     cpu.CacheLinePad // prevents false-sharing between empty and nproc/nwait
 ```

`多核⼼给我们带来的问题`

* 单变量的并发操作也必须⽤同步⼿段，⽐如 atomic
* 全局视⻆下观察到的多变量读写的顺序可能会乱序

**单变量的原⼦读/写，多核⼼使⽤ mesi 协议保证正确性**
 
 ![1](https://pic.imgdb.cn/item/648a85bb1ddac507ccb7e379.png)
 ![2](https://pic.imgdb.cn/item/648a85bb1ddac507ccb7e3a9.png)

*Mesi 协议是以整个 cache line 为单位进⾏的*

参考资料: <https://www.scss.tcd.ie/Jeremy.Jones/VivioJS/caches/MESIHelp.htm>

**多核⼼执⾏时，CPU 和编译器可能对读写指令进⾏重排**

使⽤ Litmus 测试观察内存重排：

 ![1](https://pic.imgdb.cn/item/648a8aec1ddac507ccc34dfe.png)

上边的伪代码可以认为是：

 ![2](https://pic.imgdb.cn/item/648a8aec1ddac507ccc34de8.png)

然后检查两个核⼼的 EAX 是不是都为 0

工具链接: <https://github.com/herd/herdtools7>

`False sharing`

**因为 CPU 处理读写是以 cache line 为单位，所以在并发修改变量时，会⼀次性将其它 CPU core 中的 cache line invalidate 掉，导致未修改的内存上相邻的变量也需要同步，带来额外的性能负担。**

`True sharing`

**多线程确实在共享并更新同⼀个变量/内存区域**

`Happen-before`

* 同⼀个 goroutine 内的逻辑有依赖的语句执⾏，满⾜顺序关系。
* 编译器/CPU 可能对同⼀个 goroutine 中的语句执⾏进⾏打乱，以提⾼性能，但不能破坏其应⽤原有的逻辑。
* 不同的 goroutine 观察到的共享变量的修改顺序可能不⼀样。

 **初始化：**

* A pkg import B pkg，那么 B pkg 的 init 函数⼀定在 A pkg 的 init 函数之前执⾏。
* Init 函数⼀定在 main.main 之前执⾏

 **Goroutine 创建：**

* Goroutine 的创建(creation)⼀定先于 goroutine 的执⾏ (execution)

 **Goroutine 结束：**

* 在没有显式同步的情况下，goroutine 的结束没有任何保证，可能被执⾏，也可能不被执⾏

 **Channel 收/发：**

* A send on a channel happens before the corresponding receive from that channel completes.
 
 发送一定在接收之前完成

```go
var c = make(chan int, 10)
var a string

func f() {
	a = "hi"
	c <- 0
}

func main() {
	go f()
	<-c
	println(a)
}
```

 *这⾥ c <- 0 ⼀定先于 <- c 执⾏完所以 print ⼀定能打印出 hi*

* The closing of a channel happens before a receive that returns a zero value because the channel is closed.

 关闭一定是在接收之前

```go
var c = make(chan int, 10)
var a string

func f() {
	a = "hi"
	close(c)
}

func main() {
	go f()
	<-c
	println(a)
}
```

 *close(c) ⼀定先于 <-c 执⾏完,所以这⾥也可以保证打印出 hi*

* A receive from an unbuffered channel happens before the send on that channel completes.

```go
var c = make(chan int)
var a string

func f() {
	a = "hi"
	<-c
}

func main() {
	go f()
	c <- 0
	println(a)
}
```

 *⽆ buffer 的 chan receive 先于 send 执⾏完，这⾥也可以保证打印出 hi*

 **Lock && Unlock：**

* For any sync.Mutex or sync.RWMutex variable l and n < m, call n of l.Unlock() happens before call m of l.Lock() returns.
* Unlock ⼀定先于 Lock 函数返回前执⾏完

 **Once：**

* A single call of f() from once.Do(f) happens (returns) before any call of once.Do(f) returns.

 **总结 Happen-before**

 *本质是在⽤户不知道 memory barrier 概念和具体实现的前提下，能够按照官⽅提供的 happen-before 正确进⾏并发编程。*

`Memory barrier`

 **在并发编程中的 memory barrier 和 GC 中的 barrier 不是⼀回事。**

 **Memory barrier 是为了防⽌各种类型的读写重排：**
  
  ![mm-6.png](https://s2.loli.net/2023/06/15/MrtpGQ6VEa7fSUc.png)

 **⽽ GC 中的 read/write barrier 则是指堆上指针修改之前插⼊的⼀⼩段代码。**

## References

* <https://wudaijun.com/2018/02/go-sync-map-implement/>
* <https://github.com/kat-co/concurrency-in-go-src>
* <https://speakerdeck.com/kavya719/understanding-channels>
* <https://www.zenlife.tk/concurrency-with-keep-order.md?hmsr=joyk.com&utm_source=joyk.com&utm_medium=referral>
* <https://golang.org/ref/mem>
* <https://www.hardwaretimes.com/difference-between-l1-l2-and-l3-cache-what-is-cpu-cache/>
* <https://github.com/lotusirous/go-concurrency-patterns>
* <https://songlh.github.io/paper/go-study.pdf>
* <https://github.com/cch123/golang-notes/blob/master/memory_barrier.md>

## 未涉及

* 内置并发结构：sync.Cond
* 进阶话题：如 acquire、release、sequential consistency、Lock-Free，Wait-free 等等
* 扩展并发原语：SingleFlight，ErrGroup 等