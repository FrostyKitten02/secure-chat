package ws

import (
	"github.com/gorilla/websocket"
	"log/slog"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // adjust in prod
}

func Handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("Failed to upgrade websocket", "error", err)
		return
	}
	defer conn.Close()
	conn.WriteMessage(websocket.TextMessage, []byte("Hello, World!"))

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			slog.Error("Read error, closing connection", "error", err)
			break
		}

		slog.Info("Received message", "msg", string(msg))

		// Example: echo back
		if err := conn.WriteMessage(msgType, msg); err != nil {
			slog.Error("Write error, closing connection", "error", err)
			break
		}
	}
}
