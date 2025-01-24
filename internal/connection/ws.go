package connection

import (
	"sync"

	"github.com/gorilla/websocket"
)

var (
	ReviewChannelWs  = make(chan string, 3)
	ActiveWsChannels sync.Map
)

func CheckActiveWsChannel() (count int, isEmpty bool) {
	c := 0
	ActiveWsChannels.Range(func(key, value any) bool {
		_, ok := key.(*websocket.Conn)
		if ok {
			c++
		}
		return true
	})
	return c, c == 0
}
