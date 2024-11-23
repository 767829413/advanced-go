#include "vmlinux.h"
#include "event.h"
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_core_read.h>
#include <bpf/bpf_tracing.h>


char LICENSE[] SEC("license") = "Dual BSD/GPL";

// 定义开始时间的哈希表，键为线程 ID（TID），值为开始时间
struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, 10240);
    __type(key, u32);
    __type(value, u64);
} start_times SEC(".maps");

// 定义查询字符串的哈希表，键为线程 ID（TID），值为查询字符串
struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, 10240);
    __type(key, u32);
    __type(value, char[256]);
} comm_sql SEC(".maps");

// 定义 ringbuffer，用于向用户空间传递数据
struct {
    __uint(type, BPF_MAP_TYPE_RINGBUF);
    __uint(max_entries, 1 << 24); // 16 MiB
} events SEC(".maps");



// uprobe：在 dispatch_command 函数入口时触发
SEC("uprobe//var/lib/docker/overlay2/270506fa70024c98a81543a6a274e9f03258f39471655b69fa22fce5fada6d13/merged/sbin/mysqld:_Z16dispatch_commandP3THDPK8COM_DATA19enum_server_command")
int uprobe_mysql_dispatch_command(struct pt_regs *ctx) {
    u32 tid = bpf_get_current_pid_tgid() >> 32;
    u64 ts = bpf_ktime_get_ns();

    void* st = (void*) PT_REGS_PARM2(ctx);
    char query[256];
    char* query_ptr;
    bpf_probe_read_user(&query_ptr, sizeof(query_ptr), st);
    bpf_probe_read_user_str(query, sizeof(query), query_ptr);

    bpf_printk("slowsql: query=%s, query_ptr=%s\n", query,query_ptr);

    // 记录comm_sql
    bpf_map_update_elem(&comm_sql, &tid, query, BPF_ANY);

    // 记录开始时间
    bpf_map_update_elem(&start_times, &tid, &ts, BPF_ANY);

    return 0;
}

// uretprobe：在 dispatch_command 函数返回时触发
SEC("uretprobe//var/lib/docker/overlay2/270506fa70024c98a81543a6a274e9f03258f39471655b69fa22fce5fada6d13/merged/sbin/mysqld:_Z16dispatch_commandP3THDPK8COM_DATA19enum_server_command")
int uretprobe_mysql_dispatch_command(struct pt_regs *ctx) {
    u32 tid = bpf_get_current_pid_tgid() >> 32;

    char *query = bpf_map_lookup_elem(&comm_sql, &tid);
    if (!query) {
        // 未找到查询字符串，可能是因为在进入时未记录
        return 0;
    }

    // 删除查询字符串记录，防止内存泄漏
    bpf_map_delete_elem(&comm_sql, &tid);



    u64 *start_ts = bpf_map_lookup_elem(&start_times, &tid);
    if (!start_ts) {
        // 未找到开始时间，可能是因为在进入时未记录
        return 0;
    }

    u64 delta = bpf_ktime_get_ns() - *start_ts;
    // 删除开始时间记录，防止内存泄漏
    bpf_map_delete_elem(&start_times, &tid);

    // 如果延迟小于等于 10 毫秒（10000000 纳秒），则不处理
    if (delta <= 10000000) {
        return 0;
    }

    // 分配事件数据
    struct event *e;
    e = bpf_ringbuf_reserve(&events, sizeof(*e), 0);
    if (!e) {
        // 分配失败
        return 0;
    }

    // 填充事件数据
    e->pid = bpf_get_current_pid_tgid() & 0xFFFFFFFF;
    e->tid = tid;
    e->delta_ns = delta;

    if (query) {
        bpf_probe_read_kernel_str(e->query, sizeof(e->query), query);
    } else {
        // 未能获取查询字符串
        e->query[0] = '\0';
    }

    bpf_printk("slowsql: pid=%d, tid=%d, delta=%lld, query=%s\n", e->pid, e->tid, e->delta_ns, e->query);

    // 提交事件到 ringbuffer
    bpf_ringbuf_submit(e, 0);


    return 0;
}