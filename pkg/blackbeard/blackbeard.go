package blackbeard

import (
	"encoding/json"
	"net/http"
)

//Inventory represents a group of variable to use in the templates.
// Namespace is the namespace dedicated files where to apply the variables contains into Values
// Values is map of string can contains whatever the user set in the default.json file inside a playbook
type Inventory struct {
	Namespace string                 `json:"namespace"`
	Values    map[string]interface{} `json:"values"`
}

//NamespaceService define the way kubernetes namespace should be managed.
type NamespaceService interface {
	Create(Inventory) error
	Apply(Inventory) error
}

//ResourceService defines the way kubernetes resources such as pods, services, etc. should be managed.
type ResourceService interface {
	GetPods(string) (Pods, error)
	GetNamespaceStatus(string) (string, error)
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

//KubectlClient is an interface that represents the way kubernetes is managed using kubectl.
type KubectlClient interface {
	NamespaceService() NamespaceService
}

//KubernetesClient creates a client that use the kubernetes-go-client.
type KubernetesClient interface {
	ResourceService() ResourceService
}

//WebsocketHandler defines the way Websocket should be handled
type WebsocketHandler interface {
	Handle(http.ResponseWriter, *http.Request, string)
}

//NewInventory create a new inventory for a given namespace.
func NewInventory(namespace string, defaults []byte) Inventory {

	var inventory Inventory

	if err := json.Unmarshal(defaults, &inventory); err != nil {
		panic(err)
	}

	inventory.Namespace = namespace

	return inventory
}

//Pods represent a list of pods.
type Pods []Pod

//Pod represent a Kubernetes pod.
type Pod struct {
	Name   string
	Status string
}
