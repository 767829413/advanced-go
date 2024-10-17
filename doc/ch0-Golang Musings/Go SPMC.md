# Golang 生产者消费者模式杂谈之SPMC

## 前言概述
-----------

Go 标准库和运行库中，有一些专门针对特定场景优化的数据结构，这些数据结构并没有暴露出来，这里就是简单介绍一下这些结构. 

## 生产者消费者模式介绍
-----------

众所周知,生产者消费者模式是一种常见的并发模式，根据生产者的数量和消费者的数量，可以分为四种情况：

* 单生产者-单消费者模式: `SPSC`
    
* 单生产者-多消费者模式: `SPMC`
    
* 多生产者-单消费者模式: `MPSC`
    
* 多生产者-多消费者模式: `MPMC`

## SPMC (单生产者多消费者)
-----------

首先介绍的是 `SPMC`(单生产者多消费者), 主要拿 Golang 内部的几个队列随便聊聊,看看所谓的 `lock-free`、`高性能`到底是个啥,这里主要是以下两个队列的实现:

* `PoolDequeue`:  是一个固定尺寸，使用 `ringbuffer (环形队列)` 方式实现的队列

* `PoolChain`: 基于 `PoolDequeue` 实现的一个动态尺寸的队列

`Channel` 基本上可以看做是一种 `MPMC (多生产者多消费者模式)`的队列. 可以同时允许多个生产者发送数据，也可以允许多个消费者消费数据，它也可以应用在其他模式的场景，比如在 `Go` 的 `rpc` 包中，`RPC（远程过程调用）`机制允许客户端调用远程服务器上的方法，就像调用本地方法一样. 每次调用都是一次请求-响应的过程，这种机制类似于 `oneshot` 模式(oneshot 通常是指一种一次性请求-响应的通信模式，即客户端发送一个请求，服务器处理后返回一个响应，然后通信结束)、通知情况下的的单生产者多消费者模式、rpc 和服务端单连接通讯时的消息处理，就是多生产者单消费者模式. 

但是 Go 标准库的 sync 包下，有一个针对单生产者多消费者的数据结构，它是一个 `lock-free` 的数据结构，针对这个场景做了优化，被使用在 `sync.Pool` 中. 

`sync.Pool` 采用了一种类似 Go 运行时调度的机制，针对每个 p 有一个 `private` 的数据，同时还有一个 `shared` 的数据，如果在本地 `private`、`shared` 中没有数据，就去其他 P 对应的 `shared` 去偷取. 这就可能有多个 P 偷取同一个 `shared`, 这是多消费者. 

同时对 `shared` 的写只有它隶属的 p 执行 `Put` 的时候才会发生：

```go
 l, _ := p.pin()
 if l.private == nil {
  l.private = x
 } else {
  l.shared.pushHead(x)
 }
 runtime_procUnpin()
```

针对于单生产者模式, `sync.Pool` 使用了 `PoolDequeue` 和 `PoolChain` 来做优化. 

下面先聊聊 `poolDequeue`. 

## poolDequeue
-----------

`poolDequeue` 是一个 `lock-free` 的数据结构，必然会使用 `atomic`, 同时它要求必须使用单生产者，否则会有并发问题. 消费者可以是并发多个，当然你用一个也没问题. 

其中，生产者可以使用下面的方法：

* **pushHead**: 在队列头部新增加一个数据. 如果队列满了，增加失败
    
* **popHead**： 在队列头部弹出一个数据. 生产者总是弹出新增加的数据，除非队列为空
    
消费者可以使用下面的一个方法：

* **popTail**: 从队尾处弹出一个数据，除非队列为空. 所以消费者总是消费最老的数据，这也正好符合大部分的场景

所以一般要结合代码分析了，很无聊,随便看看得了,建议跳过吧. 

### 代码分析

首先看这个 `struct` 的定义：

```go
type poolDequeue struct {
 headTail atomic.Uint64
 vals []eface
}
```

结构体有两个重要的属性：

* `headTail`： 一个 `atomic.Uint64` 类型的字段，它的高 32 位是 `head`，低 32 位是 `tail`. `head` 是下一个要填充的位置，`tail` 是最老的数据的位置. 
    
* `vals`： 一个 `eface` 类型的切片，它是一个环形队列，大小必须是 `2` 的幂次方. 
    
生产者增加数据的逻辑如下：

```go
func (d *poolDequeue) pushHead(val any) bool {
	ptrs := d.headTail.Load()
	head, tail := d.unpack(ptrs)
	if (tail+uint32(len(d.vals)))&(1<<dequeueBits-1) == head {
		// 队列满
		return false
	}
	slot := &d.vals[head&uint32(len(d.vals)-1)]

	// 检查 head slot 是否被 popTail 释放
	typ := atomic.LoadPointer(&slot.typ)
	if typ != nil {
		// 另一个 goroutine 正在清理 tail，所以队列还是满的
		return false
	}

	// 如果值为空，那么设置一个特殊值
	if val == nil {
		val = dequeueNil(nil)
	}
	// 队列头是空的，将数据写入 slot
	*(*any)(unsafe.Pointer(slot)) = val // ①

	// 增加 head，这样 popTail 就可以消费这个 slot 了
	// 同时也是一个 store barrier，保证了 slot 的写入
	d.headTail.Add(1 << dequeueBits)
	return true
}
```

① 处会有并发问题吗？万一有两个 goroutine 同时执行到这里，会不会有问题？这里没有问题，因为要求只有一个生产者，不会有另外一个 goroutine 同时写这个槽位. 

注意它还实现了`pack`和`unpack`方法，用于将 `head` 和 `tail` 打包到一个 `uint64` 中，或者从 `uint64` 中解包出 `head` 和 `tail`. 

消费者消费数据的逻辑如下：

```go
func (d *poolDequeue) popTail() (any, bool) {
	var slot *eface
	for { // ②
		ptrs := d.headTail.Load()
		head, tail := d.unpack(ptrs)
		if tail == head {
			// 队列为空
			return nil, false
		}

		// 确认头部和尾部（用于我们之前的推测性检查），并递增尾部。如果成功，那么我们就拥有了尾部的插槽。
		ptrs2 := d.pack(head, tail+1)
		if d.headTail.CompareAndSwap(ptrs, ptrs2) {
			// 成功读取了一个 slot
			slot = &d.vals[tail&uint32(len(d.vals)-1)]
			break
		}
	}

	// 剩下来就是读取槽位的值
	val := *(*any)(unsafe.Pointer(slot))
	if val == dequeueNil(nil) { // 如果本身就存储的nil
		val = nil
	}

	// 释放 slot，这样 pushHead 就可以继续写入这个 slot 了
	slot.val = nil                      // ③
	atomic.StorePointer(&slot.typ, nil) // ④

	return val, true
}
```

② 处是一个 for 循环，这是一个自旋的过程，直到成功读取到一个 slot 为止. 在有大量的 goroutine 的时候，这里可能会是一个瓶颈点，但是少量的消费者应该还不算大问题. 

③ 和 ④ 处是释放 slot 的过程，这样生产者就可以继续写入这个 slot 了. 

生产者还可以调用`popHead`方法，用来弹出刚刚压入还没有消费的数据:

```go
func (d *poolDequeue) popHead() (any, bool) {
	var slot *eface
	for {
		ptrs := d.headTail.Load()
		head, tail := d.unpack(ptrs)
		if tail == head {
			// 队列为空
			return nil, false
		}

		// 确认头部和尾部（用于我们之前的推测性检查），并递减头部。如果成功，那么我们就拥有了头部的插槽。
		head--
		ptrs2 := d.pack(head, tail)
		if d.headTail.CompareAndSwap(ptrs, ptrs2) {
			// 成功取回了一个 slot
			slot = &d.vals[head&uint32(len(d.vals)-1)]
			break
		}
	}

	val := *(*any)(unsafe.Pointer(slot))
	if val == dequeueNil(nil) {
		val = nil
	}

	// 释放 slot，这样 pushHead 就可以继续写入这个 slot 了
	*slot = eface{}
	return val, true
}
```

需要注意这是一个固定大小的队列，如果队列满了，生产者就会生产失败,需要等待. 同时这个队列的大小是 `2` 的幂次方，这样可以用 `&` 来取模，而不用 `%`，这样可以提高性能. 

## PoolChain
-----------

`PoolChain` 是在 `PoolDequeue` 的基础上实现的一个动态尺寸的队列，它的实现和 `PoolDequeue` 类似，只是增加了一个 `headTail` 的链表，用于存储多个 `PoolDequeue`. 

```go
type poolChain struct {
	// head 是生产者用来push的 poolDequeue。只有生产者访问，所以不需要同步
	head *poolChainElt

	// tail 是消费者用来pop的 poolDequeue。消费者访问，所以需要原子操作
	tail atomic.Pointer[poolChainElt]
}

type poolChainElt struct {
	poolDequeue

	// next由生产者原子写入，消费者原子读取。它只能从nil转换为非nil。
	// prev由消费者原子写入，生产者原子读取。它只能从非nil转换为nil。
	next, prev atomic.Pointer[poolChainElt]
}
```

这个都是抄这个文章 <https://github.com/golang/go/blob/master/src/sync/poolqueue.go#L220-L302> 里，想查看具体的实现可以细细品味.  整体的思想就是将多个`poolDequeue`串联起来，生产者在`head`处增加数据，消费者在`tail`处消费数据，当`tail`的`poolDequeue`为空时，就从`head`处获取一个`poolDequeue`.  当`head`满了的时候，就增加一个新的`poolDequeue`.  这样就实现了动态尺寸的队列. 

`sync.Pool`中就是使用的`PoolChain`来实现的，它是一个单生产者多消费者的队列，可以同时有多个消费者消费数据，但是只有一个生产者生产数据. 

要是想在实际业务中来使用这个结构完成一些事情，可以参考一下这个: 

<https://github.com/767829413/advanced-go/structure/poolqueue.go> , 单元测试和性能测试也凑合. 

可以学这个方法，使用类似的技术，创建一个 look-free 无限长度的 byte buffer. 在一些 Go 的网络优化库中就使用这种方法，避免频繁的 grow 和 copy 既有数据. 

## 与 channel 的性能比较
-----------

当然了,能不能打,打过了才知道,直接 `Benchmark` 跑一下就行,将 `poolDequeue`、`PoolChain` 和 `channel` 使用 1 个 goroutine 进行写入，10 个 goroutine 进行读取,粗略的评估一下性能优劣: 

```go
func BenchmarkPoolDequeue(b *testing.B) {
	const size = 1024
	pd := NewPoolDequeue(size)
	var wg sync.WaitGroup

	// Producer
	go func() {
		for i := 0; i < b.N; i++ {
			pd.PushHead(i)
		}
		wg.Done()
	}()

	// Consumers
	numConsumers := 10
	wg.Add(numConsumers + 1)
	for i := 0; i < numConsumers; i++ {
		go func() {
			for {
				if _, ok := pd.PopTail(); !ok {
					break
				}
			}
			wg.Done()
		}()
	}

	wg.Wait()
}

func BenchmarkPoolChain(b *testing.B) {
	pc := NewPoolChain()
	var wg sync.WaitGroup

	// Producer
	go func() {
		for i := 0; i < b.N; i++ {
			pc.PushHead(i)
		}
		wg.Done()
	}()

	// Consumers
	numConsumers := 10
	wg.Add(numConsumers + 1)
	for i := 0; i < numConsumers; i++ {
		go func() {
			for {
				if _, ok := pc.PopTail(); !ok {
					break
				}
			}
			wg.Done()
		}()
	}

	wg.Wait()
}

func BenchmarkChannel(b *testing.B) {
	ch := make(chan interface{}, 1024)
	var wg sync.WaitGroup

	// Producer
	go func() {
		for i := 0; i < b.N; i++ {
			ch <- i
		}
		close(ch)
		wg.Done()
	}()

	// Consumers
	numConsumers := 10
	wg.Add(numConsumers + 1)
	for i := 0; i < numConsumers; i++ {
		go func() {
			for range ch {
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
```

运行这个 benchmark,我们可以看到`poolDequeue`和`PoolChain`的性能要比`channel`高很多，大约是`channel`的 10 倍. `poolDequeue` 比 `PoolChain` 要好一些，性能是后者的两倍.  

测试结果如下:

```bash
goos: linux
goarch: amd64
pkg: github.com/767829413/advanced-go/structure
cpu: Intel(R) Core(TM) i7-8809G CPU @ 3.10GHz
BenchmarkPoolDequeue
BenchmarkPoolDequeue-8          63151449                17.87 ns/op            8 B/op          0 allocs/op
BenchmarkPoolChain
BenchmarkPoolChain-8            19922865                52.11 ns/op           34 B/op          0 allocs/op
BenchmarkChannel
BenchmarkChannel-8               9283756               142.3 ns/op             8 B/op          0 allocs/op
PASS
```
