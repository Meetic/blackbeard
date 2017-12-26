package http

import (
	"net/http"

	"github.com/Meetic/blackbeard/pkg/blackbeard"
	"github.com/Meetic/blackbeard/pkg/files"

	"github.com/gin-gonic/gin"
)

//createQuery represents the POST payload send to the create handler
type createQuery struct {
	Namespace string `json:"namespace,required" binding:"required"`
}

// @Summary Create an inventory
// @Description Create an inventory for the given namespace. This will also create the inventory file and the associated namespace.
// @ID create-inventory
// @Accept  json
// @Produce  json
// @Param   namespace     body    http.createQuery     true        "Namespace"
// @Success 200 {object} blackbeard.Inventory	"The inventory"
// @Failure 400 {string} string   "The inventory already exists"
// @Failure 422 {string} string   "The inventory could not be created due to communication with kubernetes"
// @Failure 500 {string} string   "Something went wrong checking for existing inventories"
// @Router /inventories [post]
//Create handle the inventory and associated namespace creation
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

// @Summary Return an inventory
// @Description Read inventory file for a given namespace and return it as a json object.
// @ID get-inventory
// @Accept  json
// @Produce  json
// @Param   namespace     path    string     true        "Namespace"
// @Success 200 {object} blackbeard.Inventory	"The inventory"
// @Failure 404 {string} string   "Can not find the namespace/inventory"
// @Failure 500 {string} string   "Something went wrong when reading the inventory"
// @Router /inventories/{namespace} [get]
//Get return an inventory for a given namespace passed has query parameters.
func (h *Handler) Get(c *gin.Context) {

	inv, err := h.config.InventoryService().Get(c.Params.ByName("namespace"))

	if err != nil {
		if notFound, ok := err.(files.ErrorInventoryNotFound); ok {
			c.JSON(http.StatusNotFound, gin.H{"error": notFound.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, inv)
}

// @Summary Get default value for an inventory
// @Description Return the content of the defaults.json file in the used playbook.
// @ID get-defaults
// @Accept  json
// @Produce  json
// @Success 200 {object} blackbeard.Inventory	"The default inventory"
// @Failure 404 {string} string   "Defaults file not found"
// @Router /defaults [get]
//GetDefaults return default for an inventory
func (h *Handler) GetDefaults(c *gin.Context) {

	inv, err := h.config.InventoryService().GetDefaults()

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, inv)
}

// @Summary Return the list of existing inventories
// @Description Read all inventory files and return them as an array
// @ID list-inventories
// @Accept  json
// @Produce  json
// @Success 200 {object} blackbeard.Inventories "List of inventories"
// @Failure 500 {string} string   "Impossible to read the inventory list"
// @Router /inventories/ [get]
//List returns the list of existing inventories.
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
