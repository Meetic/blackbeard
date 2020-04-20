package mock

import (
	"k8s.io/client-go/kubernetes"

	"github.com/Meetic/blackbeard/pkg/resource"
)

type statefulsetRepository struct {
	kubernetes.Interface
}

func NewStatefulsetRepository(kubernetes kubernetes.Interface) resource.StatefulsetRepository {
	return &statefulsetRepository{kubernetes}
}

func (statefulsetRepository) List(namespace string) (resource.Statefulsets, error) {
	if namespace == "testko" {
		sfs := resource.Statefulsets{
			{
				Name:   "test",
				Status: resource.StatefulsetNotReady,
			},
		}

		return sfs, nil
	}

	sfs := resource.Statefulsets{
		{
			Name:   "test",
			Status: resource.StatefulsetReady,
		},
	}

	return sfs, nil
}
