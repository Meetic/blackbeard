package mock

import (
	"github.com/Meetic/blackbeard/pkg/resource"
	"k8s.io/client-go/kubernetes"
)

type namespaceRepository struct {
	kubernetes kubernetes.Interface
}

// NewNamespaceRepository returns a new NamespaceRepository.
// The parameter is a go-client Kubernetes client
func NewNamespaceRepository(kubernetes kubernetes.Interface) resource.NamespaceRepository {
	return &namespaceRepository{
		kubernetes: kubernetes,
	}
}

// Create creates a namespace
func (ns *namespaceRepository) Create(namespace string) error {
	return nil
}

// Delete deletes a given namespace
func (ns *namespaceRepository) Delete(namespace string) error {
	return nil
}

// List returns a slice of Namespace.
// Name is the namespace name from Kubernetes.
// Phase is the status phase.
// List returns an error if the namespace list could not be get from Kubernetes cluster.
func (ns *namespaceRepository) List() ([]resource.Namespace, error) {
	namespaces := []resource.Namespace{
		{
			Name:  "test",
			Phase: "active",
		},
	}

	return namespaces, nil
}

// ApplyConfig loads configuration files into kubernetes
func (ns *namespaceRepository) ApplyConfig(namespace, configPath string) error {
	return nil
}