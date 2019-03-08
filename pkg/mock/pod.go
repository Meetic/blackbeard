package mock

import (
	"errors"

	"github.com/Meetic/blackbeard/pkg/resource"
	v1 "k8s.io/api/core/v1"
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

func (rs *podRepository) Delete(n string, p resource.Pod) error {
	if p.Name == "err" {
		return errors.New("an error occurred in pod deletion")
	}
	return nil
}
