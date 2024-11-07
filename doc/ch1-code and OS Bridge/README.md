# 编程语言与操作系统桥梁

## 什么是系统调⽤

`什么是操作系统`

**操作系统是资源的管理器，其管理的资源均进⾏了抽象**

* 磁盘抽象：⽂件夹
* 内存抽象：虚拟内存
* CPU 抽象：时间⽚

 ![os_some.png](https://s2.loli.net/2023/06/12/9HyoD8FqGvQtxf6.png)

`分级保护域-protection ring`

**CPU 为操作系统提供了特殊的安全⽀持操作系统内核运⾏在特殊模式下即图中的 ring-0 应⽤运⾏在 ring-3，权限被严格限制**

 ![os_level.jpg](https://s2.loli.net/2023/06/12/NYtFMRH28yQAKI1.png)

*Intel 64 有四个特权级别，不过实际上只⽤到了其中的两个：ring-0 和 ring-3. *
*ring-1 ring-2本来计划是为驱动程序和 OS 服务⽤，不过流⾏的 OS 们都没有接受这个⽅案. *

`什么是系统调⽤`

**系统调⽤是操作系统内核为应⽤提供的 API,是内核为应⽤提供的服务,操作系统为上层的应⽤程序提供了⼀个“标准库”**

 ![os_call_1.png](https://s2.loli.net/2023/06/12/kBFwsJAD4NxGP8i.png)

*对于应⽤来说，系统调⽤可以实现超出⾃⼰能⼒以外的事情*

**Go(1.14) 语⾔调⽤规约中未使⽤寄存器**

 ![os_call_go.png](https://s2.loli.net/2023/06/12/brWdz7uDfBCTFha.png)

**寄存器: CPU内部的特殊存储单元**

 ![os_call_2.png](https://s2.loli.net/2023/06/12/ZbsYL7lPfmOk3zG.png)

 参考资料: <https://github.com/cch123/llp-trans/blob/master/part1/basic-computer-architecture/registers.md>

**系统调⽤有⾃⼰的⼀套调⽤规约，需要使⽤寄存器和 C 语⾔的调⽤规约相似**

 | arch | syscall NR | return | arg0 | arg1 | arg2 | arg3 | arg4 | arg5 |
 | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: |
 | arm | r7 | r0 | r0 | r1 | r2 | r3 | r4 | r5 |
 | arm64 | x8 | x0 | x0 | x1 | x2 | x3 | x4 | x5 |
 | x86 | eax | eax | ebx | ecx | edx | esi | edi | ebp |
 | x86_64 | rax | rax | rdi | rsi | rdx | r10 | r8 | r9 |
  
 参考资料: <https://chromium.googlesource.com/chromiumos/docs/+/master/constants/syscalls.md#x86_64-64_bit>

**系统调⽤举例**

 ![os_call_demo.jpg](https://s2.loli.net/2023/06/12/hYCAyGQdqTSmBuI.png)

**SYSCALL 之后发⽣了什么**

 ![demo2.png](https://s2.loli.net/2023/06/12/3SJUGABDpwvPWM8.png)

 内核代码查找: <https://code.woboq.org/linux/linux/fs/eventpoll.c.html#do_epoll_create>

## 常⻅系统调⽤

`go的一些常见系统调用`

```go
// getpid()
os.getpid()

// write(2,"hello world",11)
// arg0: stderr arg2: strlen
println("Hello world")

// clone(child_stack=0xc42003c000,flags=....)
startm().newm().newosproc()
```

 | Types of System Calls | Windows | Linux |
 | :---: | :---: | :---: |
 | Process Control | CreateProcess() ExitProcess() WaitForSingleObject() | fork() exit() wait() |
 | File Management | CreateFile() ReadFile() WriteFile() CloseHandle() | open() read() write() close() |
 | Device Management | SetConsoleMode() ReadConsole() WriteConsole() | ioctl() read() write() |
 | Information Maintenance | GetCurrentProcessID() SetTimer() Sleep() | getpid() alarm() sleep() |
 | Communication | CreatePipe() CreateFileMapping() MapViewOfFile() | pipe() shmget() mmap() |

## 观察系统调⽤

`strace on linux, dtruss on macOS`

```bash
# -c 统计系统调用, -f follow forks,详细用法可以 -h 
➜  advanced-go git:(main) ✗ strace  ./test            
execve("./test", ["./test"], 0x7fff3786fce0 /* 37 vars */) = 0
arch_prctl(ARCH_SET_FS, 0x4bfcd0)       = 0
sched_getaffinity(0, 8192, [0, 1, 2, 3, 4, 5, 6, 7]) = 32
openat(AT_FDCWD, "/sys/kernel/mm/transparent_hugepage/hpage_pmd_size", O_RDONLY) = 3
read(3, "2097152\n", 20)                = 8
close(3)                                = 0
...

➜  advanced-go git:(main) ✗ strace -c ./test        
0 0 false
% time     seconds  usecs/call     calls    errors syscall
------ ----------- ----------- --------- --------- ----------------
 65.53    0.005507          48       114           rt_sigaction
 15.97    0.001342          53        25           rt_sigreturn
  6.87    0.000577          72         8           rt_sigprocmask
  4.60    0.000387          64         6           write
  3.80    0.000319         106         3           clone
  1.61    0.000135          67         2           sigaltstack
  0.84    0.000071          71         1           futex
  0.79    0.000066          66         1           gettid
  0.00    0.000000           0         1           read
  0.00    0.000000           0         1           close
  0.00    0.000000           0        18           mmap
  0.00    0.000000           0         1           execve
  0.00    0.000000           0         1           uname
  0.00    0.000000           0         1           arch_prctl
  0.00    0.000000           0         1           sched_getaffinity
  0.00    0.000000           0         1           openat
------ ----------- ----------- --------- --------- ----------------
100.00    0.008404                   185           total
```

**strace 的实现依赖了 ptrace 这个 syscall 调试器(如 delve)也是⼤量使⽤了 ptrace**

 ![strace.png](https://s2.loli.net/2023/06/12/D7eQawG4PMRAHxN.png)

## Go 语⾔中的系统调⽤

`阻塞和⾮阻塞的系统调⽤`

 ![sys_call.png](https://s2.loli.net/2023/06/12/OZx51LsCcwVW92R.png)

 ![sys_call_1.png](https://s2.loli.net/2023/06/12/2TjPrfkWpnEA6NJ.png)

 **6 其实说的是 6 个参数,很多系统接⼝也有类似的命名法，如 wait4，accept4**

`Syscall 相关代码的基本结构`

* OS 相关的基础⽂件，在 syscall package 中：<https://golang.org/src/syscall/syscall_linux.go>
* 使⽤脚本⽣成的⽂件，在 syscall package 中：<https://golang.org/src/syscall/zsyscall_linux_386.go>
* 不对⽤户暴露的特殊 syscall，不受调度影响，在 runtime 中：<https://golang.org/src/runtime/sys_linux_amd64.s>

**阻塞的系统调⽤需要修改 P 的状态：running -> syscall. 这样在 sysmon 中才能发现这个 P 已经在 syscall 状态阻塞了**

 ![sys_call_2.png](https://s2.loli.net/2023/06/12/RqOy7aEDWmsb5iA.png)

 演示动画: <https://www.figma.com/proto/ounOboEYjlzBwcOhPgE2Z5/syscall?node-id=11-3&starting-point-node-id=11%3A3>

**VDSO 优化,内核负责，⾃动映射值到⽤户地址空间,⽆需⽤户/内核态切换**

```text
 The "vDSo" (virtual dynamic shared object) is a small shared library that the kernel automatically maps into the address space of all user-space applications.
```

 ![vdso.png](https://s2.loli.net/2023/06/12/naxTCvwLi3rBPRc.png)

## 系统调⽤发展历史

 ![call_history.png](https://s2.loli.net/2023/06/12/pqRszOTSI9uQlV3.png)

 Intel x86 vs x64 system call: <https://stackoverflow.com/questions/15168822/intel-x86-vs-x64-system-call>

## 学习系统调⽤

* 系统调用参考: <https://man7.org/>
* 参考书籍: 《The Linux Programming Interface》<https://sciencesoftcode.files.wordpress.com/2018/12/the-linux-programming-interface-michael-kerrisk-1.pdf>
* 本地的 man 命令

**这些知识到底有什么⽤？**

* 了解操作系统与应⽤程序的边界
* 了解内核升级导致应⽤程序⾏为变化的原因，如：madvdontneed 修改时，导致线上应⽤ RSS ⼤幅上升
* 系统⽆响应时，观察系统是卡死在什么地⽅(也许根本没卡死，只是你的错觉)

## References

* Ring 0 和 Ring 3: 
  * <https://zh.wikipedia.org/wiki/%E5%88%86%E7%BA%A7%E4%BF%9D%E6%8A%A4%E5%9F%9F>
  * <https://stackoverflow.com/questions/18717016/what-are-ring-0-and-ring-3-in-the-context-of-operating-systems>
  * <https://www.futurelearn.com/info/courses/computer-systems/0/steps/53514>
* 寄存器的概念： <https://github.com/cch123/llp-trans/blob/d6b7f46c72e83ac9145d5534c6bc4e690da8d815/part1/basic-computer-architecture/registers.md>
* 硬件中断、软件中断：<https://www.cs.montana.edu/courses/spring2005/518/Hypertextbook/jim/media/interrupts_on_linux.pdf>
* Int 80 和 syscall 有什么区别: <https://reverseengineering.stackexchange.com/questions/16702/difference-between-int-0x80-and-syscall>
* 为什么 sysenter/syscall ⽐ int 80 开销⼩：<https://stackoverflow.com/questions/15168822/intel-x86-vs-x64-system-call>
* 这⼀篇把内核中的代码返回位置都说清楚了，不需要⾃⼰去翻代码了：<https://blog.packagecloud.io/eng/2016/04/05/the-definitive-guide-to-linux-system-calls/>
* Strace 在 docker ⾥做实验时，可能要做⼀点配置：<https://jvns.ca/blog/2020/04/29/why-strace-doesnt-work-in-docker/>
* Anatomy of a system call
  * <https://lwn.net/Articles/604287/>
  * <https://lwn.net/Articles/604515/>
* VDSO: <https://man7.org/linux/man-pages/man7/vdso.7.html>
* time.Now 的性能衰退，其中有讲到 vdso：<https://mp.weixin.qq.com/s/06SDQLzDprJf2AEaDnX-QQ>
* vdso 和普通 syscall 的性能对⽐：<http://arkanis.de/weblog/2017-01-05-measurements-of-system-call-performance-and-overhead>