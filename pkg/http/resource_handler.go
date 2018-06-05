package http

import (
	"net/http"

	"github.com/Meetic/blackbeard/pkg/blackbeard"
	"github.com/gin-gonic/gin"
)

//ListServices returns the list of exposed services (NodePort and ingress configuration) of a given inventory
func (h *Handler) ListServices(c *gin.Context) {

	services, err := h.api.GetExposedServices(c.Params.ByName("namespace"))

	if err != nil {
		if notFound, ok := err.(blackbeard.ErrorInventoryNotFound); ok {
			c.JSON(http.StatusNotFound, gin.H{"error": notFound.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, services)
}

//GetStatus returns the namespace status (ready or not) for a given namespace
func (h *Handler) GetStatus(c *gin.Context) {

	_, err := h.api.Inventories().Get(c.Params.ByName("namespace"))

	if err != nil {
		if notFound, ok := err.(blackbeard.ErrorInventoryNotFound); ok {
			c.JSON(http.StatusNotFound, gin.H{"error": notFound.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	status, err := h.api.Namespaces().GetStatus(c.Params.ByName("namespace"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, struct {
		Status int `json:"status"`
	}{
		Status: status,
	})
}

//GetStatuses returns an array of namespaces and their associated status
func (h *Handler) GetStatuses(c *gin.Context) {

	invs, _ := h.api.Inventories().List()

	var statuses []struct {
		Namespace string `json:"namespace"`
		Status    int    `json:"status"`
	}

	for _, i := range invs {
		s, err := h.api.Namespaces().GetStatus(i.Namespace)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		statuses = append(statuses, struct {
			Namespace string `json:"namespace"`
			Status    int    `json:"status"`
		}{
			Namespace: i.Namespace,
			Status:    s,
		})
	}

	c.JSON(http.StatusOK, statuses)
}
