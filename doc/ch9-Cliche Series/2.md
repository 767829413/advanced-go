# 老生常谈2: Go pprof use

## pprof使用介绍

检查Go程序内存的使用情况最常用的就是Go标准库自带的pprof库了，可以通过http暴露出这个profile, 然后通过go tool pprof或者pprof工具命令行/web方式查看。下面主要介绍下如何使用

### 采集数据

1. 通过暴露端口路由

```go
import (
    "fmt"
    "net/http"
    _ "net/http/pprof"  // 导入 pprof 包

   	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func main() {
    // 注册 pprof 的 HTTP 路由
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()

    // 其他服务器代码
    http.HandleFunc("/", handler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}

// 如果使用gin,那么也可以这样
// 添加性能监控
func RegisterHttpForProfiling(path string) {
	/*初始化http服务*/
	port := 82
	address := fmt.Sprintf("0.0.0.0:%d", port)
	router := gin.New()
	pprof.Register(router, path)
	if err := router.Run(address); err != nil {
		fmt.Printf("net pprof is start failed: %s", err.Error())
	}
}

func RegisterRouterForProfiling(router *gin.Engine) {
	// 添加基本的校验机制
	debugGroup := router.Group("/mypprof/:pass", func(c *gin.Context) {
		if c.Param("pass") != "212123453" {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Next()
	})
	pprof.RouteRegister(debugGroup, "debug/pprof")
}
```

2. 输出文件

```go
package monitor

import (
	"fmt"
	"net/http"
	"os"
	"runtime/pprof"
	"runtime/trace"

	"github.com/gin-gonic/gin"
)

const (
	MonitorPath = "/monitor/notify"
)

var (
	notify = gin.H{
		"code": 0,
		"msg":  "notify success",
	}
	open = make(chan struct{})
)

func NotifyHandler(c *gin.Context) {
	open <- struct{}{}
	c.JSON(http.StatusOK, notify)
}

func RegisterNotifyForProfiling(serverName string, domain string, dir string) {
	go func() {
		var memoryProfile, cpuProfile, traceProfile *os.File
		started := false
		for range open {
			if started {
				pprof.StopCPUProfile()
				trace.Stop()
				pprof.WriteHeapProfile(memoryProfile)
				memoryProfile.Close()
				cpuProfile.Close()
				traceProfile.Close()
				started = false
			} else {
				fileNamePre := fmt.Sprintf("%s/%s_%s", dir, serverName, domain)
				cpuProfile, _ = os.Create(fmt.Sprintf("%s.%s", fileNamePre, "cpu.pprof"))
				memoryProfile, _ = os.Create(fmt.Sprintf("%s.%s", fileNamePre, "memory.pprof"))
				traceProfile, _ = os.Create(fmt.Sprintf("%s.%s", fileNamePre, "runtime.trace"))
				pprof.StartCPUProfile(cpuProfile)
				trace.Start(traceProfile)
				started = true
			}
		}
	}()
}
```

当然,更建议使用这个库: <https://github.com/mosn/holmes>

### 分析数据

上述输出的分析文件可以通过 `go tool pprof` 命令进行分析

```bash
# 直接分析采样地址
#获取内存采样
go tool pprof -http :8848 "http://host/debug/pprof/heap?debug=1"
#获取goroutine采样
go tool pprof -http :8849 "http://host/debug/pprof/goroutine?debug=1"

# 直接分析采样文件
go tool pprof -http :8849 "your_path/pprof.goroutine.001.pb.gz"
```

## 举个例子

```bash
go tool pprof -http :8848 "http://host/debug/pprof/heap?debug=1"
```

可以获取服务器 `http://host` 的堆信息，并且在本机 8848 端口启动一个服务器展示堆的信息。

![img](https://pic.imgdb.cn/item/66485d24d9c307b7e97f7b85.jpg)

正常使用`pprof`可以直观发现哪些代码分配内存多，有助于识别潜在内存泄漏。但高内存分配不必然意味泄漏，它可能指这部分内存曾经或正在频繁分配，且可能随后被回收。

类似Java的profiler，pprof能通过比较不同时间点内存分配的差异，揭示哪些区域的内存未被释放，这可能是内存泄漏的迹象。

操作步骤如下：

1. 确保`pprof`的HTTP路径已配置，访问http://ip:port/debug/pprof/以确认。
2. 导出初始时间点的堆`profile`作为基准：`curl -s http://host/debug/pprof/heap > base.heap`. 
3. 经过一段时间后，导出第二时间点的堆`profile`：`curl -s http://host/debug/pprof/heap > current.heap`. 
4. 使用`pprof`比较两个时间点的堆差异：`go tool pprof --base base.heap current.heap`. 

操作和正常的go tool pprof操作一样， 比如使用top查看使用堆内存最多的几处地方的内存增删情况：

![img](https://pic.imgdb.cn/item/6648605ad9c307b7e98427c9.jpg)

或者你直接使用命令打开`web`界面: `go tool pprof --http :9090 --base base.heap current.heap`. 

![img](https://pic.imgdb.cn/item/66486108d9c307b7e984f349.jpg)

## 后续思考和问题

### 问题: 打开的 heap profile 不能正确反应当前分配的内存吗

#### 问题代码

```go
package main

import (
 "fmt"
 "net/http"
 _ "net/http/pprof"
 "runtime"
 "time"
)

func main() {
 go func() {
  http.ListenAndServe("localhost:8080", nil)
 }()

 go func() {
  // 每秒打印内存分配情况
  for {
   var m runtime.MemStats
   runtime.ReadMemStats(&m)
   fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
   fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
   fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
   fmt.Printf("\tNumGC = %v\n", m.NumGC)
   time.Sleep(1 * time.Second)
  }

 }()
 time.Sleep(5 * time.Second)

 fmt.Println("start test")

 // 创建一个 200 MiB 的切片
 var memoryLeaks [][]int32
 for i := 0; i < 10; i++ {
  leak := make([]int32, 5*1024*1024) // 分配 5*1M*4bytes = 20 MiB
  memoryLeaks = append(memoryLeaks, leak)
  time.Sleep(1 * time.Second) // 延迟一秒观察内存分配情况
 }
 // 期望至少分配了 200 MiB 内存
 fmt.Println("end test")
 // 看到上面的文字后，打开go pprof 工具，查看工具的分析
 // go tool pprof -http :8972 http://127.0.0.1:8080/debug/pprof/heap


 time.Sleep(1 * time.Hour)
 fmt.Println("test", memoryLeaks[9][5*1024*1024-1]) // 避免垃圾回收和优化

}

func bToMb(b uint64) uint64 {
 return b / 1024 / 1024
}
```

上面这个程序首先创建一个监听 `8080` 端口的 `web` 服务，主要利用它访问 `http heap profile`，所以导入了`_ "net/http/pprof"`包。

然后创建了了一个切片，切片包含 `10` 个元素，每个元素又是一个 `[]int32` 的切片。每个元素占用 `20 MiB` 的大小 (`int32` 是 `4` 个字节，`5*1024*1024` 是 `5M`, 所以是 `4*5 MiB=20 Mib`). 

所以期望 `10` 个循环创建完毕后，至少此程序会占用 `200 MiB` 内存，另外在加一些栈、网络、运行时等一些额外的内存的话，略微会比 `200MiB` 大一些。

执行这个程序后，可以看到程序每秒的内存占用统计：

![img](https://pic.imgdb.cn/item/66486573d9c307b7e98a2f4e.jpg)

可以看到使用 `runtime.MemStats` 统计的内存分配情况还是比较准的， `200MiB`，符合预期。

当看到 `end test` 时，表明 `200MiB` 已经分配，可以运行下面的命令打开 `heap`

```go
go tool pprof -http :8972 "http://127.0.0.1:8080/debug/pprof/heap"
```

![img](https://pic.imgdb.cn/item/664862f1d9c307b7e9873913.jpg)

***这个时候 `heap profile` 才显示 `80 MB` 的内存!!少了 `120 MB`。和期望的严重不符!!!***

这 `80 MB` 内存的确是在创建元素的时候分配的：

![img](https://pic.imgdb.cn/item/664868d1d9c307b7e98e2a5d.jpg)

可那 `120 MB` 呢？如果 `Heap Profile` 数据不准，那官方团队早就出来洗地了!

应该不是 `Go` 实现的问题，那么就要想一下是不是采样的机制导致的,下面开始针对这个问题进行分析。

#### 缺少的分配内存

首先再等一会，大约在命令行中看到 `end test` 再等待两到五分钟吧，当然等的时间越长越好，在等待时间结束后,开始重新执行命令，查看 `heap profile`:

```bash
go tool pprof -http :8972 "http://127.0.0.1:8080/debug/pprof/heap"
```

![img](https://pic.imgdb.cn/item/66486b1ad9c307b7e991328a.jpg)

这个时候浏览器中显示出了 `200 MB` 的内存分配，符合预期：

等待一段时间之后，`heap profile` 就显示正常了。难道过了一会才分配的内存？肯定不可能。因为命令行中 `runtime.MemStats` 已经显示结束测试前内存已经分配好了。

看到 `end test` 之后，程序只是休眠了 1 小时。有没有可能是垃圾回收导致的呢?毕竟要经过一段时间就要执行垃圾回收的.

可以直接验证一下，在 `end test` 之后强制垃圾回收一下，再立即打开 `heap profile` 是不是显示 `200 MiB`. 

```go
	// 期望至少分配了 200 MiB 内存
	fmt.Println("end test")
	// 这里强制GC一下
	runtime.GC()
```

![img](https://pic.imgdb.cn/item/66486d22d9c307b7e99354df.jpg)

果然是垃圾回收的原因，如果直接在代码里强制垃圾回收，来立即显示当前已分配的内存了，确实不用再等待。

但是线上运行的程序中也不可能随心所欲的让我们在任意的地方动态加 `runtime.GC()` 吧。

所以，访问 `heap profile` 需要加上 `gc=1` 的参数即可 (大于 `0` 的数都可以)：

```bash
go tool pprof -http :8972 "http://127.0.0.1:8080/debug/pprof/heap?gc=1"
```

在程序中把 runtime.GC() 那一行去掉，使用上面的方式访问 `heap profile`.

```text
还可以加上 debug=1 (非 0 的参数)，可以通过文字的方式查看内存分配以及 runtime.MemStats 数据。如 http://127.0.0.1:8080/debug/pprof/heap?gc=1&debug=1
```

#### 原因探讨

`Heap profile` 其实调用的 `runtime.MemProfile` 进行统计:

```go
func writeHeapInternal(w io.Writer, debug int, defaultSampleType string) error {
 ......
	var p []runtime.MemProfileRecord
	n, ok := runtime.MemProfile(nil, true)
	for {
		p = make([]runtime.MemProfileRecord, n+50)
		n, ok = runtime.MemProfile(p, true)
		if ok {
			p = p[0:n]
			break
		}
	}
    ....
}
```

`runtime.MemProfile` 的方法签名如下，它负责统计内存的分配情况：

![img](https://pic.imgdb.cn/item/6648766bd9c307b7e99d871f.png)

核心的一句翻译

```text
返回的性能分析数据可能最多延迟两个垃圾回收周期。这样做是为了避免对内存分配产生偏差;因为内存分配是实时发生的,但释放操作直到垃圾回收器进行清理时才会延迟进行,所以性能分析只会统计那些已经有机会被垃圾回收器回收的内存分配
```

一开始的问题根因就在这里, 看到 `end test` 立即查看 `heap profile`, 还没有进行垃圾回收，相关的内存统计数据还没有计算出来，所以才看到 `80` MB 的内存，而不是 `200MB` 的内存。

`system.GC`会强制进行垃圾回收，并且会发布新的 `heap profile`，所以可以看到 `200MB` 的内存。

```go
func GC() {
...
	// Now we're really done with sweeping, so we can publish the
	// stable heap profile. Only do this if we haven't already hit
	// another mark termination.
	mp := acquirem()
	cycle := work.cycles.Load()
	if cycle == n+1 || (gcphase == _GCmark && cycle == n+2) {
		mProf_PostSweep()
	}
	releasem(mp)
}
```

#### Heap Profile 到底怎么采样的?

`heap profile` 实际上在内存分配时做采样统计的，默认情况下并不会记录所有的内存分配。

这里需要关注 `runtime.MemProfileRate` 这个变量。

```go
// MemProfileRate controls the fraction of memory allocations
// that are recorded and reported in the memory profile.
// The profiler aims to sample an average of
// one allocation per MemProfileRate bytes allocated.
//
// To include every allocated block in the profile, set MemProfileRate to 1.
// To turn off profiling entirely, set MemProfileRate to 0.
//
// The tools that process the memory profiles assume that the
// profile rate is constant across the lifetime of the program
// and equal to the current value. Programs that change the
// memory profiling rate should do so just once, as early as
// possible in the execution of the program (for example,
// at the beginning of main).
var MemProfileRate int = 512 * 1024
```

翻译一下表示: 

```text
MemProfileRate 控制了在内存分析中记录和报告的内存分配的比例。分析器的目标是每分配 MemProfileRate 字节时,采样记录一次分配。

如果要在分析中包含每一个分配的内存块,可以将 MemProfileRate 设置为 1。如果要完全关闭内存分析,可以将 MemProfileRate 设置为 0。

处理内存分析数据的工具假设分析的采样率在程序的整个生命周期中保持不变,并等于当前设置的值。如果程序需要改变内存分析的采样率,应该只改变一次,并且尽可能早地在程序执行的开始阶段(例如在 main 函数的开始部分)进行设置。
```

如果把每次分配的大小改小一点，如下：

```go
func main() {
...
	for i := 0; i < 10; i++ {
		// leak := make([]int32, 5*1024*1024) // 分配 5*1M*4bytes = 20 MiB
		leak := make([]int32, 60*1024) // 分配 5*1M*4bytes = 20 MiB
		memoryLeaks = append(memoryLeaks, leak)
		time.Sleep(1 * time.Second) // 延迟一秒观察内存分配情况
	}
...
}
```

这个时候查看 `heap profile`，会出现 `heap profile` 和 `runtime.MemStats` 统计的内存分配情况不一样，因为 `runtime.MemStats` 是实时统计的，而 `heap profile` 是采样统计的。

如果在程序一开始加上 `runtime.MemProfileRate = 1`, `heap profle` 就会保持一致了。 

但是也需要注意，将 `runtime.MemProfileRate` 设置为 1，对每一次内存都进行采样统计，对性能的影响也相应的增大了，非必要不用更改这个值。

这样就表示默认看到 `heap profile`, 都是在默认值 `512 * 1024` 这个阈值下采样统计的，和实际的 `heap` 上的内存分配可能会有些出入，但是对于分析内存泄露，影响不大。

具体操作这里就不赘叙了.
