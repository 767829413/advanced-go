# Golang 语⾔的内存管理与垃圾回收

## 基本概念

`提出问题`

 ![golang_deep_0.jpg](https://s2.loli.net/2023/06/13/frV7n4zF1NtIpTS.png)

`栈上内存分配`

 **栈分配，函数调⽤返回后，函数栈帧⾃动销毁(SP 下移)**

 演示动画: <https://www.figma.com/proto/tSl3CoSWKitJtvIhqLd8Ek/memory-management-%26%26-garbage-collection?node-id=117-202&starting-point-node-id=117%3A21>

`堆分配`

 在 C 语⾔中，返回函数的局部变量会怎么样？

  ```C
  int *func(void)
  {
    int num = 123;
    /* ... */
    return &num;
  }
  ```

 上面的可能导致: 

    悬垂指针：Dangling pointer
    可能触发：Segmentation fault!!!!

 堆分配，在 Go 语⾔⾥，为什么我们不⽤担⼼ dangling pointer? 

 主要是go语言编译器做了 Escape analysis (逃逸分析)

 ```go
 package main

 func main() {
    println(getNP())
 }
 
 func getNP() *int {
    var n = 1234
    return &n
 }
 ```

 ```bash
 # -m 越多,给出的分析结果越多
 go build -o test -gcflags="-m -m" ./main.go 
 # command-line-arguments
 ./main.go:7:6: can inline getNP with cost 8 as: func() *int { var n int; n = 1234; return &n }
 ./main.go:3:6: can inline main with cost 11 as: func() { println(getNP()) }
 ./main.go:4:15: inlining call to getNP
 ./main.go:8:6: n escapes to heap:
 ./main.go:8:6:   flow: ~r0 = &n:
 ./main.go:8:6:     from &n (address-of) at ./main.go:9:9
 ./main.go:8:6:     from return &n (return) at ./main.go:9:2
 ./main.go:8:6: moved to heap: n
 ```

`Escape analysis (逃逸分析)`

 golang中逃逸分析的代码可以在去这里看: <https://github.com/golang/go/blob/master/src/cmd/compile/internal/escape/escape.go>

* 逃逸分析是一种优化技术，用于确定哪些Go变量可以在堆栈上分配，以减少内存分配和垃圾回收成本. 
* 逃逸分析需要遵守两个关键规则：不允许在堆中存储指向堆栈对象的指针，不允许指向堆栈对象的指针存在于对象已销毁的情况下. 
* 在Go语言中，逃逸分析通过对抽象语法树（AST）进行静态数据流分析来实现. 建立带权重的有向图，遍历图，查找可能违反上述不变式的赋值路径. 
* EA（Escape Analysis）将每个分配语句或表达式映射到唯一的“位置”，并将每个赋值建模为位置之间的有向边. EA（Escape Analysis）记录每个边的解引用操作数减去寻址操作数作为边的权重. 
* 为支持跨程序的分析，EA还记录堆到每个函数的参数和结果参数的数据流，并将此信息总结为“参数标签”，用于执行静态调用点的函数参数的逃逸分析. 

 还可以通过关于 escape.go 的单元测试来看: <https://github.com/golang/go/tree/master/test>

 ![escape_test.png](https://s2.loli.net/2023/06/13/eakyInRu8NDMf4V.png)

`内存管理中的⻆⾊`

*内存需要分配，谁来分配：*

* ⾃动 allocator，⼿⼯分配

*内存需要回收，谁来回收：*

* ⾃动 collector，⼿⼯回收

**⾃动内存回收技术=垃圾回收技术**

`内存管理中的三个⻆⾊`

* Mutator：fancy(花哨的) word for application，其实就是你写的应⽤程序，它会不断地修改对象的引⽤关系，即对象图. 
* Allocator：内存分配器，负责管理从操作系统中分配出的内存空间，C语言中 malloc 其实底层就有⼀个内存分配器的实现(glibc 中)，tcmalloc 是 malloc 多线程改进版. Go 中的实现类似 tcmalloc. 
* Collector：垃圾收集器，负责清理死对象，释放内存空间. 

*Mutator、Allocator、Collector 概览：*
 
 ![role_mm.png](https://s2.loli.net/2023/06/13/Vrd4qfFLPbzYM3l.png)

`内存管理抽象`

每个操作系统都有相应的实现,如：

 ![mm_abs.png](https://s2.loli.net/2023/06/13/mSt3WfiQGLIhr2J.png)

mem_linux.go, mem_windows.go, 主要逻辑可以去看 <https://github.com/golang/go/blob/master/src/runtime/malloc.go>

## 进程虚拟内存布局

`内存的分配类型`

 ![mm_type.png](https://s2.loli.net/2023/06/13/z61nQKdFu8G7MVb.png)

 参考链接: <https://www.kernel.org/doc/html/next/x86/x86_64/mm.html>

 **进程虚拟内存分布: 多线程的情况**

 ![mm_muilty.png](https://s2.loli.net/2023/06/13/nOyVebE5LPJUCqu.png)

## Allocator 基础

`介绍`

* Bump/Sequential Allocator (线性分配器)
  * 演示动画: <https://www.figma.com/proto/tSl3CoSWKitJtvIhqLd8Ek/memory-management-%26%26-garbage-collection?node-id=175-119&starting-point-node-id=175%3A119>
* Free List Allocator (空闲链表分配器)
  * First-Fit
  * Next-Fit
  * Best-Fit
  * Segregated-Fit
  * 演示动画: <https://www.figma.com/proto/tSl3CoSWKitJtvIhqLd8Ek/memory-management-%26%26-garbage-collection?node-id=233-22&starting-point-node-id=233%3A22>

## malloc 实现

`简介`

*当你执⾏ malloc 时：*

 ![malloc.png](https://s2.loli.net/2023/06/13/2DyLXjUdh1sM8fb.png)

* brk 只能通过调整 program break位置推动堆增⻓
* mmap 可以从任意未分配位置映射内存

![malloc_1.png](https://s2.loli.net/2023/06/13/6vX5j7kR4aSr3Gm.png)
![malloc_2.png](https://s2.loli.net/2023/06/13/EOj27nNikJgBeHM.png)

*为什么⼈⾁管理内存不靠谱，复杂对象图维护时的 dangling pointer：*

演示动画: <https://www.figma.com/proto/tSl3CoSWKitJtvIhqLd8Ek/memory-management-%26%26-garbage-collection?node-id=175-17&starting-point-node-id=175%3A17>

## Go 语⾔内存分配

`版本演变`

*⽼版本，连续堆：*
 
 ![mm-1.png](https://s2.loli.net/2023/06/13/oFpBaWMwOnsb1RV.png)

*新版本，稀疏堆：*

 ![mm-2.png](https://s2.loli.net/2023/06/13/ZrKouasSwxf6WFe.png)

 *申请稀疏堆时，我们该⽤ brk 还是 mmap ? brk分配的都是连续地址,这里应该是选择mmap*

`Heap grow`

 演示动画: <https://www.figma.com/proto/tSl3CoSWKitJtvIhqLd8Ek/memory-management-%26%26-garbage-collection?node-id=151-306&starting-point-node-id=151%3A38>

`分配`

*分配⼤⼩分类：*

* Tiny : size < 16 bytes && has no pointer(noscan)
* Small ：has pointer(scan) || (size >= 16 bytes && size <= 32 KB)
* Large : size > 32 KB

*内存分配器在 Go 语⾔中维护了⼀个多级结构：*

* mcache -> mcentral -> mheap
* mcache：与 P 绑定，本地内存分配操作，不需要加锁. 
* mcentral：中⼼分配缓存，分配时需要上锁，不同 spanClass 使⽤不同的锁. 
* mheap：全局唯⼀，从 OS 申请内存，并修改其内存定义结构时，需要加锁，是个全局锁. 

 **sizeClass 分类(sizeClass = spanClass >> 1)**
 
 <https://github.com/golang/go/blob/master/src/runtime/sizeclasses.go>

`Tiny alloc`

 演示动画: <https://www.figma.com/proto/tSl3CoSWKitJtvIhqLd8Ek/memory-management-%26%26-garbage-collection?node-id=165-197&starting-point-node-id=165%3A197>

`Small alloc`

 ![mm-3.png](https://s2.loli.net/2023/06/13/ELnMYQcOIJ4Flse.png)

`Large alloc`

* ⼤对象分配会直接越过 mcache、mcentral，直接从 mheap 进⾏相应数量的 page 分配. 
* pageAlloc 结构经过多个版本的变化，从：freelist -> treap -> radix tree，查找时间复杂度越来越低，结构越来越复杂

`补充`

**Refill 流程：**

* 本地 mcache 没有时触发(mcache.refill)
* 从 mcentral ⾥的 non-empty 链表中找(mcentral.cacheSpan)
* 尝试 sweep mcentral 的 empty，insert sweeped -> nonempty(mcentral.cacheSpan)
* 增⻓ mcentral，尝试从 arena 获取内存(mcentral.grow)
* arena 如果还是没有，向操作系统申请(mheap.alloc) 最终还是会将申请到的 mspan 放在 mcache 中

`数据结构总览`
 
 ![mm-4.png](https://s2.loli.net/2023/06/13/N5Q8XADwKbkGluO.png)

`Bitmap 与 allocCache`

 ![mm-5.png](https://s2.loli.net/2023/06/13/TZRUjWn28POqSCx.png)

## 垃圾回收基础

`垃圾分类`

* 语义垃圾(semantic garbage)—有的被称作内存泄露
  * 语义垃圾指的是从语法上可达(可以通过局部、全局变量引⽤得到)的对象，但从语义上来讲他们是垃圾，垃圾回收器对此⽆能为⼒. 
* 语法垃圾(syntactic garbage)
  * 语法垃圾是讲那些从语法上⽆法到达的对象，这些才是垃圾收集器主要的收集⽬标. 

`语义垃圾(semantic garbage)`

 动画演示: <https://www.figma.com/proto/tSl3CoSWKitJtvIhqLd8Ek/memory-management-%26%26-garbage-collection?node-id=185-49&starting-point-node-id=185%3A13>

`语法垃圾(syntactic garbage)`

 ```go
 package main
 
 func main() {
    allocOnHeap()
 }
 
 func allocOnHeap() {
    var m = make([]int, 10240)
    println(m)
 }
 ```
 
 ```bash
 ➜  advanced-go git:(main) ✗ go run -gcflags="-m" ./main.go 
 # command-line-arguments
 ./main.go:7:6: can inline allocOnHeap
 ./main.go:3:6: can inline main
 ./main.go:4:13: inlining call to allocOnHeap
 ./main.go:4:13: make([]int, 10240) escapes to heap
 ./main.go:8:14: make([]int, 10240) escapes to heap
 [10240/10240]0xc000100000
 ```

`常⻅垃圾回收算法`

* 引⽤计数(Reference Counting)：某个对象的根引⽤计数变为 0 时，其所有⼦节点均需被回收. 
* 标记压缩(Mark-Compact)：将存活对象移动到⼀起，解决内存碎⽚问题. 
* 复制算法(Copying)：将所有正在使⽤的对象从 From 复制到 To 空间，堆利⽤率只有⼀半. 
* 标记清扫(Mark-Sweep)：解决不了内存碎⽚问题. 需要与能尽量避免内存碎⽚的分配器使⽤，如 tcmalloc.  <— Go 在这⾥

 演示动画和说明: <https://spin.atomicobject.com/2014/09/03/visualizing-garbage-collection-algorithms/>

## Go 语⾔垃圾回收

`演变`

 *1.8之前*

  ![mm-1.png](https://s2.loli.net/2023/06/13/9jJBpCziGuI3kZS.png)

 **1.8 后通过混合 write barrier 消除了第⼆个 stw 中的 stack re-scan，stw 时间⼤⼤减少**

  ![mm-2.png](https://s2.loli.net/2023/06/13/qxg6jQIeTtfdJ2m.png)

`程序⼊⼝ && 触发点`

 **垃圾回收⼊⼝：gcStart**

  ![mm-3.png](https://s2.loli.net/2023/06/13/4kXdtsfRizKGIcT.png)

## GC 标记流程

`基本代码流程`
 
 ![gc.svg](./gc.svg)

*标记对象从哪⾥来？*

* gcMarkWorker
* Mark assist
* mutator write/delete heap pointers

*标记对象到哪⾥去？*

* Work buffer(gcMarkWorker)
  * 本地 work buffer => p.gcw
  * 全局 work buffer => runtime.work.full
* Write barrier buffer(mutator write/delete heap pointers)
  * 本地 write barrier buffer => p.wbBuf

![mm-4.png](https://s2.loli.net/2023/06/13/CaSj7Qy5KxTnHWq.png)

*标记对象被谁消费？*
 
* gcMarkWorker

`三⾊标记抽象与解释`

* ⿊：已经扫描完毕，⼦节点扫描完毕. (gcmarkbits = 1，且在队列外. )
* 灰：已经扫描完毕，⼦节点未扫描完毕. (gcmarkbits = 1, 在队列内)
* ⽩：未扫描，collector 不知道任何相关信息. 

动画演示: <https://www.figma.com/proto/tSl3CoSWKitJtvIhqLd8Ek/memory-management-%26%26-garbage-collection?node-id=2-38&starting-point-node-id=2%3A38>

![mm-5.png](https://s2.loli.net/2023/06/13/zVAc6IdywmlgShQ.png)

1. 首先把所有的对象都放到白色的集合中
2. 从根节点开始遍历对象，遍历到的白色对象从白色集合中放到灰色集合中
3. 遍历灰色集合中的对象，把灰色对象引用的白色集合的对象放入到灰色集合中，同时把遍历过的灰色集合中的对象放到黑色的集合中
4. 循环步骤 3，知道灰色集合中没有对象
5. 步骤 4 结束后，白色集合中的对象就是不可达对象，也就是垃圾，进行回收

`⼀些问题`

* 对象在标记过程中不能丢失
* Mark 阶段 mutator 的指向堆的指针修改需要被记录下来
* GC Mark 的 CPU 控制要努⼒做到 25% 以内

**三⾊抽象的问题，标记过程中对象漏标，导致被意外回收：**

动画演示: <https://www.figma.com/proto/tSl3CoSWKitJtvIhqLd8Ek/memory-management-%26%26-garbage-collection?node-id=16-151&starting-point-node-id=16%3A151>

`解决丢失问题的理论基础`

* 强三⾊不变性 strong tricolor invariant
  * 禁⽌⿊⾊对象指向⽩⾊对象

  ![mm-6.png](https://s2.loli.net/2023/06/13/dfRgznOYoyGmC2p.png)

* 弱三⾊不变性 weak tricolor invariant
  * ⿊⾊对象指向的⽩⾊对象，如果有灰⾊对象到它的可达路径，那也可以

  ![mm-7.png](https://s2.loli.net/2023/06/13/SCb1aBeKtWLPEJ9.png)
  
`write barrier`

**barrier 本质是: snippet of code insert before pointer modify**

 ![mm-8.png](https://s2.loli.net/2023/06/13/Lvl8G5wTI7rei6o.png)
 ![mm-9.png](https://s2.loli.net/2023/06/13/JTBDkjLKyOfzuZF.png)

*这就是 write barrier，是在指针修改前插⼊的⼀个函数调⽤*

`Dijkstra barrier(插入写屏障)`

```go
// slot is the destination in Go code.
// ptr is the value that goes into the slot in Go code.
// Slot 是 Go 代码⾥的被修改的指针对象
// Ptr 是 Slot 要修改成的值

// Dijkstra barrier(插入写屏障)
// 新的对象进行标灰
// 实现的是强三色不变性
func DijkstraWB(slot *unsafe.Pointer, ptr unsafe.Pointer) {
    shade(ptr)
    *slot = ptr
}
```

**优点:**

* 能够保证堆上对象的强三⾊不变性(⽆栈对象参与时)
* 能防⽌指针从栈被隐藏进堆(因为堆上新建的连接都会被着⾊)

**缺点:**

* 不能防⽌栈上的⿊⾊对象指向堆上的⽩⾊对象(这个⽩⾊对象之前是被堆上的⿊/灰指着的)
* 所以在 mark 结束后需要 stw 重新扫描所有 goroutine 栈

 动画演示: <https://www.figma.com/proto/tSl3CoSWKitJtvIhqLd8Ek/memory-management-%26%26-garbage-collection?node-id=201-1&starting-point-node-id=201%3A1>

`Yuasa barrier(删除写屏障)`

```go
// Yuasa barrier(删除写屏障)
// 旧的对象进行标灰
// 实现的是弱三色不变性
func YuasaWB(slot *unsafe.Pointer, ptr unsafe.pointer) {
    shade(*slot)
    *slot = ptr
}
```

**优点:**

* 能够保证堆上的弱三⾊不变性(⽆栈对象参与时)
* 能防⽌指针从堆被隐藏进栈(因为堆上断开的连接都会被着⾊)

**缺点:**

* 不能防⽌堆上的⿊⾊对象指向堆上的⽩⾊对象(这个⽩⾊对象之前是由栈的⿊/灰⾊对象指着的)
* 所以需要 GC 开始时 STW 对栈做快照

 动画演示: <https://www.figma.com/proto/tSl3CoSWKitJtvIhqLd8Ek/memory-management-%26%26-garbage-collection?node-id=201-294&starting-point-node-id=201%3A294>

`Hybrid barrier(混合写屏障)`

**如果我们在所有指针操作中都加上 Dijkstra barrier 或者 Yuasa barrier，就可以避免对象丢失了，为啥实际的实现没那么简单？**

* 因为栈的操作频率极⾼，所以 Go 在栈上指针操作上是不加 barrier 的. 
  * 因为 Go 的栈上的指针编辑不加 barrier，所以单独使⽤任意⼀种 barrier 都会有问题
    * Dijkstra barrier(插入写屏障) 的问题
      * 会出现栈上⿊指向堆上⽩的情况,该⽩⾊对象之前被堆上对象所指
    * Yuasa barrier(删除写屏障) 的问题
      * 会出现堆上⿊指向堆上⽩的情况,该⽩⾊对象之前被栈上某对象所指

*为了解决插入写屏障和删除写屏障面对的问题(栈上指针操作上是不加 barrier)*

```go
// Hybrid barrier(混合写屏障)
func HybridWB(slot *unsafe.Pointer, ptr unsafe.pointer) {
    shade(*slot)
    // 检查当前栈是不是灰色的
    if current stack is grey
        shade(ptr)
    *slot = ptr
}
```

**不过golang中实际使用的是是另一种方法,主要原因是栈检查的成本太高,不利于性能**

```go
// 混合屏障会将指针(指向堆的)修改前指向的位置和修改后指向的位置都标灰
func RealityWB(slot *unsafe.Pointer, ptr unsafe.pointer) {
    // 没有进行栈的检查
    shade(*slot)
    shade(ptr)
    *slot = ptr
}
```

Go 实际的混合写屏障：代码在 gcWriteBarrier 汇编函数中: <https://github.com/golang/go/blob/master/src/runtime/asm_386.s>

```go
// gcWriteBarrier informs the GC about heap pointer writes.
//
// gcWriteBarrier returns space in a write barrier buffer which
// should be filled in by the caller.
// gcWriteBarrier does NOT follow the Go ABI. It accepts the
// number of bytes of buffer needed in DI, and returns a pointer
// to the buffer space in DI.
// It clobbers FLAGS. It does not clobber any general-purpose registers,
// but may clobber others (e.g., SSE registers).
// Typical use would be, when doing *(CX+88) = AX
//     CMPL    $0, runtime.writeBarrier(SB)
//     JEQ     dowrite
//     CALL    runtime.gcBatchBarrier2(SB)
//     MOVL    AX, (DI)
//     MOVL    88(CX), DX
//     MOVL    DX, 4(DI)
// dowrite:
//     MOVL    AX, 88(CX)
TEXT gcWriteBarrier<>(SB),NOSPLIT,$28
    // Save the registers clobbered by the fast path. This is slightly
    // faster than having the caller spill these.
    MOVL	CX, 20(SP)
    MOVL	BX, 24(SP)
retry:
    // TODO: Consider passing g.m.p in as an argument so they can be shared
    // across a sequence of write barriers.
    get_tls(BX)
    MOVL	g(BX), BX
    MOVL	g_m(BX), BX
    MOVL	m_p(BX), BX
    // Get current buffer write position.
    MOVL	(p_wbBuf+wbBuf_next)(BX), CX	// original next position
    ADDL	DI, CX				// new next position
    // Is the buffer full?
    CMPL	CX, (p_wbBuf+wbBuf_end)(BX)
    JA	flush
    // Commit to the larger buffer.
    MOVL	CX, (p_wbBuf+wbBuf_next)(BX)
    // Make return value (the original next position)
    SUBL	DI, CX
    MOVL	CX, DI
    // Restore registers.
    MOVL	20(SP), CX
    MOVL	24(SP), BX
    RET
```

`垃圾回收代码流程`

* gcStart -> gcBgMarkWorker && gcRootPrepare，这时 gcBgMarkWorker 在休眠中
* schedule -> findRunnableGCWorker 唤醒适宜数量的 gcBgMarkWorker
* gcBgMarkWorker -> gcDrain -> scanobject -> greyobject(set mark bit and put to gcw)
* 在 gcBgMarkWorker 中调⽤ gcMarkDone 排空各种 wbBuf 后，使⽤分布式 termination 检查算法，进⼊ gcMarkTermination -> gcSweep 唤醒后台沉睡的 sweepg 和 scvg -> sweep -> wake bgsweep && bgscavenge

`CPU 使⽤控制`

* GC 的 CPU 控制⽬标是整体 25%
* 当 P = 4 * N 时，只要启动 N 个 worker 就可以. 
* 但 P ⽆法被 4 整除时，需要兼职的 gcMarkWorker 来帮助做⼀部分⼯作
  * 全职 GC 员⼯：Dedicated worker，需要⼀直⼲活，直到被抢占. 
  * 兼职 GC 员⼯：Fractional worker，达到业绩⽬标(fractionalUtilizationGoal)时可以主动让出. 
  * 还有⼀种 IDLE 模式，在调度循环中发现找不到可执⾏的 g，但此时有标记任务未完成，就⽤开启 IDLE 模式去帮忙. 
* Worker 运⾏模式在：\_p\_.gcMarkWorkerMode

 ![mm-11.png](https://s2.loli.net/2023/06/13/3eKk4HwQgd58ExM.png)

## Go 语⾔的栈内存管理

 **栈本身的内存：**

  ![mm-12.png](https://s2.loli.net/2023/06/13/aZADGI3dXCU4FM1.png)

* newstack
* shrinkstack

* 使⽤ allocManual 和 freeManual 相当于⼿动管理内存，不计⼊ heap_inuse 和 heap_sys；
* 计⼊ stackinuse 和 stacksys
* 栈上变量的内存：SP 移动销毁，简单快速

## References

* Memory Management Reference: <https://www.memorymanagement.org/>
* <https://my.eng.utah.edu/~cs4400/malloc.pdf>
* <https://cboard.cprogramming.com/linux-programming/101090-what-differences-between-brk-mmap.html>
* <https://medium.com/a-journey-with-go/go-memory-management-and-memory-sweep-cc71b484de05>
* <https://medium.com/a-journey-with-go/go-memory-management-and-allocation-a7396d430f44>
* <https://spin.atomicobject.com/2014/09/03/visualizing-garbage-collection-algorithms/>
* <https://golangpiter.com/system/attachments/files/000/001/718/original/GP_2019_-_An_Insight_Into_Go_Garbage_Collection.pdf?1572419303>
* <https://www.cnblogs.com/zkweb/p/7880099.html?utm_source=tuicool&utm_medium=referral>
* <https://go.googlesource.com/proposal/+/master/design/17503-eliminate-rescan.md>
* <https://blog.golang.org/ismmkeynote>
* <https://www.ardanlabs.com/blog/2018/12/garbage-collection-in-go-part1-semantics.html>
* <https://docs.google.com/document/d/1wmjrocXIWTr1JxU-3EQBI6BK6KgtiFArkG47XK73xIQ/edit#>
* <https://zhuanlan.zhihu.com/p/95056679>
* <https://github.com/golang-design/under-the-hood>

## 未涉及

* Runtime 中 Not in heap 类对象的内存管理
* Page alloc && Page cache
* Stack alloc && persistentAlloc
* Tiny 分配时的 offset 对⻬逻辑
* 为防⽌⻓时间标记导致应⽤请求阻塞，引⼊的 oblet 概念
* 对象被回收时的 finalizer 执⾏逻辑
* Fd 类内置数据结构的引⽤计数逻辑
* Sweep 和 Scvg 流程未涉及

**菜狗要有自知之明!!!**
**菜狗要有自知之明!!!**
**菜狗要有自知之明!!!**