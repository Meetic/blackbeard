package http

import (
	"log"

	"github.com/Meetic/blackbeard/pkg/blackbeard"
	"github.com/Meetic/blackbeard/pkg/websocket"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

//Handler actually handle http requests.
//It use a router to map uri to HandlerFunc
type Handler struct {
	api        blackbeard.Api
	websocket  websocket.Handler
	configPath string

	engine *gin.Engine
}

//NewHandler create an Handler using defined routes.
//It takes a client as argument in order to be passe to the handler and be accessible to the HandlerFunc
//Typically in a CRUD API, the client manage connections to a storage system.
func NewHandler(api blackbeard.Api, websocket websocket.Handler, configPath string, corsEnable bool) *Handler {
	h := &Handler{
		api:        api,
		websocket:  websocket,
		configPath: configPath,
	}

	h.engine = gin.Default()

	if corsEnable == true {
		config := cors.DefaultConfig()
		config.AllowAllOrigins = true
		config.AddAllowHeaders("authorization")
		h.engine.Use(cors.New(config))
		log.Println("cors are enabled")
	}

	h.engine.POST("/inventories", h.Create)
	h.engine.GET("/inventories/:namespace", h.Get)
	h.engine.GET("/inventories/:namespace/status", h.GetStatus)
	h.engine.POST("/inventories/:namespace/reset", h.Reset)
	h.engine.GET("/inventories/:namespace/services", h.ListServices)
	h.engine.GET("/inventories", h.List)
	//h.engine.GET("/inventories/status", h.GetStatuses)
	h.engine.GET("/defaults", h.GetDefaults)
	h.engine.PUT("/inventories/:namespace", h.Update)
	h.engine.DELETE("/inventories/:namespace", h.Delete)
	h.engine.GET("/ws", func(c *gin.Context) {
		websocket.Handle(c.Writer, c.Request)
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
