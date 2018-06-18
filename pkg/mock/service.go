package mock

import (
	"github.com/Meetic/blackbeard/pkg/resource"
	"k8s.io/client-go/kubernetes"
)

type serviceRepository struct {
	kubernetes kubernetes.Interface
	host       string
}

// NewServiceRepository retuns a new ServiceRespository
// It takes as parameter a go-client kubernetes client and the kubernetes cluster host (domain name or ip).
func NewServiceRepository(kubernetes kubernetes.Interface, host string) resource.ServiceRepository {
	return &serviceRepository{
		kubernetes: kubernetes,
		host:       host,
	}
}

// ListNodePort returns a list of kubernetes services exposed as NodePort.
func (sr *serviceRepository) ListNodePort(n string) ([]resource.Service, error) {
	services := []resource.Service{
		{
			Name: "testPort",
		},
	}

	return services, nil
}

// ListIngress returns a list of Kubernetes services exposed throw Ingress.
func (sr *serviceRepository) ListIngress(n string) ([]resource.Service, error) {
	services := []resource.Service{
		{
			Name: "testIngress",
		},
	}

	return services, nil
}
