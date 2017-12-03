package ws

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/websocket"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	// Time allowed to write the file to the client.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Poll file for changes with this period.
	kubePeriod = 10 * time.Second
)

type Handler struct {
	upgrader websocket.Upgrader
}

func NewHandler() *Handler {
	up := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	h := Handler{
		upgrader: up,
	}

	return &h

}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request, namespace string) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade: ", err)
		return
	}

	go writer(conn, namespace)
	reader(conn)
}

func reader(ws *websocket.Conn) {
	defer ws.Close()
	ws.SetReadLimit(512)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}
}

func writer(ws *websocket.Conn, namespace string) {
	lastError := ""
	pingTicker := time.NewTicker(pingPeriod)
	kubeTicker := time.NewTicker(kubePeriod)
	defer func() {
		pingTicker.Stop()
		kubeTicker.Stop()
		ws.Close()
	}()
	for {
		select {
		case <-kubeTicker.C:

			status, err := readStatus(namespace)

			if err != nil {
				if s := err.Error(); s != lastError {
					lastError = s
					status = []byte(lastError)
				}
			} else {
				lastError = ""
			}

			if status != nil {
				ws.SetWriteDeadline(time.Now().Add(writeWait))
				if err := ws.WriteMessage(websocket.TextMessage, status); err != nil {
					return
				}
			}
		case <-pingTicker.C:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func readStatus(namespace string) ([]byte, error) {
	home := homeDir()
	config, _ := clientcmd.BuildConfigFromFlags("", filepath.Join(home, ".kube", "config"))

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	podsList, err := clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var running int

	for _, pod := range podsList.Items {
		log.Printf("Pod status : %s", pod.Status.Phase)
		if pod.Status.Phase == "Running" {
			running++
		}
	}

	log.Printf("Running count : %d", running)
	log.Printf("Pods count : %d", len(podsList.Items))

	if running != len(podsList.Items) {
		return []byte("creating"), nil
	}

	return []byte("running"), nil
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
