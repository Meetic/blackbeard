package http

import (
	"net/http"

	"github.com/Meetic/blackbeard/pkg/files"
	"github.com/gin-gonic/gin"
)

//ListServices returns the list of exposed services (NodePort and ingress configuration) of a given inventory
func (h *Handler) ListServices(c *gin.Context) {

	_, err := h.config.InventoryService().Get(c.Params.ByName("namespace"))

	if err != nil {
		if notFound, ok := err.(files.ErrorInventoryNotFound); ok {
			c.JSON(http.StatusNotFound, gin.H{"error": notFound.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	services, err := h.kubernetes.ResourceService().GetExposedServices(c.Params.ByName("namespace"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, services)
}
