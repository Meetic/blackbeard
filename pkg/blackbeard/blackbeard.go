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

//NamespaceConfigurationService apply configuration file to a namespace.
type NamespaceConfigurationService interface {
	Apply(string) error
}

//NamespaceService defined the way namespace should be managed.
type NamespaceService interface {
	Create(string) error
	Delete(string) error
}

//ResourceService defines the way kubernetes resources such as pods, services, etc. should be managed.
type ResourceService interface {
	GetPods(string) (Pods, error)
	GetNamespaceStatus(string) (int, error)
	GetExposedServices(string) ([]Service, error)
}

//InventoryService define the way inventory should be managed.
type InventoryService interface {
	Create(string) (Inventory, error)
	Update(string, Inventory) error
	Get(string) (Inventory, error)
	GetDefaults() (Inventory, error)
	List() ([]Inventory, error)
	Delete(string) error
	Reset(string) error
}

//ConfigService define the way configuration should be managed
type ConfigService interface {
	Apply(inv Inventory) error
	Delete(string) error
}

//ConfigClient is an interface that must be implemented by any kind of blackbeard client.
//ConfigClient could be an http client, a tcp client or event a io.writer that log informations.
type ConfigClient interface {
	InventoryService() InventoryService
	ConfigService() ConfigService
}

//KubectlClient is an interface that represents the way kubernetes is managed using kubectl.
type KubectlClient interface {
	NamespaceConfigurationService() NamespaceConfigurationService
}

//KubernetesClient creates a client that use the kubernetes-go-client.
type KubernetesClient interface {
	ResourceService() ResourceService
	NamespaceService() NamespaceService
}

//WebsocketHandler defines the way Websocket should be handled
type WebsocketHandler interface {
	Handle(http.ResponseWriter, *http.Request)
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

//Port represent a kubernetes service port.
//This mean an internal port and a exposed port
type Port struct {
	Port        int32 `json:"port"`
	ExposedPort int32 `json:"exposedPort"`
}

//Service represent a kubernetes service
type Service struct {
	Name  string `json:"name"`
	Ports []Port `json:"ports"`
	Addr  string `json:"addr"`
}
