# Golang 1.23 新包：unique

> `unique` 包提供了一些工具，用来对"可比较的值""进行规范化处理（即“驻留”）. 
> 
> 具体来说，“规范化”（canonicalization）或“驻留”（interning）指的是将多个相同的值（例如相同内容的字符串或结构体）通过某种机制合并成一个唯一的副本. 这样，当有多个相同的值时，它们在内存中只会保存一个规范化的版本，其他相同的值都指向这个唯一的副本，从而节省内存并加速相等性比较操作. 

在 Go 官方博客上，unique 包的主刀 Michael Knyszek 写了一篇关于 unique 包的介绍<https://go.dev/blog/unique>，并说明了实现这个包过程中一些新发现 (弱指针、finalizer 替代者)，以下是翻译：

Go 1.23 标准库中引入了一个新包，名为 `unique`. 这个包的目标是实现可比值的规范化，简单来说，它允许你对值进行去重，从而让它们指向唯一的、规范化的副本，并在底层高效管理这些规范化副本. 你可能已经对这种概念有所了解，称之为“驻留（interning）”. 让我们深入了解一下它是如何工作的，以及它为什么有用. 

## 一个简单的驻留实现
---------

从宏观上看，驻留非常简单. 以下代码示例展示了如何通过普通的 Map 来对字符串进行去重. 

```go
var internPool map[string]string

// Intern 返回一个与 s 相等的字符串，但可能与之前传给 Intern 的字符串共享存储. 
func Intern(s string) string {
    pooled, ok := internPool[s]
    if !ok {
        // 克隆字符串以防它是某个更大字符串的一部分. 
        // 如果驻留使用得当，这种情况应该很少发生. 
        pooled = strings.Clone(s)
        internPool[pooled] = pooled
    }
    return pooled
}
```

当你构建许多可能重复的字符串时，例如解析文本格式时，这非常有用.   
然而，这种实现虽然简单，但存在一些问题：

1. 它永远不会从对象池中移除字符串. 
    
2. 它无法安全地在多个 goroutine 中并发使用. 
    
3. 它仅适用于字符串，而这个想法其实是普遍适用的. 
    
此外，这个实现还错过了一个微妙的优化机会. 字符串在底层是不可变的结构，包含一个指针和一个长度. 当比较两个字符串时，如果指针不相等，就必须比较它们的内容以确定是否相等. 但如果我们知道两个字符串是规范化的，那么只需比较它们的指针即可. 

## 引入 `unique` 包
---------

新引入的 `unique` 包提供了一个类似于 `Intern` 的函数 `Make`，它的工作方式与 `Intern` 类似. 在内部，它也有一个全局 Map（一个快速的 [泛型并发 Map](https://pkg.go.dev/internal/concurrent@go1.23.0)），并在该 Map 中查找值. 然而，它与 `Intern` 有两个重要的区别：首先，它接受任何可比较类型的值；其次，它返回一个包装值 `Handle[T]`，可以从中检索规范化的值. 

`Handle[T]` 是设计的关键. `Handle[T]` 有这样一个特性：只有当用来创建它的两个值相等时，两个 `Handle[T]` 才相等. 更重要的是，两个 `Handle[T]` 的比较是非常廉价的：只需进行指针比较. 相比之下，比较两个长字符串的成本要高得多！

到目前为止，这些功能都可以通过普通的 Go 代码实现. 然而，`Handle[T]` 还有第二个作用：只要某个值存在一个 `Handle[T]`，Map 就会保留该值的规范化副本. 一旦所有 Map 到特定值的 `Handle[T]` 都消失，该包就会将内部 Map 项标记为可删除，供垃圾回收器在未来回收. 这为何时从 Map 中移除条目设定了明确的策略：当规范化条目不再被使用时，垃圾回收器可以自由清理它们. 

如果你曾经使用过 Lisp，这一切可能听起来很熟悉. Lisp 中的符号是驻留的字符串，但它们本身并不是字符串，所有符号的字符串值都保证位于同一个池中. 这种符号与字符串的关系类似于 `Handle[string]` 与 `string` 的关系. 

## 一个实际例子
---------

如何使用 `unique`？可以看看标准库中的 `net/netip` 包，它对 `netip.Addr` 结构中的 `addrDetail` 类型的值进行了驻留. 以下是 `net/netip` 中实际代码的简化版本，它使用了 `unique` 包. 

```go
type Addr struct {
    // 与地址相关的详细信息，被打包在一起并进行了规范化. 
    z unique.Handle[addrDetail]
}

type addrDetail struct {
    isV6   bool   // 如果是 IPv4，则为 false；如果是 IPv6，则为 true. 
    zoneV6 string // 如果是 IPv6，可能不等于 "". 
}

var z6noz = unique.Make(addrDetail{isV6: true})

// WithZone 返回一个与 ip 相同的 IP，但带有指定的 zone. 如果 zone 为空，则移除 zone. 
func (ip Addr) WithZone(zone string) Addr {
    if !ip.Is6() {
        return ip
    }
    if zone == "" {
        ip.z = z6noz
        return ip
    }
    ip.z = unique.Make(addrDetail{isV6: true, zoneV6: zone})
    return ip
}
```

由于许多 IP 地址可能使用相同的 zone，且该 zone 是它们标识的一部分，因此对它们进行规范化非常合理. Zone 的去重减少了每个 `netip.Addr` 的平均内存占用量，而它们被规范化后，比较 zone 名称只需进行简单的指针比较，这使得值的比较更加高效. 

## 关于字符串驻留的注脚
---------

尽管 `unique` 包很有用，但它与字符串的驻留不太一样，因为为了防止字符串被从内部 Map 中删除使用 `Handle[T]` 是必须的. 这意味着你需要修改代码以同时保留 `Handle[T]` 和字符串. 

但字符串特殊之处在于，虽然它们表现得像值，但实际上它们的底层包含指针. 因此，理论上可以只对字符串的底层存储进行规范化，而将 `Handle[T]` 的细节隐藏在字符串内部. 因此，未来仍然有可能实现所谓的透明字符串驻留，即可以在不需要 `Handle[T]` 的情况下对字符串进行驻留，类似于 `Intern` 函数，但语义更像 `Make`. 

目前，`unique.Make("my string").Value()` 是一种可能的解决方法. 即使没有保留 `Handle[T]`，字符串也会被允许从 `unique` 的内部 Map 中删除，但不会立即删除. 实际上，条目至少会在下一次垃圾回收完成后才被删除，因此这种解决方法在回收之间的时间段内仍然允许一定程度的去重. 

## 一些历史与展望
---------

事实上，`net/netip` 包自引入以来就已经对 zone 字符串进行了驻留. 它使用的驻留包是 [go4.org/intern](https://pkg.go.dev/go4.org/intern) 的内部副本. 与 `unique` 包类似，它有一个 `Value` 类型（在泛型之前看起来很像 `Handle[T]`），其内部 Map 中的条目会在不再被引用后被移除. 

为了实现这种行为，旧的 `intern` 包做了一些不安全的事情，特别是在运行时之外实现了弱指针. 而弱指针是 `unique` 包的核心抽象. 弱指针是一种不会阻止垃圾回收器回收变量的指针；当变量被回收时，弱指针会自动变成 `nil`. 

在实现 `unique` 包时，我们为垃圾回收器添加了适当的弱指针支持. 经过设计决策的考验后，我们惊讶地发现这一切竟然如此简单且直接. 弱指针现在已经成为一个[公开提案](https://go.dev/issue/67552). 

这项工作还促使我们重新审视终结器，最终提出了一个更易于使用且效率更高的[终结器替代方案](https://go.dev/issue/67535). 随着可比较值的哈希函数即将推出，Go 中构建内存[高效缓存](https://go.dev/issue/67552#issuecomment-2200755798)的未来充满希望！

## 参考资料
---------

1. unique 包的介绍: https://go.dev/blog/unique

2. 泛型并发 Map: https://pkg.go.dev/internal/concurrent@go1.23.0

3. go4.org/intern: https://pkg.go.dev/go4.org/intern

4. 公开提案: https://go.dev/issue/67552

5. 终结器替代方案: https://go.dev/issue/67535

6. 高效缓存: https://go.dev/issue/67552#issuecomment-2200755798