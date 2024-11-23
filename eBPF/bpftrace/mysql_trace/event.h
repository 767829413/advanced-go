#ifndef EBPF_EXAMPLE_EVENT_H
#define EBPF_EXAMPLE_EVENT_H

// 定义发送给用户空间的数据结构
struct event {
    u32 pid;
    u32 tid;
    u64 delta_ns;
    char query[256]; // 假设最大查询长度为 256 字节
};

#endif