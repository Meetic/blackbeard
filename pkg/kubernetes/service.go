package kubernetes

import (
	"fmt"
	"strings"

	"github.com/Meetic/blackbeard/pkg/resource"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type serviceRepository struct {
	kubernetes kubernetes.Interface
	host       string
}

// NewServiceRepository returns a new ServiceRepository
// It takes as parameter a go-client kubernetes client and the kubernetes cluster host (domain name or ip).
func NewServiceRepository(kubernetes kubernetes.Interface, host string) resource.ServiceRepository {
	return &serviceRepository{
		kubernetes: kubernetes,
		host:       host,
	}
}

// ListExternal returns a list of kubernetes services exposed as NodePort or LoadBalancer.
func (sr *serviceRepository) ListExternal(n string) ([]resource.Service, error) {
	// unfortunately, we cant filter service by type using field selector
	svcs, err := sr.kubernetes.CoreV1().Services(n).List(metav1.ListOptions{})

	if err != nil {
		return nil, fmt.Errorf("kubernetes api list services : %s", err.Error())
	}

	var services []resource.Service

	for _, svc := range svcs.Items {
		if svc.Spec.Type == v1.ServiceTypeNodePort || svc.Spec.Type == v1.ServiceTypeLoadBalancer {
			var ports []resource.Port

			for _, p := range svc.Spec.Ports {
				ports = append(ports, resource.Port{
					Port:        p.Port,
					ExposedPort: p.NodePort,
				})
			}

			addr := sr.host

			if svc.Spec.Type == v1.ServiceTypeLoadBalancer {
				var ips []string
				for _, lbi := range svc.Status.LoadBalancer.Ingress {
					ips = append(ips, lbi.IP)
				}

				addr = strings.Join(ips, ",")
			}

			services = append(services, resource.Service{
				Name:  svc.Name,
				Ports: ports,
				Addr:  addr,
			})

		}
	}

	return services, nil

}

// ListIngress returns a list of Kubernetes services exposed throw Ingress.
func (sr *serviceRepository) ListIngress(n string) ([]resource.Service, error) {
	ingressList, err := sr.kubernetes.ExtensionsV1beta1().Ingresses(n).List(metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	var services []resource.Service

	for _, ing := range ingressList.Items {

		for _, rules := range ing.Spec.Rules {
			for _, path := range rules.HTTP.Paths {
				svc := resource.Service{
					Name: path.Backend.ServiceName,
					Addr: rules.Host,
					Ports: []resource.Port{
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
