#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>

char LICENSE[] SEC("license") = "Dual BSD/GPL";

// 定义一个简单的数据结构来存储事件数据
struct event {
    u32 pid;
    u8 comm[80];
};

// 定义一个 BPF map 来存储事件数据
struct {
    __uint(type, BPF_MAP_TYPE_PERF_EVENT_ARRAY);
    __uint(key_size, sizeof(u32));
    __uint(value_size, sizeof(u32));
} events SEC(".maps");

// kprobe 程序入口点
SEC("kprobe/__x64_sys_execve")
int kprobe_exec(struct pt_regs *ctx)
{
    struct event event = {};

    // 获取当前进程 ID
    event.pid = bpf_get_current_pid_tgid() >> 32;

    // 获取当前进程名称
    bpf_get_current_comm(&event.comm, sizeof(event.comm));

    // 将事件数据发送到用户空间
    bpf_perf_event_output(ctx, &events, BPF_F_CURRENT_CPU, &event, sizeof(event));

    return 0;
}