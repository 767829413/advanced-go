package goroutineState

import "fmt"

type GoroutineState uint32

const (
	IDLE GoroutineState = iota
	RUNNABLE
	RUNNING
	SYSCALL
	WAITING
	MORIBUND_UNUSED
	DEAD
	ENQUEUE_UNUSED
	COPYSTACK
	PREEMPTED
)

func (s GoroutineState) String() string {
	switch s {
	case IDLE:
		return "IDLE"
	case RUNNABLE:
		return "RUNNABLE"
	case RUNNING:
		return "RUNNING"
	case SYSCALL:
		return "SYSCALL"
	case WAITING:
		return "WAITING"
	case MORIBUND_UNUSED:
		return "MORIBUND_UNUSED"
	case DEAD:
		return "DEAD"
	case ENQUEUE_UNUSED:
		return "ENQUEUE_UNUSED"
	case COPYSTACK:
		return "COPYSTACK"
	case PREEMPTED:
		return "PREEMPTED"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", s)
	}
}

// GoroutineExecuteData 与 C 语言的 struct goroutine_execute_data 对应
type GoroutineExecuteData struct {
	OldState GoroutineState
	NewState GoroutineState
	Goid     uint64
	Pid      uint32
	Tgid     uint32
}
