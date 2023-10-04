package messaging

import (
	"errors"
	"github.com/gorilla/websocket"
	"net/http"
	"strings"
)

// Upgrader 协议升级器
type Upgrader struct {
	wsUpgrader websocket.Upgrader
}

// Upgrade 从http连接升级到stomp连接
func (u *Upgrader) Upgrade(w http.ResponseWriter, r *http.Request) (*Conn, error) {
	//建立websocket连接
	wsConn, err := u.wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	conn := NewConn(wsConn)

	//读取CONNECT/STOMP帧
	frame, err := conn.ReadFrame()
	if err != nil {
		return nil, err
	}
	//不是CONNECT帧
	if frame.Command != CONNECT && frame.Command != STOMP {
		conn.Write(
			ERROR,
			[]byte(frame.String()),
			"message", "can not sent other frame before CONNECT",
		)
		return nil, errors.New("not allowed frame before CONNECT")
	}

	//开始协议协商
	v, ok := frame.Headers["accept-version"]

	if !ok { //没有accept version头
		conn.Write(
			ERROR,
			[]byte(frame.String()),
			"message", "CONNECT frame must has an accept-version header",
		)
		return nil, errors.New("CONNECT frame must has an accept-version header")
	}

	//协议版本协商，目前只支持1.2
	cVersions := strings.Split(v, ",")
	var support bool
	for _, version := range cVersions {
		if version == "1.2" {
			support = true
		}
	}

	if !support {
		conn.Write(
			ERROR,
			[]byte(frame.String()),
			"message", "unsupported protocol version, the supported version is 1.2",
		)
		return nil, errors.New("unsupported protocol version")
	}

	//连接成功
	conn.Write(
		CONNECTED,
		nil,
		"version", "1.1",
		"heart-beat", "0,0",
	)
	return conn, nil
}
