package kubernetes

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/Meetic/blackbeard/pkg/resource"
)

type statefulsetRepository struct {
	kubernetes.Interface
}

func NewStatefulsetRepository(kubernetes kubernetes.Interface) resource.StatefulsetRepository {
	return &statefulsetRepository{
		kubernetes,
	}
}

// List return a list of statefulset with their status Ready or NotReady
func (r *statefulsetRepository) List(namespace string) (resource.Statefulsets, error) {
	sfl, err := r.AppsV1().StatefulSets(namespace).List(v1.ListOptions{})

	if err != nil {
		return nil, fmt.Errorf("unable to list statefulsets: %v", err)
	}

	sfs := make(resource.Statefulsets, 0)

	for _, dp := range sfl.Items {
		status := resource.StatefulsetNotReady

		if dp.Status.ReadyReplicas == dp.Status.Replicas {
			status = resource.StatefulsetReady
		}

		sfs = append(sfs, resource.Statefulset{
			Name:   dp.Name,
			Status: status,
		})
	}

	return sfs, nil
}
