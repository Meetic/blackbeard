package http

import (
	"github.com/Meetic/blackbeard/pkg/blackbeard"
	"github.com/Meetic/blackbeard/pkg/ws"
	"github.com/gin-gonic/gin"
)

//Handler actually handle http requests.
//It use a router to map uri to HandlerFunc
type Handler struct {
	config blackbeard.ConfigClient
	kube   blackbeard.KubeClient
	engine *gin.Engine
}

//NewHandler create an Handler using defined routes.
//It takes a client as argument in order to be passe to the handler and be accessible to the HandlerFunc
//Typically in a CRUD API, the client manage connections to a storage system.
func NewHandler(c blackbeard.ConfigClient, k blackbeard.KubeClient) *Handler {
	h := &Handler{
		config: c,
		kube:   k,
	}

	h.engine = gin.Default()

	h.engine.POST("/inventories", h.Create)
	h.engine.GET("/inventories/:namespace", h.Get)
	h.engine.GET("/inventories", h.List)
	h.engine.GET("/defaults", h.GetDefaults)
	h.engine.PUT("/inventories/:namespace", h.Update)
	h.engine.GET("/ws/:namespace", func(c *gin.Context) {
		ws.NewHandler().Handle(c.Writer, c.Request, c.Params.ByName("namespace"))
	})

	return h
}

//Engine returns the defined router for the Handler
func (h *Handler) Engine() *gin.Engine { return h.engine }

//Server represents an http server that handle request
type Server struct {
	handler *Handler
}

//NewServer return an http server with a given handler
func NewServer(h *Handler) *Server {
	return &Server{
		handler: h,
	}
}

//Serve launch the webserver
func (s *Server) Serve() {
	s.handler.Engine().Run(":8080")
}
