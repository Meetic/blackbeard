package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Meetic/blackbeard/pkg/playbook"
	"github.com/Meetic/blackbeard/pkg/resource"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write the file to the client.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Poll kubernetes for changes with this period.
	kubernetesPeriod = 10 * time.Second
)

// WebsocketHandler defines the way Websocket should be handled
type Handler interface {
	Handle(http.ResponseWriter, *http.Request)
}

// Handler represent a websocket handler
type handler struct {
	upgrader    websocket.Upgrader
	namespaces  resource.NamespaceService
	inventories playbook.InventoryService
	conn        *websocket.Conn
}

type namespaceStatus struct {
	Namespace  string
	Status     int
	PodsStatus resource.Pods
}

// NewHandler creates a websocket server
func NewHandler(namespace resource.NamespaceService, inventories playbook.InventoryService) Handler {
	up := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	h := handler{
		upgrader:    up,
		namespaces:  namespace,
		inventories: inventories,
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
	lastError := ""
	var lastStatus []namespaceStatus
	pingTicker := time.NewTicker(pingPeriod)
	kubeTicker := time.NewTicker(kubernetesPeriod)
	defer func() {
		pingTicker.Stop()
		kubeTicker.Stop()
		h.conn.Close()
	}()

	for {
		select {
		case <-kubeTicker.C:

			status, err := h.readNamespacesStatus()

			if err != nil {
				if s := err.Error(); s != lastError {
					lastError = s
					h.conn.SetWriteDeadline(time.Now().Add(writeWait))
					if err := h.conn.WriteMessage(websocket.TextMessage, []byte(lastError)); err != nil {
						return
					}
				}
			} else {
				lastError = ""
			}

			returnStatus := diff(status, lastStatus)

			lastStatus = status

			if returnStatus != nil {
				h.conn.SetWriteDeadline(time.Now().Add(writeWait))
				jsonStatus, _ := json.Marshal(returnStatus)
				if err := h.conn.WriteMessage(websocket.TextMessage, jsonStatus); err != nil {
					return
				}
			}
		case <-pingTicker.C:
			h.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := h.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func (h *handler) readNamespacesStatus() ([]namespaceStatus, error) {

	invs, _ := h.inventories.List()

	var status []namespaceStatus

	for _, i := range invs {
		s, err := h.namespaces.GetStatus(i.Namespace)
		if err != nil {
			return nil, err
		}

		pods, err := h.namespaces.GetPods(i.Namespace)
		if err != nil {
			return nil, err
		}

		status = append(status, namespaceStatus{
			Namespace:  i.Namespace,
			Status:     s,
			PodsStatus: pods,
		})
	}

	return status, nil
}

func diff(now []namespaceStatus, before []namespaceStatus) []namespaceStatus {
	var diff []namespaceStatus

	for i := 0; i < 2; i++ {
		for _, s1 := range now {
			found := false
			statusDiff := true
			for _, s2 := range before {
				if s1.Namespace == s2.Namespace {
					found = true
					if s1.Status == s2.Status {
						statusDiff = false
					}
					break
				}
			}

			if !found || statusDiff {
				diff = append(diff, s1)
			}
		}

		if i == 0 {
			now, before = before, now
		}
	}

	return diff
}
