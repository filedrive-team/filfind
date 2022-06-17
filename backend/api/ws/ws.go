package ws

import (
	"github.com/gorilla/websocket"
	logger "github.com/sirupsen/logrus"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// ServeWs handles websocket requests from the peer.
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// TODO: Limit the number of connections by IP
	//if set, ok := hub.clients[uid]; ok {
	//	if set.Size() >= MaxWsConnsOneUser {
	//		w.WriteHeader(http.StatusTooManyRequests)
	//		return
	//	}
	//}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.WithError(err).Error("call Upgrade failed")
		return
	}
	client := &Client{
		hub:        hub,
		conn:       conn,
		send:       make(chan []byte, 1024),
		channelSet: NewChannelSet(),
	}

	hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
