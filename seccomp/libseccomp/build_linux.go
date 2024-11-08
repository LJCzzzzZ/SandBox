package libseccomp

import (
	"Sandbox/seccomp"
	"syscall"

	libseccomp "github.com/elastic/go-seccomp-bpf"
	"golang.org/x/net/bpf"
)

type Builder struct {
	Allow, Trace []string
	Default      Action
}

var actTrace = libseccomp.ActionTrace

func (b *Builder) Build() (seccomp.Filter, error) {
	policy := libseccomp.Policy{
		DefaultAction: ChooseDefaultAction(b.Default),
		Syscalls: []libseccomp.SyscallGroup{
			{
				Action: libseccomp.ActionAllow,
				Names:  b.Allow,
			},
			{
				Action: actTrace,
				Names:  b.Trace,
			},
		},
	}
	program, err := policy.Assemble()
	if err != nil {
		return nil, err
	}
	return ExportBPF(program)
}

func ExportBPF(filter []bpf.Instruction) (seccomp.Filter, error) {
	raw, err := bpf.Assemble(filter)
	if err != nil {
		return nil, err
	}
	return sockFilter(raw), nil
}

func sockFilter(raw []bpf.RawInstruction) []syscall.SockFilter {
	filter := make([]syscall.SockFilter, 0, len(raw))
	for _, instruction := range raw {
		filter = append(filter, syscall.SockFilter{
			Code: instruction.Op,
			Jt:   instruction.Jt,
			Jf:   instruction.Jf,
			K:    instruction.K,
		})
	}
	return filter
}

func ChooseDefaultAction(a Action) libseccomp.Action {
	var action libseccomp.Action
	switch a.Action() {
	case ActionAllow:
		action = libseccomp.ActionAllow
	case ActionErrno:
		action = libseccomp.ActionErrno
	case ActionTrace:
		action = libseccomp.ActionTrace
	default:
		action = libseccomp.ActionKillProcess
	}
	return action
}