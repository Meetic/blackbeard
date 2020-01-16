package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"

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
}

type client struct {
	socket   *websocket.Conn
	mutex    sync.Mutex
	listener string
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
		logrus.Warnf("Failed to set websocket upgrade: %s", err.Error())
		return
	}

	listener := fmt.Sprintf("ws_client_%s", uuid.NewV4())
	h.api.Namespaces().AddListener(listener)

	logrus.WithFields(logrus.Fields{
		"component": "listener",
		"listener":  listener,
	}).Debugf("adding new listener")

	client := &client{
		socket:   conn,
		listener: listener,
	}

	go h.writer(client)
	h.reader(client)
}

func (h *handler) reader(client *client) {
	client.socket.SetReadLimit(512)
	client.socket.SetReadDeadline(time.Now().Add(pongWait))
	client.socket.SetPongHandler(func(string) error { client.socket.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, _, err := client.socket.ReadMessage()

		if err != nil {
			h.close(client)
			return
		}
	}
}

func (h *handler) writer(client *client) {
	pingTicker := time.NewTicker(pingPeriod)

	defer pingTicker.Stop()
	defer h.close(client)

	for {
		select {
		case e := <-h.api.Namespaces().Events(client.listener):
			client.socket.SetWriteDeadline(time.Now().Add(writeWait))
			jsonEvent, _ := json.Marshal(e)
			h.send(client, websocket.TextMessage, jsonEvent)
		case <-pingTicker.C:
			client.socket.SetWriteDeadline(time.Now().Add(writeWait))
			h.send(client, websocket.PingMessage, []byte{})
		}
	}
}

// Prevent concurrent write to websocket connection
func (h *handler) send(client *client, messageType int, message []byte) error {
	client.mutex.Lock()
	defer client.mutex.Unlock()

	return client.socket.WriteMessage(messageType, message)
}

func (h *handler) close(client *client) {
	if err := client.socket.Close(); err != nil {
		logrus.Warnf("Impossible to close socket connection : %s", err.Error())
	}

	if err := h.api.Namespaces().RemoveListener(client.listener); err != nil {
		logrus.Warnf("Impossible to remove listener: %s", err.Error())
	}

	logrus.WithFields(logrus.Fields{
		"component": "listener",
		"listener":  client.listener,
	}).Debugf("removing listener")
}
