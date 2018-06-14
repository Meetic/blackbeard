package mock

import (
	"github.com/Meetic/blackbeard/pkg/resource"
	"k8s.io/client-go/kubernetes"
)

type podRepository struct {
	kubernetes kubernetes.Interface
}

// NewPodRepository returns a new PodRepository.
// The parameter is a go-client kubernetes client.
func NewPodRepository(kubernetes kubernetes.Interface) resource.PodRepository {
	return &podRepository{
		kubernetes: kubernetes,
	}
}

// GetPods of all the pods in a given namespace.
// This method returns a Pods slice containing the pod name and the pod status (pod status phase).
func (rs *podRepository) List(n string) (resource.Pods, error) {
	pods := resource.Pods{
		{
			Name:   "test",
			Status: "running",
		},
	}

	return pods, nil
}
