package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Version(c *gin.Context) {
	version, err := h.api.GetVersion()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Unable to get blackbeard version",
			"message": err.Error(),
		})
	}

	c.JSON(http.StatusOK, version)
}

func (h *Handler) HealthCheck(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusOK)
}
