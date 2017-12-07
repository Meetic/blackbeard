package kubernetes

import (
	"github.com/Meetic/blackbeard/pkg/blackbeard"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	//Namespace status
	ready    = "ready"
	notReady = "not ready"
)

type ResourceService struct {
	client *kubernetes.Clientset
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
