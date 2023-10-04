package messaging

import (
	"path"
	"strings"
)

// RouterGroup 路由组
type RouterGroup struct {
	prefix string
	router *Router
}

func (rg *RouterGroup) Group(prefix string) *RouterGroup {
	return &RouterGroup{prefix: path.Join(rg.prefix, prefix), router: rg.router}
}

func (rg *RouterGroup) Send(prefix string, handler MessageHandler) {
	rg.router.sendMap[path.Join(rg.prefix, prefix)] = handler
}
func (rg *RouterGroup) Subscribe(prefix string, handler MessageHandler) {
	rg.router.subscribeMap[path.Join(rg.prefix, prefix)] = handler
}

// Router 消息路由器
type Router struct {
	RouterGroup
	//send帧处理函数
	sendMap map[string]MessageHandler
	//subscribe帧处理函数
	subscribeMap map[string]MessageHandler
}

func NewRouter() *Router {
	r := &Router{
		sendMap:      make(map[string]MessageHandler),
		subscribeMap: make(map[string]MessageHandler),
	}
	rg := RouterGroup{
		router: r,
		prefix: "",
	}
	r.RouterGroup = rg
	return r
}

func (router *Router) getMessageHandler(context *Context) MessageHandler {
	destArr := strings.Split(context.Frame.Headers["destination"], "/")
	var hFunc MessageHandler
	for key, value := range router.sendMap {
		keyArr := strings.Split(key, "/")
		if len(destArr) != len(keyArr) {
			continue
		}
		match := true
		for i := 0; i < len(keyArr); i++ {
			if strings.HasPrefix(keyArr[i], ":") {
				context.Params[strings.TrimPrefix(keyArr[i], ":")] = destArr[i]
				continue
			}
			if keyArr[i] != destArr[i] {
				match = false
				break
			}
		}
		if match {
			hFunc = value
			break
		}
	}

	for key, value := range router.subscribeMap {
		keyArr := strings.Split(key, "/")
		if len(destArr) != len(keyArr) {
			continue
		}
		match := true
		for i := 0; i < len(keyArr); i++ {
			if strings.HasPrefix(keyArr[i], ":") {
				context.Params[strings.TrimPrefix(keyArr[i], ":")] = destArr[i]
				continue
			}
			if keyArr[i] != destArr[i] {
				match = false
				break
			}
		}
		if match {
			hFunc = value
			break
		}
	}
	return hFunc
}
