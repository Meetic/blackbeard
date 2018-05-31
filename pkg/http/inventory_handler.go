package http

import (
	"net/http"

	"github.com/Meetic/blackbeard/pkg/blackbeard"

	"github.com/gin-gonic/gin"
)

//createQuery represents the POST payload send to the create handler
type createQuery struct {
	Namespace string `json:"namespace" binding:"required"`
}

//Create handle the namespace creation.
func (h *Handler) Create(c *gin.Context) {

	var createQ createQuery

	if err := c.BindJSON(&createQ); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Create inventory
	inv, err := h.api.Create(createQ.Namespace)
	if err != nil {
		if alreadyExist, ok := err.(blackbeard.ErrorInventoryAlreadyExist); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": alreadyExist.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, inv)
}

//Get return an inventory for a given namespace passed has query parameters.
func (h *Handler) Get(c *gin.Context) {

	inv, err := h.api.Inventories().Get(c.Params.ByName("namespace"))

	if err != nil {
		if notFound, ok := err.(blackbeard.ErrorInventoryNotFound); ok {
			c.JSON(http.StatusNotFound, gin.H{"error": notFound.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, inv)
}

//GetDefaults return default for an inventory
// $ curl -xGET defaults/
func (h *Handler) GetDefaults(c *gin.Context) {

	inv, err := h.api.Inventories().GetDefault()

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, inv)
}

//List returns the list of existing inventories.
func (h *Handler) List(c *gin.Context) {

	invList, err := h.api.Inventories().List()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, invList)
}

//Update will update inventory for a given namespace
func (h *Handler) Update(c *gin.Context) {

	var uQ blackbeard.Inventory

	if err := c.BindJSON(&uQ); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.api.Update(c.Params.ByName("namespace"), uQ, h.configPath); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

//Reset reset a inventory to default and apply changes into kubernetes
func (h *Handler) Reset(c *gin.Context) {

	n := c.Params.ByName("namespace")

	if err := h.api.Reset(n, h.configPath); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)

}

//Delete handle the namespace deletion.
func (h *Handler) Delete(c *gin.Context) {
	namespace := c.Params.ByName("namespace")

	//Delete inventory
	if err := h.api.Delete(namespace); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
