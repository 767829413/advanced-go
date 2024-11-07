#ifndef EBPF_EXAMPLE_GOROUTINE_H
#define EBPF_EXAMPLE_GOROUTINE_H

enum goroutine_state
{
    IDLE,
    RUNNABLE,
    RUNNING,
    SYSCALL,
    WAITING,
    MORIBUND_UNUSED,
    DEAD,
    ENQUEUE_UNUSED,
    COPYSTACK,
    PREEMPTED,
};

struct goroutine_execute_data
{
    enum goroutine_state old_state;
    enum goroutine_state new_state;
    u64 goid;
    u32 pid;
    u32 tgid;
};

#endif