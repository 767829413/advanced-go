# Golang语法背后的故事

## 编译原理基础

`业务场景`

1. 代码效率

    ![Go相关分析.jpg](https://s2.loli.net/2023/05/26/1CwABDeU6oqhJxT.png)

2. 类型转换

```go
package main

func main() {
	var a = "sdf"
	var b = []byte(a)
	println(b)
}
```

3. 查找某些方法的实现

```go
// The make built-in function allocates and initializes an object of type
// slice, map, or chan (only). Like new, the first argument is a type, not a
// value. Unlike new, make's return type is the same as the type of its
// argument, not a pointer to it. The specification of the result depends on
// the type:
//
//	Slice: The size specifies the length. The capacity of the slice is
//	equal to its length. A second integer argument may be provided to
//	specify a different capacity; it must be no smaller than the
//	length. For example, make([]int, 0, 10) allocates an underlying array
//	of size 10 and returns a slice of length 0 and capacity 10 that is
//	backed by this underlying array.
//	Map: An empty map is allocated with enough space to hold the
//	specified number of elements. The size may be omitted, in which case
//	a small starting size is allocated.
//	Channel: The channel's buffer is initialized with the specified
//	buffer capacity. If zero, or the size is omitted, the channel is
//	unbuffered.
func make(t Type, size ...IntegerType) Type

// The new built-in function allocates memory. The first argument is a type,
// not a value, and the value returned is a pointer to a newly
// allocated zero value of that type.
func new(Type) *Type
```

3. 业务需求大体一致细节多样

```text
会员服务，给⽤户分等级：

初级会员，发贴数 > 10
中级会员，充值 > 1000 RMB
⾼级会员，发帖数 > 100，充值 > 10000 RMB

如果项⽬数 = ⼏百，每个项⽬都有⾃⼰的会员规则，怎么办
```

4. 封装统一的数据查询服务

```text
⽤户提供查询条件，会经常变
代理去不同的模块查数据
外部模块没有统⼀的数据获取规范

每次我们的⽤户提了需求，我们就⼀定要写⼀遍代码么？
```

5. 数据配置更换

```text
公司想从 Thrift 切换到 gRPC

已经有了⼤量的 Thrift IDL
想提供 gRPC 接⼝

⼿⼯把 Thrift IDL 抄写成 pb ⽂件效率太低
```

6. SQL审计

```text
我是 SQL 专家
我知道怎么获取到表的索引
我可以把⽤户代码⾥的 SQL 扫描出来

我想在上线的时候能⾃动做⼀些拦截，提醒⽤户去给表加索引
```

7. 可定制化需求

```text
公司卖出去的服务软件，客户要⾃⼰定制，但不想给客户源码

要⽀持客户写的扩展代码能在我们的模块上运⾏
```

`回顾一下go的编译过程`

```bash
[root@~]# go build -x main.go 
WORK=/tmp/go-build674607137
mkdir -p $WORK/b001/
cat >$WORK/b001/importcfg << 'EOF' # internal
# import config
packagefile runtime=/usr/local/go/pkg/linux_amd64/runtime.a
EOF
cd /home/advanced-go/ch2-Golang syntax
# 编译过程
/usr/local/go/pkg/tool/linux_amd64/compile -o $WORK/b001/_pkg_.a -trimpath "$WORK/b001=>" -p main -complete -buildid KpoOO4emO2I6CPS0W0eq/KpoOO4emO2I6CPS0W0eq -goversion go1.14.12 -D _/home/advanced-go/ch2-Golang_syntax -importcfg $WORK/b001/importcfg -pack -c=4 "/home/advanced-go/ch2-Golang syntax/main.go"
/usr/local/go/pkg/tool/linux_amd64/buildid -w $WORK/b001/_pkg_.a # internal
cp $WORK/b001/_pkg_.a /root/.cache/go-build/1f/1fea6e6271bf218d7b3bd1efe45d8c9a94a23a19cefa2ef7e0dc54ca738e7861-d # internal
cat >$WORK/b001/importcfg.link << 'EOF' # internal
packagefile command-line-arguments=$WORK/b001/_pkg_.a
packagefile runtime=/usr/local/go/pkg/linux_amd64/runtime.a
packagefile internal/bytealg=/usr/local/go/pkg/linux_amd64/internal/bytealg.a
packagefile internal/cpu=/usr/local/go/pkg/linux_amd64/internal/cpu.a
packagefile runtime/internal/atomic=/usr/local/go/pkg/linux_amd64/runtime/internal/atomic.a
packagefile runtime/internal/math=/usr/local/go/pkg/linux_amd64/runtime/internal/math.a
packagefile runtime/internal/sys=/usr/local/go/pkg/linux_amd64/runtime/internal/sys.a
EOF
mkdir -p $WORK/b001/exe/
cd .
# 链接过程
/usr/local/go/pkg/tool/linux_amd64/link -o $WORK/b001/exe/a.out -importcfg $WORK/b001/importcfg.link -buildmode=exe -buildid=-CJAJtpR0IoGh9Ubwiod/KpoOO4emO2I6CPS0W0eq/45AKO0l05Rj8A6cOFtji/-CJAJtpR0IoGh9Ubwiod -extld=gcc $WORK/b001/_pkg_.a
/usr/local/go/pkg/tool/linux_amd64/buildid -w $WORK/b001/exe/a.out # internal
cp $WORK/b001/exe/a.out main
rm -r $WORK/b001/
[root@~]# 
```

`编译流程`

   ![Go相关分析 _1_.jpg](https://s2.loli.net/2023/05/26/MQ7zVXq2eoGSWkH.png)

1. 词法分析

    ![Go相关分析 _2_.jpg](https://s2.loli.net/2023/05/26/GopWdknvHBbMfAe.png)

2. 语法分析

    ![Go相关分析 _3_.jpg](https://s2.loli.net/2023/05/26/uCMOKfV2qIHzGsw.png)

    参考地址: <https://astexplorer.net/>

3. 语义分析

    ![Go相关分析 _4_.jpg](https://s2.loli.net/2023/05/26/LuyemCl4E1tPbDB.png)

    **在抽象语法树 AST 上做类型检查**

4. 中间代码(SSA)⽣成与优化

    ![Go相关分析 _5_.jpg](https://s2.loli.net/2023/05/26/tCKncvUuhYaqp3g.png)

    ```text
    SSA(Single Static Assignment)的两⼤要点是：

        • Static: 每个变量只能赋值⼀次(因此应该叫常量更合适)；
        • Single: 每个表达式只能做⼀个简单运算,对于复杂的表达式a*b+c*d,要拆分成: t0=a*b; t1=c*d; t2=t0+t1; 三个简单表达式；
    ```

    参考地址: <https://golang.design/gossa>

5. 机器码⽣成

    ![ccode.png](https://s2.loli.net/2023/05/26/wRv3yqdV7rKXWf6.png)

    参考地址: <https://godbolt.org/>

`链接过程`

 ![aaa6.png](https://s2.loli.net/2023/05/26/txd5wVKBYRN49gj.png)

**最重要的就是进⾏虚拟地址重定位(Relocation)**

1. 编译后: 所有函数地址都是从 0 开始每条指令是相对函数第⼀条指令的偏移

2. 链接后: 所有指令都有了全局唯⼀的地址

## 编译与反编译⼯具

`go tool compile`

**go tool compile -S ./main.go | grep "main.go:"**

**该命令会生成 main.o 目标文件,并把目标的汇编内容输出**

```go
package main

func main() {
	var str = "hello world"
	var b = []byte(str)
	println(b)
}
```

```bash
[root@~]# go tool compile -S main.go | grep 'main.go:'
        0x0000 00000 (main.go:3)        TEXT    "".main(SB), ABIInternal, $112-0
        0x0000 00000 (main.go:3)        MOVQ    (TLS), CX
        0x0009 00009 (main.go:3)        CMPQ    SP, 16(CX)
        0x000d 00013 (main.go:3)        PCDATA  $0, $-2
        0x000d 00013 (main.go:3)        JLS     157
        0x0013 00019 (main.go:3)        PCDATA  $0, $-1
        0x0013 00019 (main.go:3)        SUBQ    $112, SP
        0x0017 00023 (main.go:3)        MOVQ    BP, 104(SP)
        0x001c 00028 (main.go:3)        LEAQ    104(SP), BP
        0x0021 00033 (main.go:3)        PCDATA  $0, $-2
        0x0021 00033 (main.go:3)        PCDATA  $1, $-2
        0x0021 00033 (main.go:3)        FUNCDATA        $0, gclocals·69c1753bd5f81501d95132d08af04464(SB)
        0x0021 00033 (main.go:3)        FUNCDATA        $1, gclocals·9fb7f0986f647f17cb53dda1484e0f7a(SB)
        0x0021 00033 (main.go:3)        FUNCDATA        $2, gclocals·9fb7f0986f647f17cb53dda1484e0f7a(SB)
        0x0021 00033 (main.go:5)        PCDATA  $0, $1
        0x0021 00033 (main.go:5)        PCDATA  $1, $0
        0x0021 00033 (main.go:5)        LEAQ    ""..autotmp_2+64(SP), AX
        0x0026 00038 (main.go:5)        PCDATA  $0, $0
        0x0026 00038 (main.go:5)        MOVQ    AX, (SP)
        0x002a 00042 (main.go:5)        PCDATA  $0, $1
        0x002a 00042 (main.go:5)        LEAQ    go.string."hello world"(SB), AX
        0x0031 00049 (main.go:5)        PCDATA  $0, $0
        0x0031 00049 (main.go:5)        MOVQ    AX, 8(SP)
        0x0036 00054 (main.go:5)        MOVQ    $11, 16(SP)
        0x003f 00063 (main.go:5)        CALL    runtime.stringtoslicebyte(SB)
        0x0044 00068 (main.go:5)        PCDATA  $0, $1
        0x0044 00068 (main.go:5)        MOVQ    24(SP), AX
        0x0049 00073 (main.go:5)        PCDATA  $0, $0
        0x0049 00073 (main.go:5)        PCDATA  $1, $1
        0x0049 00073 (main.go:5)        MOVQ    AX, "".b.ptr+96(SP)
        0x004e 00078 (main.go:5)        MOVQ    32(SP), CX
        0x0053 00083 (main.go:5)        MOVQ    CX, "".b.len+48(SP)
        0x0058 00088 (main.go:5)        MOVQ    40(SP), DX
        0x005d 00093 (main.go:5)        MOVQ    DX, "".b.cap+56(SP)
        0x0062 00098 (main.go:6)        CALL    runtime.printlock(SB)
        0x0067 00103 (main.go:6)        PCDATA  $0, $1
        0x0067 00103 (main.go:6)        PCDATA  $1, $0
        0x0067 00103 (main.go:6)        MOVQ    "".b.ptr+96(SP), AX
        0x006c 00108 (main.go:6)        PCDATA  $0, $0
        0x006c 00108 (main.go:6)        MOVQ    AX, (SP)
        0x0070 00112 (main.go:6)        MOVQ    "".b.len+48(SP), AX
        0x0075 00117 (main.go:6)        MOVQ    AX, 8(SP)
        0x007a 00122 (main.go:6)        MOVQ    "".b.cap+56(SP), AX
        0x007f 00127 (main.go:6)        MOVQ    AX, 16(SP)
        0x0084 00132 (main.go:6)        CALL    runtime.printslice(SB)
        0x0089 00137 (main.go:6)        CALL    runtime.printnl(SB)
        0x008e 00142 (main.go:6)        CALL    runtime.printunlock(SB)
        0x0093 00147 (main.go:7)        MOVQ    104(SP), BP
        0x0098 00152 (main.go:7)        ADDQ    $112, SP
        0x009c 00156 (main.go:7)        RET
        0x009d 00157 (main.go:7)        NOP
        0x009d 00157 (main.go:3)        PCDATA  $1, $-1
        0x009d 00157 (main.go:3)        PCDATA  $0, $-2
        0x009d 00157 (main.go:3)        CALL    runtime.morestack_noctxt(SB)
        0x00a2 00162 (main.go:3)        PCDATA  $0, $-1
        0x00a2 00162 (main.go:3)        JMP     0
```

`go tool objdump`

**寻找make的底层实现**

官方文档: <https://go.dev/ref/spec#Making_slices_maps_and_channels>

```go
package main

func main() {
	// make slice
	// 为了统一都分配到堆上, 栈上的 slice 结果会有出入
	var sl = make([]int, 100000)
	println(sl)

	// make channel
	var ch = make(chan int, 5)
	println(ch)

	// make map
	var m = make(map[int]int, 22)
	println(m)
}
```

```bash
[root@~]# go build main.go && go tool objdump ./main | grep -E "main.go:6|main.go:10|main.go:14"
  main.go:6             0x458791                488d05c8980000          LEAQ 0x98c8(IP), AX
  main.go:6             0x458798                48890424                MOVQ AX, 0(SP)
  main.go:6             0x45879c                48c7442408a0860100      MOVQ $0x186a0, 0x8(SP)
  main.go:6             0x4587a5                48c7442410a0860100      MOVQ $0x186a0, 0x10(SP)
  main.go:6             0x4587ae                e87d52feff              CALL runtime.makeslice(SB)
  main.go:6             0x4587b3                488b442418              MOVQ 0x18(SP), AX
  main.go:6             0x4587b8                4889442430              MOVQ AX, 0x30(SP)
  main.go:10            0x4587ec                488d056d960000          LEAQ 0x966d(IP), AX
  main.go:10            0x4587f3                48890424                MOVQ AX, 0(SP)
  main.go:10            0x4587f7                48c744240805000000      MOVQ $0x5, 0x8(SP)
  main.go:10            0x458800                e85bb1faff              CALL runtime.makechan(SB)
  main.go:10            0x458805                488b442410              MOVQ 0x10(SP), AX
  main.go:10            0x45880a                4889442428              MOVQ AX, 0x28(SP)
  main.go:14            0x45882c                0f57c0                  XORPS X0, X0
  main.go:14            0x45882f                0f11442438              MOVUPS X0, 0x38(SP)
  main.go:14            0x458834                0f11442448              MOVUPS X0, 0x48(SP)
  main.go:14            0x458839                0f11442458              MOVUPS X0, 0x58(SP)
  main.go:14            0x45883e                488d055bd00000          LEAQ 0xd05b(IP), AX
  main.go:14            0x458845                48890424                MOVQ AX, 0(SP)
  main.go:14            0x458849                48c744240816000000      MOVQ $0x16, 0x8(SP)
  main.go:14            0x458852                488d442438              LEAQ 0x38(SP), AX
  main.go:14            0x458857                4889442410              MOVQ AX, 0x10(SP)
  main.go:14            0x45885c                e8ff28fbff              CALL runtime.makemap(SB)
  main.go:14            0x458861                488b442418              MOVQ 0x18(SP), AX
  main.go:14            0x458866                4889442420              MOVQ AX, 0x20(SP)
```

## 使⽤调试⼯具

`dlv`

官方文档: <https://github.com/go-delve/delve/tree/master/Documentation/cli>

1. 调试汇编时使⽤ si 到 JMP ⽬标位置
2. 使⽤ c(continue) 从⼀个断点到下⼀个断点
3. ⽤ disass 反汇编

```bash
[root@~]# dlv exec ./main
Type 'help' for list of commands.
(dlv) b *0x455780
Breakpoint 1 (enabled) set at 0x455780 for _rt0_amd64_linux() /usr/local/go/src/runtime/rt0_linux_amd64.s:8
(dlv) c
> _rt0_amd64_linux() /usr/local/go/src/runtime/rt0_linux_amd64.s:8 (hits total:1) (PC: 0x455780)
Warning: debugging optimized function
     3: // license that can be found in the LICENSE file.
     4:
     5: #include "textflag.h"
     6:
     7: TEXT _rt0_amd64_linux(SB),NOSPLIT,$-8
=>   8:         JMP     _rt0_amd64(SB)
     9:
    10: TEXT _rt0_amd64_linux_lib(SB),NOSPLIT,$0
    11:         JMP     _rt0_amd64_lib(SB)
(dlv) si
> _rt0_amd64() /usr/local/go/src/runtime/asm_amd64.s:15 (PC: 0x451bd0)
Warning: debugging optimized function
    10: // _rt0_amd64 is common startup code for most amd64 systems when using
    11: // internal linking. This is the entry point for the program from the
    12: // kernel for an ordinary -buildmode=exe program. The stack holds the
    13: // number of arguments and the C-style argv.
    14: TEXT _rt0_amd64(SB),NOSPLIT,$-8
=>  15:         MOVQ    0(SP), DI       // argc
    16:         LEAQ    8(SP), SI       // argv
    17:         JMP     runtime·rt0_go(SB)
    18:
    19: // main is common startup code for most amd64 systems when using
    20: // external linking. The C startup code will call the symbol "main"
(dlv) disass
TEXT _rt0_amd64(SB) /usr/local/go/src/runtime/asm_amd64.s
=>      asm_amd64.s:15  0x451bd0        488b3c24        mov rdi, qword ptr [rsp]
        asm_amd64.s:16  0x451bd4        488d742408      lea rsi, ptr [rsp+0x8]
        asm_amd64.s:17  0x451bd9        e902000000      jmp $runtime.rt0_go
```

## 语法实现分析

`go func`

```go
package main

import "time"

func main() {
	go func() {
		println("hello world!")
	}()

	time.Sleep(time.Second)
}
```

```bash
[root@~]# go tool compile -S ./main.go | grep "main.go:"
...
        0x001d 00029 (./main.go:6)      MOVL    $0, (SP)
        0x0024 00036 (./main.go:6)      PCDATA  $0, $1
        0x0024 00036 (./main.go:6)      LEAQ    "".main.func1·f(SB), AX
        0x002b 00043 (./main.go:6)      PCDATA  $0, $0
        0x002b 00043 (./main.go:6)      MOVQ    AX, 8(SP)
        0x0030 00048 (./main.go:6)      CALL    runtime.newproc(SB)
...
```

`channel send && recv`

```go
package main

func main() {
	var a = make(chan int, 1)
	a <- 100
	x := <-a
	println(x)
}
```

```bash
[root@~]# go tool compile -S main.go | grep -E "main.go:4|main.go:5|main.go:6"
...
        0x0035 00053 (main.go:4)        CALL    runtime.makechan(SB)
...
        0x0054 00084 (main.go:5)        CALL    runtime.chansend1(SB)
...
        0x0075 00117 (main.go:6)        CALL    runtime.chanrecv1(SB)
...
```

`channel 非阻塞 recv`

```go
package main

func main() {
	var a = make(chan int)
	select {
	case <-a:
	default:
	}
}
```

```bash
[root@~]# go tool compile -S main.go | grep -E "main.go:5|main.go:6|main.go:7"
        0x003b 00059 (main.go:6)        MOVQ    $0, (SP)
        0x0043 00067 (main.go:6)        PCDATA  $0, $0
        0x0043 00067 (main.go:6)        MOVQ    AX, 8(SP)
        0x0048 00072 (main.go:6)        CALL    runtime.selectnbrecv(SB)
        0x004d 00077 (main.go:6)        MOVQ    24(SP), BP
        0x0052 00082 (main.go:6)        ADDQ    $32, SP
        0x0056 00086 (main.go:6)        RET
        0x0057 00087 (main.go:6)        NOP
```

## Parser 应⽤场景示例

`内置 AST ⼯具-简单的规则引擎`

```text
初级会员，发贴数 > 10
中级会员，充值 > 1000 RMB
⾼级会员，发帖数 > 100，充值 > 10000 RMB
```

 ![Go相关分析 _6_.jpg](https://s2.loli.net/2023/05/27/4rFmViJun7SyI6D.png)

 ![qqqqq1.png](https://s2.loli.net/2023/05/27/ohlI5kX1MS7FLvf.png)

 ![222222907.png](https://s2.loli.net/2023/05/27/Nd5YRWSJlhpOcCA.png)

 ![1000.png](https://s2.loli.net/2023/05/27/l1AekZGLviHpuNz.png)

 ![16.png](https://s2.loli.net/2023/05/27/os2c8XewCqu6gIN.png)

 ![14.png](https://s2.loli.net/2023/05/27/UHqPvRD5wzjk37u.png)

 ![20.png](https://s2.loli.net/2023/05/27/rvOGUVnpbqmT7i3.png)

## 函数调⽤规约

 ![1.png](https://s2.loli.net/2023/05/27/mZwqNMzsPUvB6rx.png)

 ![2.png](https://s2.loli.net/2023/05/27/FaUPsVN8e6JKqiM.png)

`函数调⽤规约`

* The order in which atomic (scalar) parameters, or individual parts of a complex parameter, are allocated
* How parameters are passed (pushed on the stack, placed in registers, or a mix of both)
* Which registers the called function must preserve for the caller (also known as: callee-saved registers or non-volatile registers)
* How the task of preparing the stack for, and restoring after, a function call is divided between the caller and the callee

## References

* Go 的词法分析和语法/语义分析过程：<<https://dev.to/nakabonne/digging-deeper-into-the-analysis-of-go-code-31af>
* 编译器各阶段的简单介绍：
<https://www.tutorialspoint.com/compiler_design/compiler_design_phases_of_compiler.htm>
* Linkers and loaders，只看内部对 linker 的职责描述就⾏，不⽤看原理
<https://golearn.coding.net/p/gonggongbanji/files/all/DF9>
* SSA 的简单介绍：
<https://mp.weixin.qq.com/s/UhxFOQBpW8EUVpFvqH2tMg>
* ⽼外的写的如何定制 Go 编译器，⾥⾯对 Go 的编译过程介绍更详细，SSA 也说明得很好：
<https://eli.thegreenplace.net/2019/go-compiler-internals-adding-a-new-statement-to-go-part-2/>
* 如何阅读 go 的 SSA：
<https://sitano.github.io/2018/03/18/howto-read-gossa/>
* CMU 的编译器课，讲 SSA(*难，只做了解)
<https://www.cs.cmu.edu/~fp/courses/15411-f08/lectures/09-ssa.pdf>
* 对逆向感兴趣的话(扩展内容，与本课程⽆关)：
<https://malwareunicorn.org/#/workshops>
* Vitess 的 SQL Parser：
<https://github.com/xwb1989/sqlparser>
* PingCAP 的 TiDB 的 SQL Parser：
<https://github.com//pingcap/parser>
* GoCN 上的 dlv 的新译⽂
<https://gocn.vip/topics/12090>
* C语⾔调⽤规约
<https://github.com/cch123/llp-trans/blob/master/part3/translation-details/function-calling-sequence/callingconvention.md>
* Go 语⾔新版调⽤规约
<https://go.googlesource.com/proposal/+/refs/changes/78/248178/1/design/40724-register-calling.md>