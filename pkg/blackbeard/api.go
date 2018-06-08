package blackbeard

import (
	"log"
)

// Api represents the blackbeard entrypoint by defining the list of actions
// blackbeard is able to perform.
type Api interface {
	Inventories() InventoryService
	Namespaces() NamespaceService
	Playbooks() PlaybookService
	Create(namespace string) (Inventory, error)
	Delete(namespace string) error
	GetExposedServices(namespace string) ([]Service, error)
	Reset(namespace string, configPath string) error
	Apply(namespace string, configPath string) error
	Update(namespace string, inventory Inventory, configPath string) error
}

type api struct {
	inventories InventoryService
	configs     ConfigService
	playbooks   PlaybookService
	namespaces  NamespaceService
	services    ServiceService
}

// NewApi creates a blackbeard api. The blackbeard api is responsible for managing playbooks and namespaces.
// Parameters are struct implementing respectively Inventory, Config, Namespace, Pod and Service interfaces.
func NewApi(inventories InventoryRepository, configs ConfigRepository, playbooks PlaybookRepository, namespaces NamespaceRepository, pods PodRepository, services ServiceRepository) Api {
	return &api{
		inventories: NewInventoryService(inventories, NewPlaybookService(playbooks)),
		configs:     NewConfigService(configs, NewPlaybookService(playbooks)),
		playbooks:   NewPlaybookService(playbooks),
		namespaces:  NewNamespaceService(namespaces, pods),
		services:    NewServiceService(services),
	}
}

// Inventories returns the Inventory Service from the api
func (api *api) Inventories() InventoryService {
	return api.inventories
}

// Namespaces returns the Namespace Service from the api
func (api *api) Namespaces() NamespaceService {
	return api.namespaces
}

// Playbooks returns the Playbook Service from the api
func (api *api) Playbooks() PlaybookService {
	return api.playbooks
}

// Create is responsible for creating an inventory, a set of kubernetes configs and a kubernetes namespace
// for a given namespace.
// If an inventory already exist, Create will log the error and continue the process. Configs will be override.
func (api *api) Create(namespace string) (Inventory, error) {
	inv, err := api.inventories.Create(namespace)
	if err != nil {
		switch e := err.(type) {
		default:
			return Inventory{}, e
		case *ErrorInventoryAlreadyExist:
			log.Println(e.Error())
			log.Println("Process continue.")
		}
	}

	if err := api.configs.Generate(inv); err != nil {
		return Inventory{}, err
	}

	if err := api.namespaces.Create(namespace); err != nil {
		return Inventory{}, err
	}

	return inv, nil
}

// Delete deletes the inventory, configs and kubernetes namespace for the given namespace.
func (api *api) Delete(namespace string) error {
	if err := api.inventories.Delete(namespace); err != nil {
		return err
	}

	if err := api.configs.Delete(namespace); err != nil {
		return err
	}

	if err := api.namespaces.Delete(namespace); err != nil {
		return err
	}

	return nil
}

// GetExposedServices returns a list of services exposed somehow outside of the kubernetes cluster.
// Exposed services could be :
// * NodePort type services
// * Http services exposed throw Ingress
func (api *api) GetExposedServices(namespace string) ([]Service, error) {
	return api.services.ListExposed(namespace)
}

// Reset resets an inventory, the associated configs and the kubernetes namespaces to default values.
// Defaults values are defines by the InventoryService GetDefault() method.
func (api *api) Reset(namespace string, configPath string) error {
	//Reset inventory file
	inv, err := api.inventories.Reset(namespace)
	if err != nil {
		return err
	}

	//Apply inventory to configuration
	if err := api.configs.Generate(inv); err != nil {
		return err
	}

	//Apply changes to Kubernetes
	if err = api.namespaces.ApplyConfig(namespace, configPath); err != nil {
		return err
	}

	return nil
}

// Apply override configs with new generated configs and apply the new configs to the kubernetes namespace.
// Warning : For now, Apply require a configPath as parameter.
// configPath is the location of configs for each namespace. This will change in the future since high level
// api should not be aware that configs are stored in files.
func (api *api) Apply(namespace string, configPath string) error {
	inv, err := api.inventories.Get(namespace)
	if err != nil {
		return err
	}

	if err := api.configs.Generate(inv); err != nil {
		return err
	}

	if err := api.namespaces.ApplyConfig(inv.Namespace, configPath); err != nil {
		return err
	}

	return nil

}

// Update replace the inventory associated to the given namespace by the one set in parameters
// and apply the changes to configs and kubernetes namespace (using the Apply method)
func (api *api) Update(namespace string, inventory Inventory, configPath string) error {
	if err := api.inventories.Update(namespace, inventory); err != nil {
		return err
	}

	if err := api.Apply(namespace, configPath); err != nil {
		return err
	}

	return nil
}
