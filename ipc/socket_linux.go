package ipc

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"syscall"
)

// buff maxSize
const msgSize = 4 << 10

type Socket struct {
	*net.UnixConn
	sendBuff []byte
	recvBuff []byte
}

type Message struct {
	Fds  []int
	Cred *syscall.Ucred
}

func newSocket(conn *net.UnixConn) *Socket {
	return &Socket{
		UnixConn: conn,
		sendBuff: make([]byte, msgSize),
		recvBuff: make([]byte, msgSize),
	}
}

func NewSocket(fd int) (*Socket, error) {
	syscall.SetNonblock(fd, true)
	syscall.CloseOnExec(fd)

	file := os.NewFile(uintptr(fd), "unix-socket")
	if file == nil {
		return nil, fmt.Errorf("NewSocket: %d is not a vaild fd", fd)
	}
	defer file.Close()

	conn, err := net.FileConn(file)
	if err != nil {
		return nil, err
	}

	UnixConn, ok := conn.(*net.UnixConn)
	if !ok {
		conn.Close()
		return nil, fmt.Errorf("NewSocket: %d is not a valid unix socket connection", fd)
	}
	return newSocket(UnixConn), nil
}

func NewSocketPair() (*Socket, *Socket, error) {
	fd, err := syscall.Socketpair(syscall.AF_LOCAL, syscall.SOCK_SEQPACKET|syscall.SOCK_CLOEXEC, 0)
	if err != nil {
		return nil, nil, fmt.Errorf("NewSocketPair: failed to call sockerpair %v", err)
	}

	ins, err := NewSocket(fd[0])
	if err != nil {
		syscall.Close(fd[0])
		syscall.Close(fd[1])
		return nil, nil, fmt.Errorf("NewSocketPair: failed to call sockerpair on sender %v", err)
	}
	outs, err := NewSocket(fd[1])
	if err != nil {
		syscall.Close(fd[0])
		syscall.Close(fd[1])
		return nil, nil, fmt.Errorf("NewSocketPair: failed to call sockerpair on receiver %v", err)
	}

	return ins, outs, err
}

func (s *Socket) SetPassCred(option int) error {
	sysconn, err := s.SyscallConn()
	if err != nil {
		return err
	}
	return sysconn.Control(func(fd uintptr) {
		syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_PASSCRED, option)
	})
}

func (s *Socket) SendMsg(b []byte, m Message) error {
	msg := bytes.NewBuffer(s.sendBuff[:0])
	if len(m.Fds) > 0 {
		msg.Write(syscall.UnixRights(m.Fds...))
	}

	if m.Cred != nil {
		msg.Write(syscall.UnixCredentials(m.Cred))
	}
	_, _, err := s.WriteMsgUnix(b, msg.Bytes(), nil)
	if err != nil {
		return err
	}
	return nil
}

func (s *Socket) Recvmsg(b []byte) (int, Message, error) {
	var message Message
	n, oobn, _, _, err := s.ReadMsgUnix(b, s.recvBuff)
	if err != nil {
		return 0, message, err
	}
	msgs, err := syscall.ParseSocketControlMessage(s.recvBuff[:oobn])
	if err != nil {
		return 0, message, err
	}
	// msg, err := parseMsg
}
