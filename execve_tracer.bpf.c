#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>
#include <linux/ptrace.h>
#include <linux/sched.h>
#include <linux/types.h>

#define TASK_COMM_LEN 16
struct data_t {
	__u32 pid;
	__u32 uid;
	char comm[TASK_COMM_LEN];
};

struct {
	__uint(type, BPF_MAP_TYPE_PERF_EVENT_ARRAY);
} events SEC(".maps");

SEC("kprobe.sys_execve")
int catch_success_cmd(struct pt_regs *ctx) {
	struct data_t data = {};
	__u32 pid = bpf_get_current_pid_tgid() >> 32;
	__u32 uid = bpf_get_current_uid_gid();

	bpf_get_current_comm(&data.comm, sizeof(data.comm));
	data.pid = pid;
	data.uid = uid;

	bpf_perf_event_output(ctx, &events, BPF_F_CURRENT_CPU, &data, sizeof(data));
	return 0;
}

char LICENSE[] SEC("license") = "GPL";
