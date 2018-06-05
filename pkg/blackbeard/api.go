package blackbeard

import (
	"log"
)

type Api interface {
	Inventories() InventoryService
	Namespaces() NamespaceService
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
	namespaces  NamespaceService
	services    ServiceService
}

func NewApi(inventories InventoryRepository, configs ConfigRepository, namespaces NamespaceRepository, pods PodRepository, services ServiceRepository) Api {
	return &api{
		inventories: NewInventoryService(inventories),
		configs:     NewConfigService(configs),
		namespaces:  NewNamespaceService(namespaces, pods),
		services:    NewServiceService(services),
	}
}

func (api *api) Inventories() InventoryService {
	return api.inventories
}

func (api *api) Namespaces() NamespaceService {
	return api.namespaces
}

func (api *api) Create(namespace string) (Inventory, error) {
	inv, err := api.inventories.Create(namespace)
	if err != nil {
		return Inventory{}, err
	}

	if err != nil {
		switch e := err.(type) {
		default:
			return Inventory{}, e
		case *ErrorReadingDefaultsFile:
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

func (api *api) GetExposedServices(namespace string) ([]Service, error) {
	return api.services.ListExposed(namespace)
}

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

func (api *api) Update(namespace string, inventory Inventory, configPath string) error {
	if err := api.inventories.Update(namespace, inventory); err != nil {
		return err
	}

	if err := api.Apply(namespace, configPath); err != nil {
		return err
	}

	return nil
}
