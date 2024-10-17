package structure

import (
	"sync/atomic"
	"unsafe"
)

// PoolDequeue 是一个无锁的固定大小单生产者、多消费者队列。
// 单个生产者可以从头部推入和弹出，消费者可以从尾部弹出。
// 它具有一个额外的功能，即将未使用的槽位置为 nil，以避免不必要的对象保留。
// 这对于 sync.Pool 很重要，但通常在文献中不被考虑。

type PoolDequeue struct {
	// headTail 将 32 位的头索引和 32 位的尾索引打包在一起。
	// 两者都是 vals 的索引，模 len(vals)-1。
	//
	// tail = 队列中最旧数据的索引
	// head = 下一个要填充的槽位的索引
	//
	// [tail, head) 范围内的槽位由消费者拥有。
	// 消费者在将槽位置为 nil 之前继续拥有该槽位，此时所有权转移给生产者。
	//
	// 头索引存储在最高有效位中，以便我们可以原子地增加它，并且溢出是无害的。
	headTail atomic.Uint64

	// vals 是存储在此队列中的 interface{} 值的环形缓冲区。
	// 其大小必须是 2 的幂。
	//
	// vals[i].typ 为 nil 表示槽位为空，否则非空。
	// 槽位在尾索引移动到其之外并且 typ 被设置为 nil 之前仍在使用。
	// 这由消费者原子地设置为 nil，并由生产者原子地读取。
	vals []eface
}

type eface struct {
	typ, val unsafe.Pointer
}

const dequeueBits = 32

// dequeueLimit 是 PoolDequeue 的最大大小。
// 这必须最多为 (1<<dequeueBits)/2，因为检测满时依赖于环形缓冲区的环绕而不环绕索引。
// 我们除以 4 以便在 32 位上适合 int。
const dequeueLimit = (1 << dequeueBits) / 4

// dequeueNil 用于在 PoolDequeue 中表示 interface{}(nil)。
// 由于我们使用 nil 表示空槽位，因此需要一个哨兵值来表示 nil。
type dequeueNil *struct{}

// NewPoolDequeue 返回一个可以容纳 n 个元素的新 PoolDequeue。
func NewPoolDequeue(n int) *PoolDequeue {
	d := &PoolDequeue{
		vals: make([]eface, n),
	}
	return d
}

func (d *PoolDequeue) unpack(ptrs uint64) (head, tail uint32) {
	const mask = 1<<dequeueBits - 1
	head = uint32((ptrs >> dequeueBits) & mask)
	tail = uint32(ptrs & mask)
	return
}

func (d *PoolDequeue) pack(head, tail uint32) uint64 {
	const mask = 1<<dequeueBits - 1
	return (uint64(head) << dequeueBits) |
		uint64(tail&mask)
}

// PushHead 在队列的头部添加 val。如果队列已满，则返回 false。
// 只能由单个生产者调用。
func (d *PoolDequeue) PushHead(val any) bool {
	ptrs := d.headTail.Load()
	head, tail := d.unpack(ptrs)
	if (tail+uint32(len(d.vals)))&(1<<dequeueBits-1) == head {
		// 队列已满。
		return false
	}
	slot := &d.vals[head&uint32(len(d.vals)-1)]

	// 检查头槽位是否已被 popTail 释放。
	typ := atomic.LoadPointer(&slot.typ)
	if typ != nil {
		// 另一个 goroutine 仍在清理尾部，因此队列实际上仍然满。
		return false
	}

	// 头槽位是空闲的，因此我们拥有它。
	if val == nil {
		val = dequeueNil(nil)
	}
	*(*any)(unsafe.Pointer(slot)) = val

	// 增加头部。这将槽位的所有权传递给 popTail，并作为写入槽位的存储屏障。
	d.headTail.Add(1 << dequeueBits)
	return true
}

// PopHead 移除并返回队列头部的元素。如果队列为空，则返回 false。
// 只能由单个生产者调用。
func (d *PoolDequeue) PopHead() (any, bool) {
	var slot *eface
	for {
		ptrs := d.headTail.Load()
		head, tail := d.unpack(ptrs)
		if tail == head {
			// 队列为空。
			return nil, false
		}

		// 确认尾部并减少头部。我们在读取值之前执行此操作以收回此槽位的所有权。
		head--
		ptrs2 := d.pack(head, tail)
		if d.headTail.CompareAndSwap(ptrs, ptrs2) {
			// 我们成功收回了槽位。
			slot = &d.vals[head&uint32(len(d.vals)-1)]
			break
		}
	}

	val := *(*any)(unsafe.Pointer(slot))
	if val == dequeueNil(nil) {
		val = nil
	}
	// 清零槽位。与 popTail 不同，这里没有与 pushHead 竞争，因此我们不需要小心。
	*slot = eface{}
	return val, true
}

// PopTail 移除并返回队列尾部的元素。如果队列为空，则返回 false。
// 可以由任意数量的消费者调用。
func (d *PoolDequeue) PopTail() (any, bool) {
	var slot *eface
	for {
		ptrs := d.headTail.Load()
		head, tail := d.unpack(ptrs)
		if tail == head {
			// 队列为空。
			return nil, false
		}

		// 确认头部和尾部（用于我们的推测性检查）并增加尾部。如果成功，则我们拥有尾部的槽位。
		ptrs2 := d.pack(head, tail+1)
		if d.headTail.CompareAndSwap(ptrs, ptrs2) {
			// 成功。
			slot = &d.vals[tail&uint32(len(d.vals)-1)]
			break
		}
	}

	// 我们现在拥有槽位。
	val := *(*any)(unsafe.Pointer(slot))
	if val == dequeueNil(nil) {
		val = nil
	}

	// 告诉 pushHead 我们已完成此槽位。清零槽位也很重要，以免留下引用，可能会使对象保持活动状态的时间比必要的长。
	// 我们先写入 val，然后通过原子写入 typ 来发布我们已完成此槽位。
	slot.val = nil
	atomic.StorePointer(&slot.typ, nil)
	// 此时 pushHead 拥有槽位。

	return val, true
}

// PoolChain 是 PoolDequeue 的动态大小版本。
// 这是通过一个 PoolDequeues 的双向链表队列实现的，其中每个队列的大小是前一个的两倍。
// 一旦一个队列填满，就会分配一个新的队列，并且只会推入到最新的队列。
// 弹出操作从链表的另一端进行，一旦一个队列耗尽，它就会从链表中移除。

type PoolChain struct {
	// head 是要推入的 PoolDequeue。这仅由生产者访问，因此不需要同步。
	head *poolChainElt

	// tail 是要从中 popTail 的 PoolDequeue。这由消费者访问，因此读写必须是原子的。
	tail atomic.Pointer[poolChainElt]
}

// NewPoolChain 返回一个新的 PoolChain。
func NewPoolChain() *PoolChain {
	return &PoolChain{}
}

type poolChainElt struct {
	PoolDequeue

	// next 和 prev 链接到此 PoolChain 中的相邻 PoolChainElts。
	//
	// next 由生产者原子地写入，并由消费者原子地读取。它仅从 nil 过渡到非 nil。
	//
	// prev 由消费者原子地写入，并由生产者原子地读取。它仅从非 nil 过渡到 nil。
	next, prev atomic.Pointer[poolChainElt]
}

// PushHead 在队列的头部添加 val。如果队列已满，则返回 false。
// 只能由单个生产者调用。
func (c *PoolChain) PushHead(val any) bool {
	d := c.head
	if d == nil {
		// 初始化链。
		const initSize = 8 // 必须是 2 的幂
		d = new(poolChainElt)
		d.vals = make([]eface, initSize)
		c.head = d
		c.tail.Store(d)
	}

	if d.PushHead(val) {
		return true
	}

	// 当前队列已满。分配一个大小为两倍的新队列。
	newSize := len(d.vals) * 2
	if newSize >= dequeueLimit {
		// 无法再做大。
		newSize = dequeueLimit
	}

	d2 := &poolChainElt{}
	d2.prev.Store(d)
	d2.vals = make([]eface, newSize)
	c.head = d2
	d.next.Store(d2)
	return d2.PushHead(val)
}

// popHead 移除并返回队列头部的元素。如果队列为空，则返回 false。
// 只能由单个生产者调用。
func (c *PoolChain) PopHead() (any, bool) {
	d := c.head
	for d != nil {
		if val, ok := d.PopHead(); ok {
			return val, ok
		}
		// 之前的队列中可能仍有未消费的元素，因此尝试回退。
		d = d.prev.Load()
	}
	return nil, false
}

// PopTail 移除并返回队列尾部的元素。如果队列为空，则返回 false。
// 可以由任意数量的消费者调用。
func (c *PoolChain) PopTail() (any, bool) {
	d := c.tail.Load()
	if d == nil {
		return nil, false
	}

	for {
		// 重要的是，我们在弹出尾部之前加载 next 指针。
		// 一般来说，d 可能是暂时空的，但如果 next 在弹出之前是非 nil 并且弹出失败，则 d 是永久空的，
		// 这是唯一可以安全地从链中删除 d 的条件。
		d2 := d.next.Load()

		if val, ok := d.PopTail(); ok {
			return val, ok
		}

		if d2 == nil {
			// 这是唯一的队列。它现在是空的，但将来可能会被推入。
			return nil, false
		}

		// 链的尾部已被耗尽，因此继续下一个队列。
		// 尝试将其从链中删除，以便下一个弹出不必再次查看空队列。
		if c.tail.CompareAndSwap(d, d2) {
			// 我们赢得了竞争。清除 prev 指针，以便垃圾收集器可以收集空队列，并且 popHead 不会回退得更远。
			d2.prev.Store(nil)
		}
		d = d2
	}
}
