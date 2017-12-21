package http

import (
	"net/http"

	"github.com/Meetic/blackbeard/pkg/blackbeard"
	"github.com/Meetic/blackbeard/pkg/files"

	"github.com/gin-gonic/gin"
)

//createQuery represents the POST payload send to the create handler
type createQuery struct {
	Namespace string `json:"namespace" binding:"required"`
}

//Create handle the testing env creation.
//It is called on route POST /inventories/ and returns either the created inventory
//or an error if the namespace could not be created.
//The payload sent must be like :
// {
//		"namespace": "test"
// }
//This function returns a 400 status if an inventory already exist for the namespace
//Else, the function returns a 500 status or a 422 depending on the error
func (h *Handler) Create(c *gin.Context) {

	var createQ createQuery

	if err := c.BindJSON(&createQ); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Create inventory
	inv, err := h.config.InventoryService().Create(createQ.Namespace)

	if err != nil {
		if alreadyExist, ok := err.(files.ErrorInventoryAlreadyExist); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": alreadyExist.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//Generate config files
	if errc := h.config.ConfigService().Apply(inv); errc != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": errc.Error()})
		return
	}

	//Create namespace
	if errc := h.kubectl.NamespaceService().Create(inv); errc != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": errc.Error()})
		return
	}

	c.JSON(http.StatusCreated, inv)
}

//Get return an inventory for a given namespace passed has query parameters.
// $ curl -xGET inventories/:namespace/
//This function returns a 404 status if the corresponding inventory could not be found.
//Else, it returns a complete inventory read from the InventoryService.
func (h *Handler) Get(c *gin.Context) {

	inv, err := h.config.InventoryService().Get(c.Params.ByName("namespace"))

	if err != nil {
		if notFound, ok := err.(files.ErrorInventoryNotFound); ok {
			c.JSON(http.StatusNotFound, gin.H{"error": notFound.Error()})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, inv)
}

//GetDefaults return default for an inventory
// $ curl -xGET defaults/
func (h *Handler) GetDefaults(c *gin.Context) {

	inv, err := h.config.InventoryService().GetDefaults()

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, inv)
}

//List returns the list of existing inventories.
// Example :
// [
//	{
//		...
// 	},
//	{
// 		...
// 	},
//]
//
//
func (h *Handler) List(c *gin.Context) {

	invList, err := h.config.InventoryService().List()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, invList)
}

//Update will update inventory for a given namespace
//Route: /inventories/{namespace}
//Example  :
// {
// 	"inventory":{
// 				    "namespace": "seblegall",
// 				    "containers": {
// 				        "Microservices": null,
// 				        "PublicAPI": [
// 				            {
// 				                "name": "api-exposure-layer",
// 				                "version": "test",
// 				                "urls": [
// 				                    "authent.ilius.net",
// 				                    "apixl.ilius.net"
// 				                ]
// 				            }
// 				        ]
// 				    }
// 				}
// }
// If the namespace value of the inventory is different from the namespace passed as uri parameters
// The function will rename the corresponding inventory file to match the new namespace
func (h *Handler) Update(c *gin.Context) {

	var uQ blackbeard.Inventory

	if err := c.BindJSON(&uQ); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.config.InventoryService().Update(c.Params.ByName("namespace"), uQ); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	if err := h.config.ConfigService().Apply(uQ); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	if err := h.kubectl.NamespaceService().Apply(uQ); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
