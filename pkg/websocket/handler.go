package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/Meetic/blackbeard/pkg/api"
)

const (
	// Time allowed to write the file to the client.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

// Handler represent a websocket handler
type handler struct {
	upgrader websocket.Upgrader
	api      api.Api
	conn     *websocket.Conn
	mutex   sync.Mutex
}

// NewHandler creates a websocket server
func NewHandler(api api.Api) *handler {
	up := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	h := handler{
		upgrader: up,
		api:      api,
	}

	return &h
}

// Handle upgrade user request to websocket and start a connexion
func (h *handler) Handle(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade: ", err)
		return
	}

	h.conn = conn

	go h.writer()
	h.reader()
}

func (h *handler) reader() {
	defer h.conn.Close()
	h.conn.SetReadLimit(512)
	h.conn.SetReadDeadline(time.Now().Add(pongWait))
	h.conn.SetPongHandler(func(string) error { h.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, _, err := h.conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (h *handler) writer() {
	pingTicker := time.NewTicker(pingPeriod)
	defer func() {
		pingTicker.Stop()
		h.conn.Close()
	}()

	for {
		select {
		case e := <-h.api.Namespaces().Events():
			h.conn.SetWriteDeadline(time.Now().Add(writeWait))
			jsonEvent, _ := json.Marshal(e)
			h.send(websocket.TextMessage, jsonEvent)
		case <-pingTicker.C:
			h.conn.SetWriteDeadline(time.Now().Add(writeWait))
			h.send(websocket.PingMessage, []byte{})
		}
	}
}

func (h *handler) send(messageType int, message []byte) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	return h.conn.WriteMessage(messageType, message)
}