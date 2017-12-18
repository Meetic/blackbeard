package kubernetes

import (
	"github.com/Meetic/blackbeard/pkg/blackbeard"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	//Namespace status
	ready    = "ready"
	notReady = "not ready"
)

type ResourceService struct {
	client kubernetes.Interface
	host   string
}

//Ensure that ResourceService implements the interface
var _ blackbeard.ResourceService = (*ResourceService)(nil)

//GetPods of all the pods in a given namespace.
//This method returns a Pods slice containing the pod name and the pod status (pod status phase).
func (rs *ResourceService) GetPods(namespace string) (blackbeard.Pods, error) {
	podsList, err := rs.client.CoreV1().Pods(namespace).List(metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	var pods blackbeard.Pods

	for _, pod := range podsList.Items {

		pods = append(pods, blackbeard.Pod{
			Name:   pod.ObjectMeta.Name,
			Status: string(pod.Status.Phase),
		})
	}

	return pods, nil
}

//IsNamespaceReady return the current status of a given namespace based on the pods status
//contains in the namespace.
func (rs *ResourceService) GetNamespaceStatus(namespace string) (string, error) {
	pods, err := rs.GetPods(namespace)

	if err != nil {
		return notReady, err
	}

	r := true
	for _, pod := range pods {
		if pod.Status != running {
			r = false
		}
	}

	if !r {
		return notReady, nil
	}

	return ready, nil

}

//GetExposedServices find services exposed as NodePort and ingress configuration and return
//an array of services containing an URL, the exposed port and the service name.
func (rs *ResourceService) GetExposedServices(namespace string) ([]blackbeard.Service, error) {

	var (
		services []blackbeard.Service
		err      error
	)

	services, err = rs.getNodePortServices(namespace)
	if err != nil {
		return nil, err
	}

	ingress, err := rs.getIngress(namespace)
	if err != nil {
		return nil, err
	}

	services = append(services, ingress...)

	return services, nil
}

func (rs *ResourceService) getNodePortServices(namespace string) ([]blackbeard.Service, error) {
	svcs, err := rs.client.CoreV1().Services(namespace).List(metav1.ListOptions{})

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
				Addr:  rs.host,
			})
		}
	}

	return services, nil

}

func (rs *ResourceService) getIngress(namespace string) ([]blackbeard.Service, error) {
	ingressList, err := rs.client.ExtensionsV1beta1().Ingresses(namespace).List(metav1.ListOptions{})

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
