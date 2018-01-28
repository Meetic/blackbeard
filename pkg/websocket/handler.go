package websocket

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Meetic/blackbeard/pkg/blackbeard"
	"github.com/gin-gonic/gin/json"
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

//Handler represent a websocket handler
type Handler struct {
	upgrader   websocket.Upgrader
	kubernetes blackbeard.KubernetesClient
	files      blackbeard.ConfigClient
	conn       *websocket.Conn
}

type namespaceStatus struct {
	Namespace string
	Status    string
}

//NewHandler creates a websocket server
func NewHandler(client blackbeard.KubernetesClient, files blackbeard.ConfigClient) *Handler {
	up := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	h := Handler{
		upgrader:   up,
		kubernetes: client,
		files:      files,
	}

	return &h

}

//Handle upgrade user request to websocket and start a connexion
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade: ", err)
		return
	}

	h.conn = conn

	go h.writer()
	h.reader()
}

func (h *Handler) reader() {
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

func (h *Handler) writer() {
	lastError := ""
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
					status = []byte(lastError)
				}
			} else {
				lastError = ""
			}

			if status != nil {
				h.conn.SetWriteDeadline(time.Now().Add(writeWait))
				if err := h.conn.WriteMessage(websocket.TextMessage, status); err != nil {
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

func (h *Handler) readNamespacesStatus() ([]byte, error) {

	invs, _ := h.files.InventoryService().List()

	var status []namespaceStatus

	for _, i := range invs {
		s, err := h.kubernetes.ResourceService().GetNamespaceStatus(i.Namespace)
		if err != nil {
			return nil, err
		}

		status = append(status, namespaceStatus{
			Namespace: i.Namespace,
			Status:    s,
		})
	}

	return json.Marshal(status)
}
