package connection

import (
	"sync"

	"github.com/labstack/echo/v4"
)

var (
	ReviewChannelSse  = make(chan string, 3)
	ActiveSseChannels sync.Map
)

func CheckActiveSseChannel() (count int, isEmpty bool) {
	c := 0
	ActiveSseChannels.Range(func(key, value any) bool {
		_, ok := key.(*echo.Response)
		if ok {
			c++
		}
		return true
	})
	return c, c == 0
}
