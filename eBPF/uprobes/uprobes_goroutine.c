#include "vmlinux.h"
#include "goroutine.h"
#include <bpf/bpf_core_read.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>

const volatile int pid_target = 0;

#define GOID_OFFSET 0x98

struct {
  __uint(type, BPF_MAP_TYPE_RINGBUF);
  __uint(max_entries, 256 * 1024);
} rb SEC(".maps");

const struct goroutine_execute_data *unused __attribute__((unused));

SEC("uprobe//home/fangyuan/code/go/src/github.com/767829413/advanced-go/eBPF/uprobes/demo/main:runtime.casgstatus")
int uprobe_runtime_casgstatus(struct pt_regs *ctx) {
    void *gp = (void *)ctx->ax;
    u32 oldval = (u32)ctx->bx;
    u32 newval = (u32)ctx->cx;

    u64 tgid_pid = bpf_get_current_pid_tgid();
    u32 pgid = tgid_pid >> 32;
    u32 pid = tgid_pid;

    if (pid_target && pid_target != pgid)
      return false;


    struct goroutine_execute_data *data;
    u64 goid;
    if (bpf_probe_read_user(&goid, sizeof(goid), gp + GOID_OFFSET) == 0) {
      data = bpf_ringbuf_reserve(&rb, sizeof(*data), 0);
      if (data) {
        data->pid = pid;
        data->tgid = pgid;
        data->goid = goid;
        data->new_state = newval;
        data->old_state = oldval; // Assuming you have added this field to the struct
        bpf_printk("pgid:%d, pid:%d, goid: %lu, old state: %d, new state: %d", data->tgid, data->pid, goid, oldval, newval);
        bpf_ringbuf_submit(data, 0);
      }
    }
    return 0;
}

char LICENSE[] SEC("license") = "GPL";