package messaging

import (
	"bufio"
	"bytes"
	"errors"
	"strings"
)

// ProtocolResolver 协议解析器
type ProtocolResolver interface {
	//Parse 将二进制数据转换为stomp帧
	Parse(buf []byte) (*Frame, error)
}

var defaultProtocolResolver ProtocolResolver = &StandardProtocolResolver{}

// StandardProtocolResolver 标准协议解析器
type StandardProtocolResolver struct {
}

func (str StandardProtocolResolver) Parse(buf []byte) (*Frame, error) {
	//获取reader
	bufReader := bufio.NewReader(bytes.NewReader(buf))

	//读取command
	commandStr, err := bufReader.ReadString('\n')
	if err != nil {
		return nil, errors.New("command must exist")
	}
	commandStr = strings.TrimSuffix(commandStr, "\n") //去除后缀
	command, err := CommandOf(commandStr)
	if err != nil {
		return nil, err
	}

	//读取header
	headers := make(map[string]string, 2) //一般的帧都有两个header
	for {
		header, err := bufReader.ReadString('\n')
		if err != nil {
			return nil, errors.New("there must be a blank line between header and body")
		}
		header = strings.TrimSuffix(header, "\n")
		if header == "" { //读取到空行
			break
		}

		kv := strings.SplitN(header, ":", 2) //获取key、value
		if len(kv) != 2 {                    //header格式是否合法
			return nil, errors.New("wrong header format")
		}
		if _, ok := headers[kv[0]]; !ok { //如果有重复的header，保留第一个
			headers[kv[0]] = kv[1]
		}
	}

	//读取body
	body, err := bufReader.ReadBytes(0x00)
	if err != nil {
		return nil, errors.New("body must be end with a NULL octet")
	}

	return NewFrame(command, headers, body), nil
}
