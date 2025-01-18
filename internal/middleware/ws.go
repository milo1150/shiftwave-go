package middleware

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	Upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			allowOrigins := []string{"http://localhost:4321"} // TODO: mapping prd domain
			header := r.Header.Get("Origin")

			for _, o := range allowOrigins {
				if header == o {
					return true
				}
			}

			return false
		},
	}
	ReviewChannel  = make(chan string, 100)
	ActiveChannels sync.Map
)
