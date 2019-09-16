package api

import (
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/Meetic/blackbeard/pkg/playbook"
	"github.com/Meetic/blackbeard/pkg/resource"
	"github.com/Meetic/blackbeard/pkg/version"
)

// Api represents the blackbeard entrypoint by defining the list of actions
// blackbeard is able to perform.
type Api interface {
	Inventories() playbook.InventoryService
	Namespaces() resource.NamespaceService
	Playbooks() playbook.PlaybookService
	Pods() resource.PodService
	Create(namespace string) (playbook.Inventory, error)
	Delete(namespace string, wait bool) error
	ListExposedServices(namespace string) ([]resource.Service, error)
	ListNamespaces() ([]Namespace, error)
	Reset(namespace string, configPath string) error
	Apply(namespace string, configPath string) error
	Update(namespace string, inventory playbook.Inventory, configPath string) error
	WaitForNamespaceReady(namespace string, timeout time.Duration, bar progress) error
	GetVersion() (*Version, error)
}

type api struct {
	inventories playbook.InventoryService
	configs     playbook.ConfigService
	playbooks   playbook.PlaybookService
	namespaces  resource.NamespaceService
	pods        resource.PodService
	services    resource.ServiceService
	cluster     resource.ClusterService
}

// NewApi creates a blackbeard api. The blackbeard api is responsible for managing playbooks and namespaces.
// Parameters are struct implementing respectively Inventory, Config, Namespace, Pod and Service interfaces.
func NewApi(
	inventories playbook.InventoryRepository,
	configs playbook.ConfigRepository,
	playbooks playbook.PlaybookRepository,
	namespaces resource.NamespaceRepository,
	pods resource.PodRepository,
	services resource.ServiceRepository,
	cluster resource.ClusterRepository,
) Api {
	api := &api{
		inventories: playbook.NewInventoryService(inventories, playbook.NewPlaybookService(playbooks)),
		configs:     playbook.NewConfigService(configs, playbook.NewPlaybookService(playbooks)),
		playbooks:   playbook.NewPlaybookService(playbooks),
		namespaces:  resource.NewNamespaceService(namespaces, pods),
		pods:        resource.NewPodService(pods),
		services:    resource.NewServiceService(services),
		cluster:     resource.NewClusterService(cluster),
	}

	go api.WatchDelete()

	return api
}

// Inventories returns the Inventory Service from the api
func (api *api) Inventories() playbook.InventoryService {
	return api.inventories
}

// Namespaces returns the Namespace Service from the api
func (api *api) Namespaces() resource.NamespaceService {
	return api.namespaces
}

// Playbooks returns the Playbook Service from the api
func (api *api) Playbooks() playbook.PlaybookService {
	return api.playbooks
}

func (api *api) Pods() resource.PodService {
	return api.pods
}

// Create is responsible for creating an inventory, a set of kubernetes configs and a kubernetes namespace
// for a given namespace.
// If an inventory already exist, Create will log the error and continue the process. Configs will be override.
func (api *api) Create(namespace string) (playbook.Inventory, error) {
	if err := api.namespaces.Create(namespace); err != nil {
		return playbook.Inventory{}, err
	}

	inv, err := api.inventories.Create(namespace)
	if err != nil {
		switch e := err.(type) {
		default:
			return playbook.Inventory{}, e
		case *playbook.ErrorInventoryAlreadyExist:
			logrus.Warn(e.Error())
			logrus.Info("Process continue")
		}
	}

	if err := api.configs.Generate(inv); err != nil {
		return playbook.Inventory{}, err
	}

	return inv, nil
}

// Delete deletes the inventory, configs and kubernetes namespace for the given namespace.
func (api *api) Delete(namespace string, wait bool) error {
	// delete namespace
	if err := api.namespaces.Delete(namespace); err != nil {
		return err
	}

	if !wait {
		api.deletePlaybook(namespace)
	}

	return nil
}

// ListExposedServices returns a list of services exposed somehow outside of the kubernetes cluster.
// Exposed services could be :
// * NodePort type services
// * Http services exposed throw Ingress
func (api *api) ListExposedServices(namespace string) ([]resource.Service, error) {
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
func (api *api) Update(namespace string, inventory playbook.Inventory, configPath string) error {
	if err := api.inventories.Update(namespace, inventory); err != nil {
		return err
	}

	if err := api.Apply(namespace, configPath); err != nil {
		return err
	}

	return nil
}

func (api *api) deletePlaybook(namespace string) {
	if inv, _ := api.inventories.Get(namespace); inv.Namespace == namespace {
		api.inventories.Delete(namespace)
		api.configs.Delete(namespace)
	}
}

func (api *api) WatchDelete() {
	api.namespaces.AddListener("http")

	// handle delete of inventories and configs files
	for event := range api.namespaces.Events("http") {
		if event.Type == "DELETED" {
			api.deletePlaybook(event.Namespace)
			logrus.WithFields(logrus.Fields{
				"component": "watcher",
				"event":     "delete",
				"namespace": event.Namespace,
			}).Debug("Playbook deleted")
		}
	}
}

type Version struct {
	Blackbeard string `json:"blackbeard"`
	Kubernetes string `json:"kubernetes"`
	Kubectl    string `json:"kubectl"`
}

func (api *api) GetVersion() (*Version, error) {
	v, err := api.cluster.GetVersion()

	if err != nil {
		return nil, err
	}

	return &Version{
		Blackbeard: version.GetVersion(),
		Kubernetes: strings.Join([]string{v.ClientVersion.Major, v.ClientVersion.Minor}, "."),
		Kubectl:    strings.Join([]string{v.ServerVersion.Major, v.ServerVersion.Minor}, "."),
	}, nil
}
