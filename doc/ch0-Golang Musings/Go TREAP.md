# Golang 互斥锁(Mutex)杂谈之Treap

## Treap 介绍

`Treap` 是一棵二叉树, 它同时维护 `二叉搜索树 (BST)` 和 `堆(heap)` 的属性, 所以由此得名 `(tree + heap   ⇒  treap)`. 

从形式上讲, `Treap  (tree + heap)`  是一棵二叉树, 其节点包含两个值, 一个  `key`  和一个  `priority`, 这样 `key` 保持 `BST` 属性, `priority` 是一个保持 `heap` 属性的随机值 (至于是最大堆还是最小堆并不重要). 相对于其他的平衡二叉搜索树, `Treap` 的特点是实现简单, 且能基本实现随机平衡的结构. 属于弱平衡树. 

`Treap` 由 `Raimund Siedel` 和 `Cecilia Aragon` 于 `1989` 年提出. 

具体来说, `Treap` 是一种在二叉树中存储键值对 `(X,Y)` 的数据结构, 其特点是：按 `X` 值满足二叉搜索树的性质, 同时按 `Y` 值满足二叉堆的性质. 如果树中某个节点包含值 `(X₀,Y₀)`, 那么：

* 左子树中所有节点的 `X` 值都满足 `X ≤ X₀` (BST 属性)
    
* 右子树中所有节点的 `X` 值都满足 `X₀ ≤ X` (BST 属性)
    
* 左右子树中所有节点的 `Y` 值都满足 `Y ≤ Y₀`  (堆属性. 这里以最大堆为例) 
    
在这种实现中, `X` 是键 (同时也是存储在 `Treap` 中的值) , 并且 `Y` 称为**优先级**. 如果没有优先级, 则 `Treap` 将是一个常规的二叉搜索树. 

优先级 (前提是每个节点的优先级都不相同) 的特殊之处在于：它们可以确定性地决定树的最终结构 (不会受到插入数据顺序的影响) . 这一点是可以通过相关定理来证明的. 这里有个巧妙的设计：如果我们随机分配这些优先级值, 就能在平均情况下得到一棵比较平衡的树 (避免树退化成链表) . 这样就能保证主要操作 (如查找、插入、删除等) 的时间复杂度保持在 `O(log N)` 水平. 正是因为这种随机分配优先级的特点, 这种数据结构也被称为"随机二叉搜索树". 

```bash
       (4,7)
      /     \
  (2,4)     (6,5)
   /  \     /   \
(1,1) (3,3) (5,2) (7,6)
```

Treap 维护堆性质的方法用到了旋转, 且只需要进行两种旋转操作, 因此编程复杂度较红黑树、AVL 树要小一些. 

### 插入

给节点随机分配一个优先级, 先和二叉搜索树的插入一样, 先把要插入的点插入到一个叶子上, 然后跟维护堆一样进行以下操作：

1. 如果当前节点的优先级比父节点大就进行 2. 或 3. 的操作
    
2. 如果当前节点是父节点的左子叶就右旋
    
3. 如果当前节点是父节点的右子叶就左旋. 

假设我们要插入键值为3，优先级为9的新节点到初始Treap中:

```bash
# 初始：
    (4,7)
   /     \
(2,4)   (6,5)

# 插入后：
    (4,7)
   /     \
(2,4)   (6,5)
   \
   (3,9)

# 第一次右旋：
    (4,7)
   /     \
(3,9)   (6,5)
 /
(2,4)

# 最终右旋：
    (3,9)
   /     \
(2,4)   (4,7)
          \
         (6,5)
```
    
### 删除

因为 `Treap` 满足堆性质, 所以只需要把要删除的节点旋转到叶节点上, 然后直接删除就可以了. 具体的方法就是每次找到优先级最大的子叶, 向与其相反的方向旋转, 直到那个节点被旋转到了叶节点, 然后直接删除. 

1. 如果要删除的节点一开始就是叶子节点，我们可以直接删除它，不需要旋转。
2. 旋转的方向取决于子节点的优先级：
    * 如果左子节点优先级更高，我们进行右旋
    * 如果右子节点优先级更高，我们进行左旋
3. 这个过程保证了删除操作不会破坏 `Treap` 的 `BST` 和 `堆` 的性质。
4. 删除操作的平均时间复杂度是`O(log n)`，因为旋转的次数预期是树的高度。

假设我们有以下 `Treap`，我们要删除键值为4的节点：

```bash
# 初始
     (4,7)
    /     \
 (2,5)   (6,3)
 /   \
(1,2) (3,1)

# 右旋:
     (2,5)
    /     \
 (1,2)   (4,7)
        /     \
     (3,1)   (6,3)

# 继续右旋:
     (2,5)
    /     \
 (1,2)   (3,1)
            \
           (4,7)
              \
             (6,3)

# 左旋:
     (2,5)
    /     \
 (1,2)   (3,1)
            \
           (6,3)
           /
         (4,7)

# 移除叶子节点:
     (2,5)
    /     \
 (1,2)   (3,1)
            \
           (6,3)
```

### 查找

和一般的二叉搜索树一样, 但是由于 `Treap` 的随机化结构, `Treap` 中查找的期望复杂度是 `O(logn)`

查找步骤：

1. 从根节点开始。
2. 比较当前节点的键值与目标键值。
3. 如果相等，则找到目标节点。
4. 如果目标键值小于当前节点，则移动到左子树。
5. 如果目标键值大于当前节点，则移动到右子树。
6. 重复步骤2-5，直到找到目标节点或达到叶子节点。

设我们有以下 `Treap` 结构，我们要查找键值为5的节点：

```bash
# 目标键值5 > 4，移动到右子树
       (4,7)  <-- 当前节点
      /     \
  (2,5)     (6,3)
   /  \     /   \
(1,2) (3,1) (5,6) (7,4)


# 目标键值5 < 6，移动到左子树
       (4,7)
      /     \
  (2,5)     (6,3)  <-- 当前节点
   /  \     /   \
(1,2) (3,1) (5,6) (7,4)


# 目标键值5 = 5，查找成功
       (4,7)
      /     \
  (2,5)     (6,3)
   /  \     /   \
(1,2) (3,1) (5,6) (7,4)
             ^
             |
        当前节点（目标找到）
```

### 更多参考

以上是 `Treap` 数据结构的背景知识, 如果你想了解更多关于 `Treap` 的知识, 你可以参考

* https://en.wikipedia.org/wiki/Treap
    
* https://medium.com/carpanese/a-visual-introduction-to-treap-data-structure-part-1-6196d6cc12ee
    
* https://cp-algorithms.com/data_structures/treap.html
    
## Go 运行时的 Treap 和用途
-----------------

在 Go 运行时 [sema.go#semaRoot](https://github.com/golang/go/blob/master/src/runtime/sema.go#L40) 中, 定义了一个数据结构 `semaRoot`:

```go
type semaRoot struct {
	lock  mutex
	treap *sudog        // 不重复的等待者(goroutine)的平衡树(treap)的根节点
	nwait atomic.Uint32 // 等待者(goroutine)的数量
}

type sudog struct {
	// The following fields are protected by the hchan.lock of the
	// channel this sudog is blocking on. shrinkstack depends on
	// this for sudogs involved in channel ops.

	g *g

	next *sudog
	prev *sudog
	elem unsafe.Pointer // data element (may point to stack)

	// The following fields are never accessed concurrently.
	// For channels, waitlink is only accessed by g.
	// For semaphores, all fields (including the ones above)
	// are only accessed when holding a semaRoot lock.

	acquiretime int64
	releasetime int64
	ticket      uint32

	// isSelect indicates g is participating in a select, so
	// g.selectDone must be CAS'd to win the wake-up race.
	isSelect bool

	// success indicates whether communication over channel c
	// succeeded. It is true if the goroutine was awoken because a
	// value was delivered over channel c, and false if awoken
	// because c was closed.
	success bool

	// waiters is a count of semaRoot waiting list other than head of list,
	// clamped to a uint16 to fit in unused space.
	// Only meaningful at the head of the list.
	// (If we wanted to be overly clever, we could store a high 16 bits
	// in the second entry in the list.)
	waiters uint16

	parent   *sudog // semaRoot binary tree
	waitlink *sudog // g.waiting list or semaRoot
	waittail *sudog // semaRoot
	c        *hchan // channel
}
```

这是 `Go` 语言互斥锁(Mutex)底层实现中的关键数据结构, 用于管理等待获取互斥锁的 `goroutine` 队列. 我们已经知道, 在获取 `sync.Mutex` 时, 如果锁已经被其它 `goroutine` 获取, 那么当前请求锁的 `goroutine` 会被 `block` 住, 就会被放入到这样一个数据结构中 (所以你也知道这个数据结构中的 `goroutine` 都是唯一的, 不重复). 

`semaRoot` 保存了一个平衡树, 树中的 `sudog` 节点都有不同的地址 `(s.elem)` ,每个 `sudog` 可能通过 `s.waitlink` 指向一个链表, 该链表包含等待相同地址的其他 `sudog`. 对具有相同地址的 `sudog` 内部链表的操作时间复杂度都是 `O(1)`. 扫描顶层 `semaRoot` 列表的时间复杂度是 `O(log n)`,其中 `n` 是具有被阻塞 `goroutine` 的不同地址的数量 (这些地址会散列到给定的 `semaRoot`) . 

`semaRoot` 的 `Treap *sudog` 其实就是一个 `Treap`, 我们来看看它的实现. 

## 增加一个元素 (入队) 
-----------------

增加一个等待的 goroutine(`sudog`) 到 `semaRoot` 的 `Treap` 中, 如果 `lifo` 为 `true`, 则将 `s` 替换到 `t` 的位置, 否则将 `s` 添加到 `t` 的等待列表的末尾. 

```go
func (root *semaRoot) queue(addr *uint32, s *sudog, lifo bool) {
   // 设置这个要加入的节点
	s.g = getg()
	s.elem = unsafe.Pointer(addr)
	s.next = nil
	s.prev = nil
	s.waiters = 0

	var last *sudog
	pt := &root.treap
   // 从根节点开始
	for t := *pt; t != nil; t = *pt { // ①
      // 如果地址已经在列表中,则加入到这个地址的链表中
		if t.elem == unsafe.Pointer(addr) {
			// 如果地址已经在列表中，并且指定了先入后出flag,这是一个替换操作
			if lifo {
				// 替换操作
				*pt = s
				s.ticket = t.ticket
            ... // 把t的各种信息复制给s
			} else {
				// 增加到到等待列表的末尾
				if t.waittail == nil {
					t.waitlink = s
				} else {
					t.waittail.waitlink = s
				}
				t.waittail = s
				s.waitlink = nil
				if t.waiters+1 != 0 {
					t.waiters++
				}
			}
			return
		}
		last = t
      // 二叉搜索树查找
		if uintptr(unsafe.Pointer(addr)) < uintptr(t.elem) { // ②
			pt = &t.prev
		} else {
			pt = &t.next
		}
	}

	// 为新节点设置ticket.这个ticket是一个随机值，作为随机堆的优先级，用于保持treap的平衡。
	s.ticket = cheaprand() | 1 // ③
	s.parent = last
	*pt = s

	// 根据优先级(ticket)旋转以保持treap的平衡
	for s.parent != nil && s.parent.ticket > s.ticket { // ④
		if s.parent.prev == s {
			root.rotateRight(s.parent) // ⑤
		} else {
			if s.parent.next != s {
				panic("semaRoot queue")
			}
			root.rotateLeft(s.parent) // ⑥
		}
	}
}
```

① 是遍历 `Treap` 的过程, 当然它是通过搜索二叉树的方式实现. `addr` 就是我们一开始讲的 `Treap` 的 `key`, 也就是 `s.elem`. 首先检查 `addr` 已经在 `Treap` 中, 如果存在, 那么就把 `s` 加入到 `addr` 对应的 `sudog` 链表中, 或者替换掉 `addr` 对应的 `sudog`. 

这个 `addr`, 如果对于 `sync.Mutex` 来说, 就是 `Mutex.sema` 字段的地址. 

```go
type Mutex struct {
	state int32
	sema  uint32
}
```

所以对于阻塞在同一个`sync.Mutex`上的 `goroutine`, 他们的 `addr` 是相同的, 所以他们会被加入到同一个`sudog` 链表中. 如果是不同的 `sync.Mutex` 锁, 他们的 `addr` 是不同的, 那么他们会被加入到这个 `Treap` 不同的节点. 

进而, 你可以知道, 这个 `rootSema` 是维护多个 `sync.Mutex` 的等待队列的, 可以快速找到不同的 `sync.Mutex` 的等待队列,也可以维护同一个 `sync.Mutex` 的等待队列. 这给了我们启发, 如果你有类似的需求, 可以参考这个实现. 

③ 就是设置这个节点的优先级, 它是一个随机值, 用于保持 `Treap` 的平衡. 这里有个技巧就是总是把优先级最低位设置为 1, 这样保证优先级不为 0.因为优先级经常和 0 做比较, 我们将最低位设置为 1, 就表明优先级已经设置. 

④ 就是将这个新加入的节点旋转到合适的位置, 以保持 `Treap` 的平衡. 这里的旋转操作就是上面提到的左旋和右旋. 稍后看. 

## 移除一个元素 (出队) 
-----------------

对应的, 还有出对的操作. 这个操作就是从 `Treap` 中移除一个节点, 这个节点就是一个等待的 goroutine(`sudog`). 

`dequeue` 搜索并找到在`semaRoot`中第一个因`addr`而阻塞的`goroutine`. 比如需要唤醒一个 goroutine, 让它继续执行(比如直接将锁交给它, 或者唤醒它去争抢锁). 

```go
func (root *semaRoot) dequeue(addr *uint32) (found *sudog, now, tailtime int64) {
	ps := &root.treap
	s := *ps
	for ; s != nil; s = *ps { // ①， 二叉搜索树查找
		if s.elem == unsafe.Pointer(addr) { // ②
			goto Found
		}
		if uintptr(unsafe.Pointer(addr)) < uintptr(s.elem) {
			ps = &s.prev
		} else {
			ps = &s.next
		}
	}
	return nil, 0, 0

Found: // ③
	now = int64(0)
	if s.acquiretime != 0 {
		now = cputicks()
	}
	if t := s.waitlink; t != nil { // ④
		// Substitute t, also waiting on addr, for s in root tree of unique addrs.
		*ps = t
		t.ticket = s.ticket
      ... // 赋值
	} else { // ⑤
		// 旋转s到叶节点，以便删除
		for s.next != nil || s.prev != nil {
			if s.next == nil || s.prev != nil && s.prev.ticket < s.next.ticket {
				root.rotateRight(s)
			} else {
				root.rotateLeft(s)
			}
		}
		// Remove s, now a leaf.
		if s.parent != nil {
			if s.parent.prev == s {
				s.parent.prev = nil
			} else {
				s.parent.next = nil
			}
		} else {
			root.treap = nil
		}
		tailtime = s.acquiretime
	}
	... // 清理s的不需要的信息
	return s, now, tailtime
}
```

① 是遍历 `Treap` 的过程, 当然它是通过搜索二叉树的方式实现. `addr` 就是我们一开始讲的 `Treap` 的 `key`, 也就是 `s.elem`. 如果找到了, 就跳到 `Found` 标签. 如果没有找到, 就返回 `nil`. 

④ 是检查这个地址上是不是有多个等待的 `goroutine`, 如果有, 就把这个节点替换成链表中的下一个节点. 把这个节点从 `Treap` 中移除并返回. 如果就一个 `goroutine`, 那么把这个移除掉后, 需要旋转 `Treap`, 直到这个节点被旋转到叶节点, 然后删除这个节点. 

这里的旋转操作就是上面提到的左旋和右旋. 

## 左旋 rotateLeft
-----------------

`rotateLeft` 函数将以 `x` 为根的子树左旋, 使其变为 `y` 为根的子树. 左旋之前的结构为 `(x a (y b c))`, 旋转后变为 `(y (x a b) c)`. 

```go
// rotateLeft rotates the tree rooted at node x.
// turning (x a (y b c)) into (y (x a b) c).
func (root *semaRoot) rotateLeft(x *sudog) {
	// p -> (x a (y b c))
	p := x.parent
	y := x.next
	b := y.prev

	y.prev = x // ①
	x.parent = y // ②
	x.next = b // ③
	if b != nil {
		b.parent = x // ④
	}

	y.parent = p // ⑤
	if p == nil {
		root.treap = y // ⑥
	} else if p.prev == x { // ⑦
		p.prev = y
	} else {
		if p.next != x {
			throw("semaRoot rotateLeft")
		}
		p.next = y
	}
}
```

具体步骤：

* 将 `y` 设为 `x` 的父节点(②), `x` 设为 `y` 的左子节点(①). 
    
* 将 `b` 设为 `x` 的右子节点(③), 并更新其父节点为 `x`(④). 
    
* 更新 `y` 的父节点为 `p`(⑤), 即 `x` 的原父节点. 如果 `p` 为 nil, 则 `y` 成为新的树根(⑥). 
    
* 根据 `y` 是 `p` 的左子节点还是右子节点, 更新对应的指针(⑦). 
    
左旋为图示具体如下:

```bash
    x            y
   / \          / \
  a   y   =>   x   c
     / \      / \
     b  c    a   b
```

## 右旋 rotateRight
-----------------

rotateRight 旋转以节点 y 为根的树. 将 `(y (x a b) c)` 变为 `(x a (y b c))`. 

```go
func (root *semaRoot) rotateRight(y *sudog) {
	// p -> (y (x a b) c)
	p := y.parent
	x := y.prev
	b := x.next

	x.next = y // ①
	y.parent = x // ②
	y.prev = b // ③
	if b != nil {
		b.parent = y // ④
	}

	x.parent = p // ⑤
	if p == nil {
		root.treap = x // ⑥
	} else if p.prev == y { // ⑦
		p.prev = x
	} else {
		if p.next != y {
			throw("semaRoot rotateRight")
		}
		p.next = x
	}
}
```

具体步骤：

* 将 `y` 设为 `x` 的右子节点(①), `x` 设为 `y` 的父节点(②)
    
* 将 `b` 设为 `y` 的左子节点(③), 并更新其父节点为 `y`(④)
    
* 更新 `x` 的父节点为 `p`(⑤), 即 `y` 的原父节点. 如果 `p` 为 nil, 则 `x` 成为新的树根(⑥)
    
* 根据 `x` 是 `p` 的左子节点还是右子节点, 更新对应的指针(⑦)
    
右旋为图示具体如下:

```bash
      y        x
     / \      / \
    x   c => a   y
   / \          / \
  a   b        b   c
```

理解了左旋和右旋, 你就理解了出队代码中这一段为什么把当前节点旋转到叶结点中了：

```go
		// 旋转s到叶节点，以便删除
		for s.next != nil || s.prev != nil {
			if s.next == nil || s.prev != nil && s.prev.ticket < s.next.ticket {
				root.rotateRight(s)
			} else {
				root.rotateLeft(s)
			}
		}
```

整体上看, `Treap` 这个数据结构确实简单可维护. 左旋和右旋的代码量很少, 结合图看起来也容易理解. 出入队的代码也很简单, 只是简单的二叉搜索树的操作, 加上旋转操作. 
