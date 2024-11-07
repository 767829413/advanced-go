// SPDX-License-Identifier: GPL-2.0 OR BSD-3-Clause
/* 版权所有 (c) 2021 Sartura */
#define BPF_NO_GLOBAL_DATA
#include "vmlinux.h"
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>
#include <bpf/bpf_core_read.h>
#include <bpf/bpf_endian.h>

// 声明许可证类型为 "Dual BSD/GPL"
char LICENSE[] SEC("license") = "Dual BSD/GPL";

// 定义一个 eBPF 程序，附加到 do_unlinkat 函数的入口点
SEC("kprobe/do_unlinkat")
int BPF_KPROBE(do_unlinkat, int dfd, struct filename *name)
{
    pid_t pid;
    const char *filename;

    // 获取当前进程的 PID
    pid = bpf_get_current_pid_tgid() >> 32;
    // 读取文件名
    filename = BPF_CORE_READ(name, name);
    // 输出调试信息（进入 kprobe）
    bpf_printk("KPROBE ENTRY pid = %d, filename = %s\n", pid, filename);
    return 0;
}

// 定义一个 eBPF 程序，附加到 do_unlinkat 函数的退出点
SEC("kretprobe/do_unlinkat")
int BPF_KRETPROBE(do_unlinkat_exit, long ret)
{
    pid_t pid;

    // 获取当前进程的 PID
    pid = bpf_get_current_pid_tgid() >> 32;
    // 输出调试信息（退出 kprobe）
    bpf_printk("KPROBE EXIT: pid = %d, ret = %ld\n", pid, ret);
    return 0;
}