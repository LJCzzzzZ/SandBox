package libseccomp

// Action represents a seccomp action applied to a syscall.
// Default value 0 is invalid.
type Action uint32

// Predefined seccomp actions.
const (
	ActionAllow Action = iota + 1 // Allow the syscall
	ActionErrno                   // Return an error
	ActionTrace                   // Trace the syscall
	ActionKill                    // Kill the process
)

// MsgType defines the type of message triggered by seccomp filters.
type MsgType int16

// Predefined message types for trapped syscalls.
const (
	MsgDisallow MsgType = iota + 1 // Syscall is disallowed
	MsgHandle                      // Syscall requires further handling
)

func (a Action) Action() Action {
	return Action(a & 0xffff)
}
