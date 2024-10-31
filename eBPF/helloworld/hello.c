/* SPDX-License-Identifier: (LGPL-2.1 OR BSD-2-Clause) */ // SPDX 许可证标识符，指定双重许可选项：LGPL-2.1 或 BSD-2-Clause
// 定义宏以禁用 BPF 程序中的全局数据使用
#define BPF_NO_GLOBAL_DATA
// 包含 Linux BPF 头文件，用于 BPF 程序定义和辅助函数                             
#include <linux/bpf.h>
// 包含 BPF 框架提供的辅助函数                               
#include <bpf/bpf_helpers.h>
// 包含 BPF 跟踪辅助函数，用于附加到跟踪点                           
#include <bpf/bpf_tracing.h> 

typedef unsigned int u32;   // 定义无符号整数类型别名 u32
typedef int pid_t;          // 定义整数类型别名 pid_t
const pid_t pid_filter = 0; // 定义用于 PID 过滤的常量，0 表示不进行过滤

char LICENSE[] SEC("license") = "Dual BSD/GPL"; // 指定 BPF 程序的许可证，允许双重 BSD/GPL 许可

// 定义一个 BPF map 来存储上次输出的时间
struct {
    __uint(type, BPF_MAP_TYPE_ARRAY);
    __uint(max_entries, 1);
    __type(key, u32);
    __type(value, __u64);
} last_output SEC(".maps");

// 定义 sys_enter_write 跟踪点的 BPF 程序部分
SEC("tp/syscalls/sys_enter_write") 
// 定义处理该跟踪点的 BPF 程序函数
int handle_tp(void *ctx)
{
    pid_t pid = bpf_get_current_pid_tgid() >> 32;
    if (pid_filter && pid != pid_filter)
        return 0;

    __u64 ts = bpf_ktime_get_ns();
    u32 key = 0;
    __u64 *last_ts = bpf_map_lookup_elem(&last_output, &key);
    
    if (!last_ts) {
        // 如果是第一次运行，直接输出并更新时间戳
        bpf_map_update_elem(&last_output, &key, &ts, BPF_ANY);
        bpf_printk("BPF triggered sys_enter_write from PID %d.\n", pid);
    } else if (ts - *last_ts >= 1000000000) {  // 1秒 = 1,000,000,000 纳秒
        // 如果距离上次输出已经过了至少1秒，则输出并更新时间戳
        bpf_map_update_elem(&last_output, &key, &ts, BPF_ANY);
        bpf_printk("BPF triggered sys_enter_write from PID %d.\n", pid);
    }

    return 0;
}