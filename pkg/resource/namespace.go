package resource

import "k8s.io/api/core/v1"

type Namespace struct {
	Name   string
	Phase  string
	Status int
}

// NamespaceService defined the way namespace are managed.
type NamespaceService interface {
	Create(namespace string) error
	ApplyConfig(namespace string, configPath string) error
	Delete(namespace string) error
	GetStatus(namespace string) (int, error)
	List() ([]Namespace, error)
}

// NamespaceService defined the way namespace area actually managed.
type NamespaceRepository interface {
	Create(namespace string) error
	ApplyConfig(namespace string, configPath string) error
	Delete(namespace string) error
	List() ([]Namespace, error)
}

type namespaceService struct {
	namespaces NamespaceRepository
	pods       PodRepository
}

// NewNamespaceService creates a new NamespaceService
func NewNamespaceService(namespaces NamespaceRepository, pods PodRepository) NamespaceService {
	return &namespaceService{
		namespaces: namespaces,
		pods:       pods,
	}
}

// Create creates a kubernetes namespace
func (ns *namespaceService) Create(n string) error {
	return ns.namespaces.Create(n)
}

// ApplyConfig apply kubernetes configurations to the given namespace.
// Warning : For now, this method takes a configPath as parameter. This parameter is the directory containing configs in a playbook
// This may change since the NamespaceService should not be aware that configs are stored in files.
func (ns *namespaceService) ApplyConfig(namespace, configPath string) error {
	return ns.namespaces.ApplyConfig(namespace, configPath)
}

// Delete deletes a kubernetes namespace
func (ns *namespaceService) Delete(namespace string) error {
	return ns.namespaces.Delete(namespace)
}

// List returns a slice of Namespace from the kubernetes package and enrich each of the
// returned namespace with its status.
func (ns *namespaceService) List() ([]Namespace, error) {
	namespaces, err := ns.namespaces.List()
	if err != nil {
		return nil, err
	}

	for i, namespace := range namespaces {
		status, err := ns.GetStatus(namespace.Name)
		if err != nil {
			return nil, err
		}

		namespaces[i].Status = status
	}

	return namespaces, nil

}

// GetStatus returns the status of an inventory
// The status is an int that represents the percentage of pods in a "running" state inside the given namespace
func (ns *namespaceService) GetStatus(namespace string) (int, error) {

	pods, err := ns.pods.List(namespace)
	if err != nil {
		return 0, err
	}

	totalPods := len(pods)

	if totalPods == 0 {
		return 0, nil
	}

	var i int

	for _, pod := range pods {
		if pod.Status == v1.PodRunning {
			i++
		}
	}

	status := i * 100 / totalPods

	return status, nil

}
