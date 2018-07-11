package mock

import (
	"github.com/Meetic/blackbeard/pkg/resource"
	"k8s.io/api/core/v1"
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

	if n == "testko" {

		pods := resource.Pods{
			{
				Name:   "test",
				Status: v1.PodPending,
			},
		}

		return pods, nil

	}

	pods := resource.Pods{
		{
			Name:   "test",
			Status: v1.PodRunning,
		},
	}

	return pods, nil
}
