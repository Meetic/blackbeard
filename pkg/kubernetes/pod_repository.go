package kubernetes

import (
	"github.com/Meetic/blackbeard/pkg/blackbeard"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type podRepository struct {
	kubernetes kubernetes.Interface
}

// NewPodRepository returns a new PodRepository.
// The parameter is a go-client kubernetes client.
func NewPodRepository(kubernetes kubernetes.Interface) blackbeard.PodRepository {
	return &podRepository{
		kubernetes: kubernetes,
	}
}

// GetPods of all the pods in a given namespace.
// This method returns a Pods slice containing the pod name and the pod status (pod status phase).
func (rs *podRepository) List(namespace string) (blackbeard.Pods, error) {
	podsList, err := rs.kubernetes.CoreV1().Pods(namespace).List(metav1.ListOptions{})

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
