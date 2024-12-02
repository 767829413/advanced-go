# Golang 运行时调度杂谈之runq

在说调度之前不得不回顾一下 `Go` 运行时的 `GPM` 模型, 网上资料一堆,这里只是稍微点一点,当然了,前面也文章也有详细说明. 

[Golang 程序是怎么跑起来的](./How%20golang%20works.md)

> “
> 
> `GPM` 模型中的 `G` 代表 `goroutine`. 每个 `goroutine` 只占用几 `KB` 的内存,可以轻松创建成千上万个. `G` 包含了 `goroutine` 的栈、指令指针和其他信息,如阻塞 `channel` 的等待队列等. 
> 
> `P` 代表 `processor`,可以理解为一个抽象的 `CPU` 核心. `P` 的数量默认等于实际的 `CPU` 核心数,但可以通过环境变量进行调整. `P` 维护了一个本地的 `goroutine` 队列,还负责执行 `goroutine` 并管理与之关联的上下文信息. 
> 
> `M` 代表 `machine`,是操作系统线程. 一个 `M` 必须绑定一个 `P` 才能执行 `goroutine`. 当一个 `M` 阻塞时,运行时会创建一个新的 `M` 或者复用一个空闲的 `M` 来保证 `P` 的数量总是等于 `GOMAXPROCS` 的值,从而充分利用 `CPU` 资源. 
> 
> 在这个模型中,`P` 扮演了承上启下的角色. 它连接了 `G` 和 `M`,实现了用户层级的 `goroutine` 到操作系统线程的映射. 这种设计允许 `Go` 在用户空间进行调度,避免了频繁的系统调用,大大提高了并发效率. 
> 
> 调度过程中,当一个 `goroutine` 被创建时,它会被放到 `P` 的本地队列或全局队列中. 如果 `P` 的本地队列已满,一些 `goroutine` 会被放到全局队列. 当 `P` 执行完当前的 `goroutine` 后,会优先从本地队列获取新的 `goroutine` 来执行. 如果本地队列为空,`P` 会尝试从全局队列或其他 `P` 的队列中偷取 `goroutine`. 
> 
> 这种`工作窃取(work-stealing)`算法确保了负载的动态平衡. 当某个 `P` 的本地队列为空时,它可以从其他 `P` 的队列中窃取一半的 `goroutine`,这有效地平衡了各个 `P` 之间的工作负载. 

`Go` 运行时这么做, 主要还是减少 `P` 之间对获取 `goroutine` 之间的竞争. 本地队列 `runq` 主要由持有它的 `P` 进行读写, 只有在"被偷"的情况下, 才可能有"数据竞争"的问题, 而这种情况发生概率较少, 所以它设计了一个高效的 `runq` 数据结构来应对这么场景. 实际看起来和上面介绍的 [PoolDequeue](./Go%20SPMC.md) 有异曲同工之妙. 

***简单介绍一下 `global queue` 等数据结构, 但不是重点.***

## runq
------------

详细代码可以参考: <https://github.com/golang/go/blob/master/src/runtime/runtime2.go>
在运行时中 `P` 是一个复杂的数据结构, 下面列出了需要关注的几个字段:

```go
// 一个goroutine的指针
type guintptr uintptr

//go:nosplit
func (gp guintptr) ptr() *g { return (*g)(unsafe.Pointer(gp)) }

//go:nosplit
func (gp *guintptr) set(g *g) { *gp = guintptr(unsafe.Pointer(g)) }

//go:nosplit
func (gp *guintptr) cas(old, new guintptr) bool {
	return atomic.Casuintptr((*uintptr)(unsafe.Pointer(gp)), uintptr(old), uintptr(new))
}

type p struct {
	id          int32
	status      uint32 // one of pidle/prunning/...
	link        puintptr
	schedtick   uint32     // incremented on every scheduler call
	syscalltick uint32     // incremented on every system call
	sysmontick  sysmontick // last tick observed by sysmon
	m           muintptr   // back-link to associated m (nil if idle)
	mcache      *mcache
	pcache      pageCache
	raceprocctx uintptr

	deferpool    []*_defer // pool of available defer structs (see panic.go)
	deferpoolbuf [32]*_defer

	// Cache of goroutine ids, amortizes accesses to runtime·sched.goidgen.
	goidcache    uint64
	goidcacheend uint64

	// Queue of runnable goroutines. Accessed without lock.
	runqhead uint32
	runqtail uint32
	runq     [256]guintptr
	// runnext, if non-nil, is a runnable G that was ready'd by
	// the current G and should be run next instead of what's in
	// runq if there's time remaining in the running G's time
	// slice. It will inherit the time left in the current time
	// slice. If a set of goroutines is locked in a
	// communicate-and-wait pattern, this schedules that set as a
	// unit and eliminates the (potentially large) scheduling
	// latency that otherwise arises from adding the ready'd
	// goroutines to the end of the run queue.
	//
	// Note that while other P's may atomically CAS this to zero,
	// only the owner P can CAS it to a valid G.
	runnext guintptr

 ...
}
```

`runq` 是一个无锁循环队列, 由数组实现, 它的长度是 `256`, 这个长度是固定的, 不会动态调整. `runqhead` 和 `runqtail` 分别是队列的头和尾, `runqhead` 指向队列的头部, `runqtail` 指向队列的尾部. `runq` 数组的每个元素是一个 `guintptr` 类型, 它是一个 `uintptr` 类型的别名, 用来存储 `g` 的指针. 

`runq` 的操作主要是 `runqput`、`runqputslow`、`runqputbatch`、`runqget`、`runqdrain`、`runqgrab`、`runqsteal`等方法. 

接下来我们捡重点的方法看一下它是怎么实现高效额度并发读写的.

### runqput

详细代码可以参考: <https://github.com/golang/go/blob/master/src/runtime/proc.go>

`runqput` 方法是向 `runq` 中添加一个 `g` 的方法, 它是一个无锁的操作, 不会阻塞. 它的实现如下：

```go
// runqput 尝试将 g 放到本地可运行队列上。
// 如果 next 为 false，runqput 将 g 添加到可运行队列的尾部。
// 如果 next 为 true，runqput 将 g 放在 pp.runnext 位置。
// 如果可运行队列已满，runnext 将 g 放到全局队列上。
// 只能由拥有 P 的所有者执行。
func runqput(pp *p, gp *g, next bool) {
	if !haveSysmon && next {
        // 如果没有 sysmon，我们必须完全避免 runnext，否则会导致饥饿。
		next = false
	}
	if randomizeScheduler && next && randn(2) == 0 {
        // 如果随机调度器打开，我们有一半的机会避免运行 runnext
		next = false
	}

    // 如果 next 为 true，优先处理 runnext
    // 将当前的goroutine放到 runnext 中, 如果原来runnext中有goroutine, 则将其放到runq中
	if next {
	retryNext:
		oldnext := pp.runnext
		if !pp.runnext.cas(oldnext, guintptr(unsafe.Pointer(gp))) {
			goto retryNext
		}
		if oldnext == 0 {
			return
		}
		// Kick the old runnext out to the regular run queue.
		gp = oldnext.ptr()
	}

    // 重点来了，将goroutine放入runq中
retry:
	h := atomic.LoadAcq(&pp.runqhead)  // ① load-acquire, synchronize with consumers
	t := pp.runqtail
	if t-h < uint32(len(pp.runq)) { // ② 如果队列未满
		pp.runq[t%uint32(len(pp.runq))].set(gp) // ③ 将goroutine放入队列
		atomic.StoreRel(&pp.runqtail, t+1) // ④ 更新队尾
		return
	}
	if runqputslow(pp, gp, h, t) { // ⑤ 如果队列满了，调用runqputslow 尝试将goroutine放入全局队列
		return
	}
    // 如果队列未满，上面的操作应该已经成功返回，否则重试
	goto retry
}
```

`runqput` 方法的实现非常简单, 它首先判断是否需要优先处理 `runnext`, 如果需要, 就将 `g` 放到 `runnext` 中, 然后再将 `g` 放到 `runq` 中. `runq` 的操作是无锁的, 它通过 `atomic` 包提供的原子操作来实现. 这里使用的内部的更精细化的原子操作, 这个后面有时间再说. 现在大概把 ①、④ 理解为`Load`、`Store`操作即可. 

②、⑤ 分别处理本地队列未满和队列已满的情况, 如果队列未满, 就将 `g` 放到队列中, 然后更新队尾；如果队列已满, 就调用 `runqputslow` 方法, 将 `g` 放到全局队列中. 

③ 处直接将 `g` 放到队列中, 这是因为只有当前的 `P` 才能操作 `runq`, 所以不会有并发问题. 同时我们也可以看到, 我们总是往尾部插入, `t`总是一直增加的,  取余操作保证了循环队列的特性. 

`runqputslow` 会把本地队列中的一半的 `g` 放到全局队列中, 包括当前要放入的 `g`. 一旦涉及到全局队列, 就会有一定的竞争, Go 运行时使用了一把锁来控制并发, 所以 `runqputslow` 方法是一个慢路径, 是性能的瓶颈点. 

### runqputbatch

`func runqputbatch(pp *p, q *gQueue, qsize int)` 是批量往本地队列中放入 `g` 的方法, 比如它从其它 `P` 那里偷来一批 `g` , 需要放到本地队列中, 就会调用这个方法. 它的实现如下：

```go
// runqputbatch 尝试将 q 上的所有 G 放到本地可运行队列上。
// 如果队列已满，它们将被放到全局队列上；在这种情况下，这将暂时获取调度器锁。
// 只能由拥有 P 的所有者执行。
func runqputbatch(pp *p, q *gQueue, qsize int) {
	h := atomic.LoadAcq(&pp.runqhead) // ①
	t := pp.runqtail
	n := uint32(0)
	for !q.empty() && t-h < uint32(len(pp.runq)) { // ② 放入的批量goroutine非空， 并且本地队列还足以放入
		gp := q.pop()
		pp.runq[t%uint32(len(pp.runq))].set(gp)
		t++
		n++
	}
	qsize -= int(n)

	if randomizeScheduler { // ③ 随机调度器, 随机打乱
		off := func(o uint32) uint32 {
			return (pp.runqtail + o) % uint32(len(pp.runq))
		}
		for i := uint32(1); i < n; i++ {
			j := cheaprandn(i + 1)
			pp.runq[off(i)], pp.runq[off(j)] = pp.runq[off(j)], pp.runq[off(i)]
		}
	}

	atomic.StoreRel(&pp.runqtail, t) // ④ 更新队尾
	if !q.empty() {
		lock(&sched.lock)
		globrunqputbatch(q, int32(qsize))
		unlock(&sched.lock)
	}
}
```

① 获取队列头,使用原子操作获取队头. 

> “
> 
> 它下面一行是获取队尾的值, 为什么不需要使用 `atomic.LoadAcq`?
>
> 直接读取 `runqtail`, 没有使用原子操作. 是因为 `runqputbatch` 函数只能由拥有 `P` 的所有者执行，如注释所述："只能由拥有 P 的所有者执行" 这意味着在这个函数执行期间，没有其他 `goroutine` 会并发修改 `runqtail`

② 逐个的将 `g` 放到队列中, 直到放完或者放满. 

如果是随机调度器, 则使用混淆算法将队列中的 `g` 随机打乱. 

最后如果队列还有剩余的 `g`, 则调用 `globrunqputbatch` 方法, 将剩余的 `g` 放到全局队列中. 

### runqget

`runqget` 方法是从 `runq` 中获取一个 `g` 的方法, 它是一个无锁的操作, 不会阻塞. 它的实现如下：

```go
// runqget 从本地可运行队列中获取一个 G。
// 如果 inheritTime 为 true，gp 应该继承当前时间片的剩余时间。
// 否则，它应该开始一个新的时间片。
// 只能由拥有 P 的所有者执行。
func runqget(pp *p) (gp *g, inheritTime bool) {
	// If there's a runnext, it's the next G to run.
	next := pp.runnext
    // 如果有 runnext，优先处理 runnext
	if next != 0 && pp.runnext.cas(next, 0) { // ①
		return next.ptr(), true
	}

	for {
		h := atomic.LoadAcq(&pp.runqhead) // ② 获取队头, load-acquire, synchronize with other consumers
		t := pp.runqtail
		if t == h { // ③ 队列为空
			return nil, false
		}
		gp := pp.runq[h%uint32(len(pp.runq))].ptr() // ④ 获取队头的goroutine
		if atomic.CasRel(&pp.runqhead, h, h+1) { // ⑤ 更新队头 cas-release, commits consume
			return gp, false
		}
	}
}
```

① 如果有 `runnext`, 则优先处理 `runnext`, 将 `runnext` 中的 `g` 取出来. 

② 获取队列头. 如果 ③ 队列为空, 直接返回. 

④ 获取队头的 `g`, 这就是要读取的 `g`. 

⑤ 更新队头, 这里使用的是 `atomic.CasRel` 方法, 它是一个原子的 `Compare-And-Swap` 操作, 用来更新队头. 

可以看到这里只使用到了队列头 `runqhead`. 

### runqdrain

`runqdrain` 方法是从 `runq` 中获取所有的 `g` 的方法, 它是一个无锁的操作, 不会阻塞. 它的实现如下：

```go
// runqdrain 从 pp 的本地可运行队列中获取所有的 G 并返回。
// 只能由拥有 P 的所有者执行。
func runqdrain(pp *p) (drainQ gQueue, n uint32) {
	oldNext := pp.runnext
	if oldNext != 0 && pp.runnext.cas(oldNext, 0) {
		drainQ.pushBack(oldNext.ptr()) // ① 将 runnext 中的goroutine放入队列
		n++
	}

retry:
	h := atomic.LoadAcq(&pp.runqhead) // ② 获取队头, load-acquire, synchronize with other consumers
	t := pp.runqtail
	qn := t - h
	if qn == 0 {
		return
	}
	if qn > uint32(len(pp.runq)) {  // ③ 超出队列的长度了 read inconsistent h and t
		goto retry
	}

	if !atomic.CasRel(&pp.runqhead, h, h+qn) { // ④ 更新队头 cas-release, commits consume
		goto retry
	}

	// We've inverted the order in which it gets G's from the local P's runnable queue
	// and then advances the head pointer because we don't want to mess up the statuses of G's
	// while runqdrain() and runqsteal() are running in parallel.
	// Thus we should advance the head pointer before draining the local P into a gQueue,
	// so that we can update any gp.schedlink only after we take the full ownership of G,
	// meanwhile, other P's can't access to all G's in local P's runnable queue and steal them.
	// See https://groups.google.com/g/golang-dev/c/0pTKxEKhHSc/m/6Q85QjdVBQAJ for more details.
    // ⑤ 将队列中的goroutine放入队列drainQ中
	for i := uint32(0); i < qn; i++ {
		gp := pp.runq[(h+i)%uint32(len(pp.runq))].ptr()
		drainQ.pushBack(gp)
		n++
	}
	return
}
```

### runqgrab

`runqgrab` 方法是从 `runq` 中获取一半的 `g` 的方法, 它是一个无锁的操作, 不会阻塞. 它的实现如下：

```go
// runqgrab 从 pp 的本地可运行队列中获取一半的 G 并返回。
// Batch 是一个环形缓冲区，从 batchHead 开始。
// 返回获取的 goroutine 数量。
// 可以由任何 P 执行。
func runqgrab(pp *p, batch *[256]guintptr, batchHead uint32, stealRunNextG bool) uint32 {
	for {
		h := atomic.LoadAcq(&pp.runqhead) // load-acquire, synchronize with other consumers
		t := atomic.LoadAcq(&pp.runqtail) // load-acquire, synchronize with the producer
		n := t - h
		n = n - n/2 // ① 取一半的goroutine
		if n == 0 {
			if stealRunNextG {
                // ② 如果要偷取runnext中的goroutine，这里会sleep一会
				// Try to steal from pp.runnext.
				if next := pp.runnext; next != 0 {
					if pp.status == _Prunning {
						// Sleep to ensure that pp isn't about to run the g
						// we are about to steal.
						// The important use case here is when the g running
						// on pp ready()s another g and then almost
						// immediately blocks. Instead of stealing runnext
						// in this window, back off to give pp a chance to
						// schedule runnext. This will avoid thrashing gs
						// between different Ps.
						// A sync chan send/recv takes ~50ns as of time of
						// writing, so 3us gives ~50x overshoot.
						if !osHasLowResTimer {
							usleep(3)
						} else {
							// On some platforms system timer granularity is
							// 1-15ms, which is way too much for this
							// optimization. So just yield.
							osyield()
						}
					}
					if !pp.runnext.cas(next, 0) {
						continue
					}
					batch[batchHead%uint32(len(batch))] = next
					return 1
				}
			}
			return 0
		}
		if n > uint32(len(pp.runq)/2) { // ③ 如果要偷取的goroutine数量超过一半, 重试 read inconsistent h and t
			continue
		}
        // ④ 将队列中至多一半的goroutine放入batch中
		for i := uint32(0); i < n; i++ {
			g := pp.runq[(h+i)%uint32(len(pp.runq))]
			batch[(batchHead+i)%uint32(len(batch))] = g
		}
		if atomic.CasRel(&pp.runqhead, h, h+n) { // ⑤ 更新队头 cas-release, commits consume
			return n
		}
	}
}
```

① 取一半的 `g`, 这里是一个简单的算法, 取一半的 `g`. 

② 如果要偷取 `runnext` 中的 `g`, 则会尝试偷取 `runnext` 中的 `g`. 

③ 如果要偷取的 `g` 数量超过一半, 则重试. 

④ 将队列中至多一半的 `g` 放入 `batch` 中. 

⑤ 更新队头, 这里使用的是 `atomic.CasRel` 方法, 它是一个原子的 `Compare-And-Swap` 操作, 用来更新队头. 

### runqsteal

`runqsteal` 方法是从其它 `P` 的 `runq` 中偷取 `g` 的方法, 它是一个无锁的操作, 不会阻塞. 它的实现如下：

```go
// Steal half of elements from local runnable queue of p2
// and put onto local runnable queue of p.
// Returns one of the stolen elements (or nil if failed).
// runqsteal 从 p2 的本地可运行队列中偷取一半的 G 并返回.
// 如果 stealRunNextG 为 true, 它还会尝试偷取 runnext 中的 G.
func runqsteal(pp, p2 *p, stealRunNextG bool) *g {
	t := pp.runqtail
	n := runqgrab(p2, &pp.runq, t, stealRunNextG) // ① 从p2中偷取一半的goroutine
	if n == 0 {
		return nil
	}
	n--
	gp := pp.runq[(t+n)%uint32(len(pp.runq))].ptr() // ② 获取偷取的一个goroutine
	if n == 0 {
		return gp
	}
	h := atomic.LoadAcq(&pp.runqhead)  // ③ 获取队头 load-acquire, synchronize with consumers
	if t-h+n >= uint32(len(pp.runq)) { // ④ 如果队列满了，重置队列
		throw("runqsteal: runq overflow")
	}
	atomic.StoreRel(&pp.runqtail, t+n)  // ⑤ 更新队尾 store-release, makes the item available for consumption
	return gp
}
```

它实际使用了 `runqgrab` 方法来偷取 `g`, 然后再从 `runq` 中取出一个 `g`. 

以上就是`runq`的主要操作, 它针对 `Go` 调度器的特点, 设计了一套特定的队列操作的函数, 这些函数都是无锁的, 不会阻塞, 保证了高效的并发读写. 

## `gQueue` 和 `gList`
------------

`gQueue` 和 `gList` 是 Go 运行时中的两个队列, 它们都是用来存储 `g` 的, 但是它们的实现方式不同. 

`gQueue` 是一个 `G` 的双端队列, 可以从首尾增加 `gp`, 通过 g.schedlink 链接. 一个 G 只能在一个 gQueue 或 gList 上. 

```go
// A gQueue is a dequeue of Gs linked through g.schedlink. A G can only
// be on one gQueue or gList at a time.
type gQueue struct {
	head guintptr
	tail guintptr
}

// empty reports whether q is empty.
func (q *gQueue) empty() bool {
	return q.head == 0
}

// push adds gp to the head of q.
func (q *gQueue) push(gp *g) {
	gp.schedlink = q.head
	q.head.set(gp)
	if q.tail == 0 {
		q.tail.set(gp)
	}
}

// pushBack adds gp to the tail of q.
func (q *gQueue) pushBack(gp *g) {
	gp.schedlink = 0
	if q.tail != 0 {
		q.tail.ptr().schedlink.set(gp)
	} else {
		q.head.set(gp)
	}
	q.tail.set(gp)
}

// pushBackAll adds all Gs in q2 to the tail of q. After this q2 must
// not be used.
func (q *gQueue) pushBackAll(q2 gQueue) {
	if q2.tail == 0 {
		return
	}
	q2.tail.ptr().schedlink = 0
	if q.tail != 0 {
		q.tail.ptr().schedlink = q2.head
	} else {
		q.head = q2.head
	}
	q.tail = q2.tail
}

// pop removes and returns the head of queue q. It returns nil if
// q is empty.
func (q *gQueue) pop() *g {
	gp := q.head.ptr()
	if gp != nil {
		q.head = gp.schedlink
		if q.head == 0 {
			q.tail = 0
		}
	}
	return gp
}

// popList takes all Gs in q and returns them as a gList.
func (q *gQueue) popList() gList {
	stack := gList{q.head}
	*q = gQueue{}
	return stack
}
```

而 `gList` 是一个 `G` 的链表, 通过 `g.schedlink` 链接. 一个 `G` 只能在一个 `gQueue` 或 `gList` 上. 

```go
// A gList is a list of Gs linked through g.schedlink. A G can only be
// on one gQueue or gList at a time.
type gList struct {
	head guintptr
}

// empty reports whether l is empty.
func (l *gList) empty() bool {
	return l.head == 0
}

// push adds gp to the head of l.
func (l *gList) push(gp *g) {
	gp.schedlink = l.head
	l.head.set(gp)
}

// pushAll prepends all Gs in q to l.
func (l *gList) pushAll(q gQueue) {
	if !q.empty() {
		q.tail.ptr().schedlink = l.head
		l.head = q.head
	}
}

// pop removes and returns the head of l. If l is empty, it returns nil.
func (l *gList) pop() *g {
	gp := l.head.ptr()
	if gp != nil {
		l.head = gp.schedlink
	}
	return gp
}
```

这是常规的数据结构中链表的实现, 你可以和教科书中的介绍和实现做对比, 看看书本中的内容如何应用到显示的工程中的. 

## global runq
------------

详细代码: <https://github.com/golang/go/blob/master/src/runtime/runtime2.go>

一个全局的 `runq` 用来处理太多的 `goroutine`, 在本地 `runq` 中的 `goroutine` 太少的情况下, 从全局队列中偷取 `goroutine`. 主要用来处理 `P` 中 `goroutine` 不均的情况. 

因为它直接使用一把锁(`sched.lock`), 而不是 `lock-free` 的数据结构, 所以代码阅读和理解起来会相对简单一些. 这里就不详细介绍了

```go
var (
 ...
	sched      schedt
 ...
)

type schedt struct {
 ...
	// Global runnable queue.
	runq     gQueue
	runqsize int32
 ...
}

// Put gp on the global runnable queue.
// sched.lock must be held.
// May run during STW, so write barriers are not allowed.
//
//go:nowritebarrierrec
func globrunqput(gp *g) {
	assertLockHeld(&sched.lock)

	sched.runq.pushBack(gp)
	sched.runqsize++
}

// Put gp at the head of the global runnable queue.
// sched.lock must be held.
// May run during STW, so write barriers are not allowed.
//
//go:nowritebarrierrec
func globrunqputhead(gp *g) {
	assertLockHeld(&sched.lock)

	sched.runq.push(gp)
	sched.runqsize++
}

// Put a batch of runnable goroutines on the global runnable queue.
// This clears *batch.
// sched.lock must be held.
// May run during STW, so write barriers are not allowed.
//
//go:nowritebarrierrec
func globrunqputbatch(batch *gQueue, n int32) {
	assertLockHeld(&sched.lock)

	sched.runq.pushBackAll(*batch)
	sched.runqsize += n
	*batch = gQueue{}
}

// Try get a batch of G's from the global runnable queue.
// sched.lock must be held.
func globrunqget(pp *p, max int32) *g {
	assertLockHeld(&sched.lock)

	if sched.runqsize == 0 {
		return nil
	}

	n := sched.runqsize/gomaxprocs + 1
	if n > sched.runqsize {
		n = sched.runqsize
	}
	if max > 0 && n > max {
		n = max
	}
	if n > int32(len(pp.runq))/2 {
		n = int32(len(pp.runq)) / 2
	}

	sched.runqsize -= n

	gp := sched.runq.pop()
	n--
	for ; n > 0; n-- {
		gp1 := sched.runq.pop()
		runqput(pp, gp1, false)
	}
	return gp
}
```
