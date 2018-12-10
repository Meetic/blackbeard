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
	GetStatus(namespace string) (*NamespaceStatus, error)
	List() ([]Namespace, error)
}

// NamespaceRepository defined the way namespace area actually managed.
type NamespaceRepository interface {
	Create(namespace string) error
	Get(namespace string) (*Namespace, error)
	ApplyConfig(namespace string, configPath string) error
	Delete(namespace string) error
	List() ([]Namespace, error)
}

type namespaceService struct {
	namespaces NamespaceRepository
	pods       PodRepository
}

type NamespaceStatus struct {
	Status int
	Phase  string
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
	err := ns.namespaces.Create(n)

	if err != nil {
		return ErrorCreateNamespace{err.Error()}
	}

	return nil
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

		namespaces[i].Status = status.Status
	}

	return namespaces, nil
}

// GetStatus returns the status of an inventory
// The status is an int that represents the percentage of pods in a "running" state inside the given namespace
func (ns *namespaceService) GetStatus(namespace string) (*NamespaceStatus, error) {

	// get namespace state
	n, err := ns.namespaces.Get(namespace)
	if err != nil {
		return &NamespaceStatus{0, ""}, err
	}

	//  get pod's namespace
	pods, err := ns.pods.List(namespace)
	if err != nil {
		return &NamespaceStatus{0, ""}, err
	}

	totalPods := len(pods)

	if totalPods == 0 {
		return &NamespaceStatus{0, ""}, nil
	}

	var i int

	for _, pod := range pods {
		if pod.Status == v1.PodRunning {
			i++
		}
	}

	status := i * 100 / totalPods

	return &NamespaceStatus{status, n.Phase}, nil
}

// ErrorCreateNamespace represents an error due to a namespace creation failure on kubernetes cluster
type ErrorCreateNamespace struct {
	Msg string
}

// Error returns the error message
func (err ErrorCreateNamespace) Error() string {
	return err.Msg
}
