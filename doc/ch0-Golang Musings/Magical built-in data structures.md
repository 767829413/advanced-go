# 神奇内置数据结构

## 内置数据结构⼀览

`runtime`

- channel
- timer
- semaphore
- map
- iface
- eface
- slice
- string

`sync`

- mutex
- cond
- pool
- once
- map
- waitgroup

`container`

- heap
- list
- ring

`os`

- os related

`context`

- context

`memory`

- allocation related
- gc related

`netpoll`

- netpoll related

## Channel

`演示动画: 基本执⾏流程`

<https://www.figma.com/proto/vfhlrTqsKicCO5ZbQZXgD4/runtime-structs?node-id=25-2&starting-point-node-id=25%3A2>

`调试代码`

```go
package main

import "time"

func main() {
    // 断点处
	var ch = make(chan int, 3)

	go func() {
		var i = 100
		for {
            // 断点处
			ch <- i
			i++
			time.Sleep(time.Second)
		}
	}()

	go func() {
		for {
			time.Sleep(3 * time.Second)
            // 断点处
			v := <-ch
			println(v)
		}
	}()
	println(1 << 16)
	select {}
}
```

![channel_recv_send_debug.png](https://s2.loli.net/2023/06/07/Jj6gwYr7W3ZnIAk.png)

`发送流程示意图`

![channel_send_flow.png](https://s2.loli.net/2023/06/06/sq3X5KtMDkLCSmz.png)

`接收流程示意图`

![channel_recv_flow.png](https://s2.loli.net/2023/06/06/clidBXekPQjRoVS.png)

`并发安全`

**channel底层还是通过用锁来实现的并发安全**

```go
// src/runtime/chan.go runtime.closechan runtime.chansend1 runtime.chanrecv1

// 这是发送
func chansend1(c *hchan, elem unsafe.Pointer) {
	chansend(c, elem, true, getcallerpc())
}

func chansend(c *hchan, ep unsafe.Pointer, block bool, callerpc uintptr) bool {
    // 省略 ...
    lock(&c.lock)
    // 省略 ...
    if c.closed != 0 {
		unlock(&c.lock)
		panic(plainError("send on closed channel"))
	}

	if sg := c.recvq.dequeue(); sg != nil {
		// Found a waiting receiver. We pass the value we want to send
		// directly to the receiver, bypassing the channel buffer (if any).
		send(c, sg, ep, func() { unlock(&c.lock) }, 3)
		return true
	}

	if c.qcount < c.dataqsiz {
		// Space is available in the channel buffer. Enqueue the element to send.
		qp := chanbuf(c, c.sendx)
		if raceenabled {
			racenotify(c, c.sendx, nil)
		}
		typedmemmove(c.elemtype, qp, ep)
		c.sendx++
		if c.sendx == c.dataqsiz {
			c.sendx = 0
		}
		c.qcount++
		unlock(&c.lock)
		return true
	}

	if !block {
		unlock(&c.lock)
		return false
	}
    // 省略 ...
}

// 这是接收
func chanrecv1(c *hchan, elem unsafe.Pointer) {
	chanrecv(c, elem, true)
}

func chanrecv(c *hchan, ep unsafe.Pointer, block bool) (selected, received bool) {
    // 省略 ...
    lock(&c.lock)
    // 省略 ...
    if c.closed != 0 {
		if c.qcount == 0 {
			if raceenabled {
				raceacquire(c.raceaddr())
			}
			unlock(&c.lock)
			if ep != nil {
				typedmemclr(c.elemtype, ep)
			}
			return true, false
		}
		// The channel has been closed, but the channel's buffer have data.
	} else {
		// Just found waiting sender with not closed.
		if sg := c.sendq.dequeue(); sg != nil {
			// Found a waiting sender. If buffer is size 0, receive value
			// directly from sender. Otherwise, receive from head of queue
			// and add sender's value to the tail of the queue (both map to
			// the same buffer slot because the queue is full).
			recv(c, sg, ep, func() { unlock(&c.lock) }, 3)
			return true, true
		}
	}

	if c.qcount > 0 {
		// Receive directly from queue
		qp := chanbuf(c, c.recvx)
		if raceenabled {
			racenotify(c, c.recvx, nil)
		}
		if ep != nil {
			typedmemmove(c.elemtype, ep, qp)
		}
		typedmemclr(c.elemtype, qp)
		c.recvx++
		if c.recvx == c.dataqsiz {
			c.recvx = 0
		}
		c.qcount--
		unlock(&c.lock)
		return true, true
	}

	if !block {
		unlock(&c.lock)
		return false, false
	}

    // 省略 ...
}

// 这是关闭
func closechan(c *hchan) {
	if c == nil {
		panic(plainError("close of nil channel"))
	}

	lock(&c.lock)
    // 省略 ...
    unlock(&c.lock)
}
```

`挂起和唤醒`

**可接管的阻塞一定是由 gopark 挂起,然后每一个 gopark 对应一个唤醒方来唤醒**

- Sender 挂起,一定是由 Receiver(或Close) 唤醒
- Receiver 挂起,一定是由 Sender(或Close) 唤醒

*gopark goready 的代码在 src/runtime/proc.go*

**gopark 对应的 goready 位置:**

- channel send -> channel recv/close
- Lock -> Unlock
- Read -> Read Ready，epoll_wait 返回了该 fd 事件时
- Timer -> checkTimers，检查到期唤醒

 ![goready_call.png](https://s2.loli.net/2023/06/07/5nwCjidXNlqg2JU.png)

 ![gopark_call.png](https://s2.loli.net/2023/06/07/qg7M3mIhdVnWyew.png)

## Timer

`四叉堆`

**早期 Timer 底层组织的数据结构**

 ![四叉堆.png](https://s2.loli.net/2023/06/07/p5zNJhwvkesRc8G.png)

性能问题: <https://github.com/golang/go/issues/15133>

`多个四叉堆`

 ![多个四叉堆.png](https://s2.loli.net/2023/06/07/zZlsLEI9HYWgrxU.png)

 **CPU 密集计算任务会导致 timer 唤醒延迟**

`在 schedule 中检查 timer`

 ![schedule_timer_call.png](https://s2.loli.net/2023/06/07/3y7n1AudprZ6kTb.png)

- 调整：
  - Timer heap 和 GMP 中的 P 绑定
  - 去除唤醒 goroutine: timerproc
- 检查：
  - 检查 timer 到期在特殊函数 checkTimers 中进⾏
  - 检查 timer 操作移⾄调度循环中进⾏
  - 主要逻辑在 src/runtime/proc.go schedule() -> findRunnable() -> checkTimers()

  ```go
  func checkTimers(pp *p, now int64) (rnow, pollUntil int64, ran bool) {
  	// If it's not yet time for the first timer, or the first adjusted
  	// timer, then there is nothing to do.
  	next := int64(atomic.Load64(&pp.timer0When))
  	nextAdj := int64(atomic.Load64(&pp.timerModifiedEarliest))
  	if next == 0 || (nextAdj != 0 && nextAdj < next) {
  		next = nextAdj
  	}
  
  	if next == 0 {
  		// No timers to run or adjust.
  		return now, 0, false
  	}
  
  	if now == 0 {
  		now = nanotime()
  	}
  	if now < next {
  		// Next timer is not ready to run, but keep going
  		// if we would clear deleted timers.
  		// This corresponds to the condition below where
  		// we decide whether to call clearDeletedTimers.
  		if pp != getg().m.p.ptr() || int(atomic.Load(&pp.deletedTimers)) <= int(atomic.Load(&pp.numTimers)/4) {
  			return now, next, false
  		}
  	}
  
  	lock(&pp.timersLock)
  
  	if len(pp.timers) > 0 {
  		adjusttimers(pp, now)
  		for len(pp.timers) > 0 {
  			// Note that runtimer may temporarily unlock
  			// pp.timersLock.
  			if tw := runtimer(pp, now); tw != 0 {
  				if tw > 0 {
  					pollUntil = tw
  				}
  				break
  			}
  			ran = true
  		}
  	}
  
  	// If this is the local P, and there are a lot of deleted timers,
  	// clear them out. We only do this for the local P to reduce
  	// lock contention on timersLock.
  	if pp == getg().m.p.ptr() && int(atomic.Load(&pp.deletedTimers)) > len(pp.timers)/4 {
  		clearDeletedTimers(pp)
  	}
  
  	unlock(&pp.timersLock)
  
  	return now, pollUntil, ran
  }
  ```

- ⼯作窃取：
  - 在 work-stealing 中，会从其它 P 那⾥偷 timer
  - 主要逻辑在 src/runtime/proc.go schedule() -> findRunnable() -> stealWork()

   ```go
   // stealWork attempts to steal a runnable goroutine or timer from any P.
   //
   // If newWork is true, new work may have been readied.
   //
   // If now is not 0 it is the current time. stealWork returns the passed time or
   // the current time if now was passed as 0.
   func stealWork(now int64) (gp *g, inheritTime bool, rnow, pollUntil int64, newWork bool) {
   	pp := getg().m.p.ptr()
   
   	ranTimer := false
   
   	const stealTries = 4
   	for i := 0; i < stealTries; i++ {
   		stealTimersOrRunNextG := i == stealTries-1
   
   		for enum := stealOrder.start(fastrand()); !enum.done(); enum.next() {
   			if sched.gcwaiting != 0 {
   				// GC work may be available.
   				return nil, false, now, pollUntil, true
   			}
   			p2 := allp[enum.position()]
   			if pp == p2 {
   				continue
   			}
   
   			// Steal timers from p2. This call to checkTimers is the only place
   			// where we might hold a lock on a different P's timers. We do this
   			// once on the last pass before checking runnext because stealing
   			// from the other P's runnext should be the last resort, so if there
   			// are timers to steal do that first.
   			//
   			// We only check timers on one of the stealing iterations because
   			// the time stored in now doesn't change in this loop and checking
   			// the timers for each P more than once with the same value of now
   			// is probably a waste of time.
   			//
   			// timerpMask tells us whether the P may have timers at all. If it
   			// can't, no need to check at all.
   			if stealTimersOrRunNextG && timerpMask.read(enum.position()) {
   				tnow, w, ran := checkTimers(p2, now)
   				now = tnow
   				if w != 0 && (pollUntil == 0 || w < pollUntil) {
   					pollUntil = w
   				}
   				if ran {
   					// Running the timers may have
   					// made an arbitrary number of G's
   					// ready and added them to this P's
   					// local run queue. That invalidates
   					// the assumption of runqsteal
   					// that it always has room to add
   					// stolen G's. So check now if there
   					// is a local G to run.
   					if gp, inheritTime := runqget(pp); gp != nil {
   						return gp, inheritTime, now, pollUntil, ranTimer
   					}
   					ranTimer = true
   				}
   			}
   
   			// Don't bother to attempt to steal if p2 is idle.
   			if !idlepMask.read(enum.position()) {
   				if gp := runqsteal(pp, p2, stealTimersOrRunNextG); gp != nil {
   					return gp, false, now, pollUntil, ranTimer
   				}
   			}
   		}
   	}
   
   	// No goroutines found to steal. Regardless, running a timer may have
   	// made some goroutine ready that we missed. Indicate the next timer to
   	// wait for.
   	return nil, false, now, pollUntil, ranTimer
   }
   ```

- 兜底：
  - runtime.sysmon 中会为 timer 未被触发(timeSleepUntil)兜底，启动新线程
  - 主要逻辑在 src/runtime/proc.go sysmon() 

  ```go
  func sysmon() {
  	lock(&sched.lock)
  	sched.nmsys++
  	checkdead()
  	unlock(&sched.lock)
  
  	lasttrace := int64(0)
  	idle := 0 // how many cycles in succession we had not wokeup somebody
  	delay := uint32(0)
  
  	for {
  		if idle == 0 { // start with 20us sleep...
  			delay = 20
  		} else if idle > 50 { // start doubling the sleep after 1ms...
  			delay *= 2
  		}
  		if delay > 10*1000 { // up to 10ms
  			delay = 10 * 1000
  		}
  		usleep(delay)
  
  		// sysmon should not enter deep sleep if schedtrace is enabled so that
  		// it can print that information at the right time.
  		//
  		// It should also not enter deep sleep if there are any active P's so
  		// that it can retake P's from syscalls, preempt long running G's, and
  		// poll the network if all P's are busy for long stretches.
  		//
  		// It should wakeup from deep sleep if any P's become active either due
  		// to exiting a syscall or waking up due to a timer expiring so that it
  		// can resume performing those duties. If it wakes from a syscall it
  		// resets idle and delay as a bet that since it had retaken a P from a
  		// syscall before, it may need to do it again shortly after the
  		// application starts work again. It does not reset idle when waking
  		// from a timer to avoid adding system load to applications that spend
  		// most of their time sleeping.
  		now := nanotime()
  		if debug.schedtrace <= 0 && (sched.gcwaiting != 0 || atomic.Load(&sched.npidle) == uint32(gomaxprocs)) {
  			lock(&sched.lock)
  			if atomic.Load(&sched.gcwaiting) != 0 || atomic.Load(&sched.npidle) == uint32(gomaxprocs) {
  				syscallWake := false
  				next := timeSleepUntil()
  				if next > now {
  					atomic.Store(&sched.sysmonwait, 1)
  					unlock(&sched.lock)
  					// Make wake-up period small enough
  					// for the sampling to be correct.
  					sleep := forcegcperiod / 2
  					if next-now < sleep {
  						sleep = next - now
  					}
  					shouldRelax := sleep >= osRelaxMinNS
  					if shouldRelax {
  						osRelax(true)
  					}
  					syscallWake = notetsleep(&sched.sysmonnote, sleep)
  					if shouldRelax {
  						osRelax(false)
  					}
  					lock(&sched.lock)
  					atomic.Store(&sched.sysmonwait, 0)
  					noteclear(&sched.sysmonnote)
  				}
  				if syscallWake {
  					idle = 0
  					delay = 20
  				}
  			}
  			unlock(&sched.lock)
  		}
          // 省略 ...
  ```

## Map

`特权语法`

```go
package main

// map 分配到栈上时,不一定会调用 makemap
// runtime.makemap n > 8
var m = make(map[int]int, 9)
// runtime.makemap_small n <= 8
var mm = make(map[int]int, 8)

func main() {
	// runtime.mapaccess1_fast64
	v1 := m[1]
	// runtime.mapaccess2_fast64
	v2, ok := m[2]
	println(v1, v2, ok)
}
```

`Map-函数⼀览`

 ![map_func.png](https://s2.loli.net/2023/06/07/FCMlnxazsXfei3Q.png)

*map 中⼤量类似但⼜冗余的函数，原因之⼀便是没有泛型*

`Map-结构图`

 ![map_struct.png](https://s2.loli.net/2023/06/09/6Bjbrn3OAZqml1J.png)

`Map-元素操作`

 **mapaccess,mapassign,mapdelete**

 ![map_op.png](https://s2.loli.net/2023/06/09/FRVqlBtzD7Q46Pw.png)

`Map-扩容`

- 触发：mapassign
- 时机：load factor 过⼤ || overflow bucket 过多
- 搬运过程是渐进进⾏的

![map_expand.png](https://s2.loli.net/2023/06/09/9wWBElDrYH7bQcP.png)

动画演示: <https://www.figma.com/proto/vfhlrTqsKicCO5ZbQZXgD4/runtime-structs?node-id=111-371&starting-point-node-id=111%3A371>

**扩容中**

- mapasssign：将命中的 bucket 从 oldbuckets 顺⼿搬运到buckets 中，顺便再多搬运⼀个 bucket
- mapdelete：将命中的 bucket 从 oldbuckets 顺⼿搬运到buckets 中，顺便再多搬运⼀个 bucket
- mapaccess: 优先在 oldbuckets 中找，如果命中，则说明这个 bucket 没有被搬运

*搬运 bucket x 时，会被该桶的 overflow 桶也⼀并搬完*

`Map-遍历`

 动画演示: <https://www.figma.com/proto/vfhlrTqsKicCO5ZbQZXgD4/runtime-structs?node-id=116-368&starting-point-node-id=116%3A368>

`Map-缺陷`

- 已经扩容的 map，⽆法收缩

 ```go
 package main
 
 import (
 	"fmt"
 	"runtime"
 )
 
 func main() {
 	m := make(map[int]int) // 创建一个map对象
 
 	var startMemStats runtime.MemStats
 	runtime.ReadMemStats(&startMemStats) // 获取map创建时的内存占用
 
 	for i := 0; i < 10000000; i++ {
 		m[i] = i // 向map中添加键值对，触发扩容
 	}
 
 	var midMemStats runtime.MemStats
 	runtime.ReadMemStats(&midMemStats) // 获取map扩容后的内存占用
 
 	for i := 0; i < 9990000; i++ {
 		delete(m, i) // 从map中删除一半的键值对，触发缩容
 	}
 
 	var endMemStats runtime.MemStats
 	runtime.ReadMemStats(&endMemStats) // 获取map缩容后的内存占用
 
 	fmt.Printf("Start mem usage: %d bytes\n", startMemStats.Alloc)
 	fmt.Printf("Mid mem usage: %d bytes\n", midMemStats.Alloc)
 	fmt.Printf("End mem usage: %d bytes\n", endMemStats.Alloc)
 
 	fmt.Printf("Memory used for expanding the map: %d bytes\n", midMemStats.Alloc-startMemStats.Alloc)
 	fmt.Printf("Memory used for shrinking the map: %d bytes\n", endMemStats.Alloc-midMemStats.Alloc)
 }
 ```

- 保证并发安全时，要⼿动读写锁，易出错

 ```go
 package main
 
 import (
 	"sync"
 )
 
 type mapWithLock struct {
 	m   map[int]int
 	mux sync.RWMutex
 }
 
 func (m *mapWithLock) readMap(idx int) int {
 	m.mux.RLock()
 	defer m.mux.RUnlock()
 	v := m.m[idx]
 	// Do some thing
 
 	return v
 }
 ```

- 多核⼼下表现差

 ```go
 package main
 
 import (
 	"fmt"
 	"runtime"
 	"sync"
 )
 
 func main() {
 	runtime.GOMAXPROCS(runtime.NumCPU())
 
 	m := make(map[int]int)
 	wg := sync.WaitGroup{}
 	mutex := sync.Mutex{}
 
 	for i := 0; i < 10000000; i++ {
 		wg.Add(1)
 		go func(j int) {
 			mutex.Lock()
 			m[j] = j
 			mutex.Unlock()
 			wg.Done()
 		}(i)
 	}
 
 	wg.Wait()
 	fmt.Println(len(m))
 }
 ```

- 难以使⽤ sync.Pool 进⾏重⽤

 ```go
 package main
 
 import "sync"
 
 var slicePool = sync.Pool{}
 var mapPool = sync.Pool{}
 
 // slice can be easily reused
 func processUserRequest1() {
 	sl := slicePool.Get().([]any)
 	defer func() {
 		sl := sl[:0]
 		slicePool.Put(sl)
 	}()
 	// processs user logic
 }
 
 // what about map?
 func processUserRequest2() {
 	m := mapPool.Get()
 	defer func() {
 		// how to reset a map
 		mapPool.Put(m)
 	}()
 }
 ```

## Context

`Context方法概览`

 ![context_some.png](https://s2.loli.net/2023/06/09/m1QtcHZSMVs6a2v.png)

- emptyCtx，所有 ctx 类型的根
- valueCtx，主要就是为了在 ctx 中嵌⼊上下⽂数据，⼀个简单的 k 和 v结构，同⼀个 ctx 内只⽀持⼀对 kv，需要更多的 kv 的话，会形成⼀棵树形结构. 
- cancelCtx，⽤来取消程序的执⾏树
- timerCtx，在 cancelCtx 上包了⼀层，⽀持基于时间的 cancel. 

`Context value结构`

 ![context_value1.png](https://s2.loli.net/2023/06/09/kg5KuqnY47zcGx3.png)
 ![context_value2.png](https://s2.loli.net/2023/06/09/woPTFsBKNXu6Mr3.png)

`树形结构`

 ![context_tree.png](https://s2.loli.net/2023/06/09/Lg5FPOnHWiDN9Ze.png)

 **⽗节点取消时，可以传导到所有⼦节点**

 ![context_cancel.png](https://s2.loli.net/2023/06/09/hfWsmHKEwjq5voz.png)

## References

- Context: <https://github.com/cch123/golang-notes/blob/master/context.md>
- ⽼版本的 timer 实现：<https://github.com/cch123/golang-notes/blob/master/timer.md>
- 1.14 timer 性能提升分析：<http://xiaorui.cc/archives/6483>
- 这位⼩姐姐的 PPT 还是做的不错的，未覆盖到的细节补充: <https://speakerdeck.com/kavya719/understanding-channels>
- DPVS 时间轮：<https://www.jianshu.com/p/f38cd8c99f70>
- Kafka 时间轮：<https://www.infoq.cn/article/erdajpj5epir65iczxzi>