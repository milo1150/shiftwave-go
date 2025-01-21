package middleware

import (
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	Upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			allowOrigins := []string{"http://localhost:4321"} // TODO: mapping prd domain
			header := r.Header.Get("Origin")

			// Allow requests without the Origin header (e.g., from Postman) if env is development
			env := os.Getenv("APP_ENV")
			if header == "" && env == "development" {
				log.Println("Allowing request from Postman (no Origin header).")
				return true
			}

			for _, origin := range allowOrigins {
				if header == origin {
					return true
				}
			}

			return false
		},
	}
	ReviewChannel  = make(chan string, 100)
	ActiveChannels sync.Map
)
