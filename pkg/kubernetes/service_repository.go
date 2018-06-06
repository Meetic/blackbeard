package kubernetes

import (
	"github.com/Meetic/blackbeard/pkg/blackbeard"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type serviceRepository struct {
	kubernetes kubernetes.Interface
	host       string
}

// NewServiceRepository retuns a new ServiceRespository
// It takes as parameter a go-client kubernetes client and the kubernetes cluster host (domain name or ip).
func NewServiceRepository(kubernetes kubernetes.Interface, host string) blackbeard.ServiceRepository {
	return &serviceRepository{
		kubernetes: kubernetes,
		host:       host,
	}
}

// ListNodePort returns a list of kubernetes services exposed as NodePort.
func (sr *serviceRepository) ListNodePort(namespace string) ([]blackbeard.Service, error) {
	svcs, err := sr.kubernetes.CoreV1().Services(namespace).List(metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	var services []blackbeard.Service

	for _, svc := range svcs.Items {
		if isNodePort(svc) {

			var ports []blackbeard.Port

			for _, p := range svc.Spec.Ports {
				ports = append(ports, blackbeard.Port{
					Port:        p.Port,
					ExposedPort: p.NodePort,
				})
			}

			services = append(services, blackbeard.Service{
				Name:  svc.Name,
				Ports: ports,
				Addr:  sr.host,
			})
		}
	}

	return services, nil

}

// ListIngress returns a list of Kubernetes services exposed throw Ingress.
func (sr *serviceRepository) ListIngress(namespace string) ([]blackbeard.Service, error) {
	ingressList, err := sr.kubernetes.ExtensionsV1beta1().Ingresses(namespace).List(metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	var services []blackbeard.Service

	for _, ing := range ingressList.Items {

		for _, rules := range ing.Spec.Rules {
			for _, path := range rules.HTTP.Paths {
				svc := blackbeard.Service{
					Name: path.Backend.ServiceName,
					Addr: rules.Host,
					Ports: []blackbeard.Port{
						{
							Port:        path.Backend.ServicePort.IntVal,
							ExposedPort: 80,
						},
					},
				}
				services = append(services, svc)
			}
		}
	}

	return services, nil

}

func isNodePort(svc v1.Service) bool {
	var nP int
	for _, p := range svc.Spec.Ports {

		if p.NodePort != 0 {
			nP++
		}
	}

	if nP > 0 {
		return true
	}

	return false

}
