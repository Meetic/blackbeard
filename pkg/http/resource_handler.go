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

func (h *Handler) GetStatus(c *gin.Context) {

	_, err := h.config.InventoryService().Get(c.Params.ByName("namespace"))

	if err != nil {
		if notFound, ok := err.(files.ErrorInventoryNotFound); ok {
			c.JSON(http.StatusNotFound, gin.H{"error": notFound.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	status, err := h.kubernetes.ResourceService().GetNamespaceStatus(c.Params.ByName("namespace"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, struct {
		Status string `json:"status"`
	}{
		Status: status,
	})
}

func (h *Handler) GetStatuses(c *gin.Context) {

	invs, _ := h.config.InventoryService().List()

	var status []struct {
		Namespace string `json:"namespace"`
		Status    string `json:"status"`
	}

	for _, i := range invs {
		s, err := h.kubernetes.ResourceService().GetNamespaceStatus(i.Namespace)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		status = append(status, struct {
			Namespace string `json:"namespace"`
			Status    string `json:"status"`
		}{
			Namespace: i.Namespace,
			Status:    s,
		})
	}

	c.JSON(http.StatusOK, status)
}
