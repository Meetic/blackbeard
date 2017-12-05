package blackbeard

import "encoding/json"

//Inventory represents the inventory file.
type Inventory struct {
	Namespace string                 `json:"namespace"`
	Values    map[string]interface{} `json:"values"`
}

//NamespaceService define the way kubernetes namespace should be managed.
type NamespaceService interface {
	Create(Inventory) error
	Apply(Inventory) error
}

//InventoryService define the way inventory should be managed.
type InventoryService interface {
	Create(namespace string) (Inventory, error)
	Update(namespace string, inv Inventory) error
	Get(namespace string) (Inventory, error)
	GetDefaults() (Inventory, error)
	List() ([]Inventory, error)
}

//ConfigService define the way configuration should be managed
type ConfigService interface {
	Apply(inv Inventory) error
}

//ConfigClient is an interface that must be implemented by any kind of blackbeard client.
//ConfigClient could be an http client, a tcp client or event a io.writer that log informations.
type ConfigClient interface {
	InventoryService() InventoryService
	ConfigService() ConfigService
}

//KubeClient is an interface that represents the way kubernetes is managed.
type KubeClient interface {
	NamespaceService() NamespaceService
}

//NewInventory create a new inventory for a given usr.
func NewInventory(namespace string, defaults []byte) Inventory {

	var inventory Inventory

	if err := json.Unmarshal(defaults, &inventory); err != nil {
		panic(err)
	}

	inventory.Namespace = namespace

	return inventory
}
