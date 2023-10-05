package messaging

import (
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"
)

func TestMessaging(t *testing.T) {
	broker := NewStompBroker()
	broker.AppDestinationPrefix = "/app"
	broker.BrokerDestinationPrefix = "/topic"
	broker.Send("/app/test", func(ctx *Context) {
		fmt.Println(ctx.Frame.String())
		ctx.Send("/topic/test", []byte("你好"))
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		broker.ServeOverHttp(w, r)
	})
	server := &http.Server{
		Addr:              ":8080",
		ReadHeaderTimeout: 3 * time.Second,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
