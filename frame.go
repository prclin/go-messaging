package messaging

import (
	"errors"
	"fmt"
	"strings"
)

// command常量,不能打乱顺序
const (
	MESSAGE Command = iota
	RECEIPT
	ERROR

	CONNECT
	STOMP
	CONNECTED
	SEND
	SUBSCRIBE
	UNSUBSCRIBE
	ACK
	NACK
	BEGIN
	COMMIT
	ABORT
	DISCONNECT
)

// 所有支持的command
var commands = []string{"MESSAGE", "RECEIPT", "ERROR", "CONNECT", "STOMP", "CONNECTED", "SEND", "SUBSCRIBE", "UNSUBSCRIBE", "ACK", "NACK", "BEGIN", "COMMIT", "ABORT", "DISCONNECT"}

// 确保Command实现了Stringer
var _ fmt.Stringer = Command(0)

// Command stomp帧的指令
type Command int8

// CommandOf 通过字符串值，获取command枚举
func CommandOf(command string) (Command, error) {
	for i, v := range commands {
		if v == command {
			return Command(i), nil
		}
	}
	return -1, errors.New("unsupported command")
}

func (c Command) String() string {
	return commands[c]
}

// Frame stomp协议帧
//
// Only the SEND, MESSAGE, and ERROR frames can have a body. All other frames MUST NOT have a body.
type Frame struct {
	Command Command
	Headers map[string]string
	Body    []byte
}

func NewFrame(command Command, headers map[string]string, body []byte) *Frame {
	return &Frame{Command: command, Headers: headers, Body: body}
}

func (f *Frame) String() string {
	var sb strings.Builder
	sb.WriteString(f.Command.String())
	sb.WriteString("\n")
	s := ""
	//command
	s += f.Command.String() + "\n"
	//headers
	for key, value := range f.Headers {
		sb.WriteString(key)
		sb.WriteString(":")
		sb.WriteString(value)
		sb.WriteString("\n")
	}
	//blank line
	sb.WriteString("\n")

	//body
	sb.Write(f.Body)
	sb.Write([]byte{0x00}) //结尾空八位

	return sb.String()
}
