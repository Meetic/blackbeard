package http

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/Meetic/blackbeard/pkg/api"
	"github.com/sirupsen/logrus"
)

// Handler actually handle http requests.
// It use a router to map uri to HandlerFunc
type Handler struct {
	api        api.Api
	configPath string

	engine *gin.Engine
}

// NewHandler create a Handler using defined routes.
// It takes a client as argument in order to be pass to the handler and be accessible to the HandlerFunc
// Typically in a CRUD API, the client manage connections to a storage system.
func NewHandler(api api.Api, configPath string, corsEnable bool) *Handler {
	h := &Handler{
		api:        api,
		configPath: configPath,
	}

	h.engine = gin.New()
	h.engine.Use(jsonLogMiddleware(), gin.Recovery())

	if corsEnable == true {
		config := cors.DefaultConfig()
		config.AllowAllOrigins = true
		config.AddAllowHeaders("authorization")
		h.engine.Use(cors.New(config))
		logrus.Info("CORS are enabled")
	}

	h.engine.GET("/ready", h.HealthCheck)
	h.engine.GET("/alive", h.HealthCheck)
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
	h.engine.DELETE("/resources/:namespace/jobs/:resource", h.DeleteResource)
	h.engine.GET("/version", h.Version)

	return h
}

// Engine returns the defined router for the Handler
func (h *Handler) Engine() *gin.Engine { return h.engine }

// Server represents a http server that handle request
type Server struct {
	handler *Handler
}

// NewServer return a http server with a given handler
func NewServer(h *Handler) *Server {
	return &Server{
		handler: h,
	}
}

// Serve launch the webserver
func (s *Server) Serve(port int) {
	s.handler.Engine().Run(fmt.Sprintf(":%d", port))
}
