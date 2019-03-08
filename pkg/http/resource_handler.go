package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Meetic/blackbeard/pkg/playbook"
)

// ListServices returns the list of exposed services (NodePort and ingress configuration) of a given inventory
func (h *Handler) ListServices(c *gin.Context) {

	services, err := h.api.ListExposedServices(c.Params.ByName("namespace"))

	if err != nil {
		if notFound, ok := err.(playbook.ErrorInventoryNotFound); ok {
			c.JSON(http.StatusNotFound, gin.H{"error": notFound.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, services)
}

// GetStatus returns the namespace status (ready or not) for a given namespace
func (h *Handler) GetStatus(c *gin.Context) {

	_, err := h.api.Inventories().Get(c.Params.ByName("namespace"))

	if err != nil {
		if notFound, ok := err.(playbook.ErrorInventoryNotFound); ok {
			c.JSON(http.StatusNotFound, gin.H{"error": notFound.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	status, err := h.api.Namespaces().GetStatus(c.Params.ByName("namespace"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}

// GetStatuses returns an array of namespaces and their associated status
func (h *Handler) GetStatuses(c *gin.Context) {

	invs, _ := h.api.Inventories().List()

	var statuses []struct {
		Namespace string `json:"namespace"`
		Status    int    `json:"status"`
		Phase     string `json:"phase"`
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
			Phase     string `json:"phase"`
		}{
			Namespace: i.Namespace,
			Status:    s.Status,
			Phase:     s.Phase,
		})
	}

	c.JSON(http.StatusOK, statuses)
}

type killQuery []string

// Kill handle the kill deployments
func (h *Handler) Kill(c *gin.Context) {

	var kQuery killQuery

	if err := c.BindJSON(&kQuery); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Kill the corresponding pods
	errs := h.api.Kill(c.Params.ByName("namespace"), kQuery)

	if len(errs) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": fmt.Sprintf("%v", errs)})
		return
	}
	c.Status(http.StatusOK)
}
