package kubernetes

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/Meetic/blackbeard/pkg/resource"
)

type deploymentRepository struct {
	kubernetes.Interface
}

func NewDeploymentRepository(kubernetes kubernetes.Interface) resource.DeploymentRepository {
	return &deploymentRepository{
		kubernetes,
	}
}

// List return a list of deployment with their status Ready or NotReady
func (r *deploymentRepository) List(namespace string) (resource.Deployments, error) {
	dl, err := r.AppsV1().Deployments(namespace).List(v1.ListOptions{})

	if err != nil {
		return nil, fmt.Errorf("unable to list deployments: %v", err)
	}

	dps := make(resource.Deployments, 0)

	for _, dp := range dl.Items {
		status := resource.DeploymentNotReady

		if dp.Status.ReadyReplicas == dp.Status.Replicas {
			status = resource.DeploymentReady
		}

		dps = append(dps, resource.Deployment{
			Name:   dp.Name,
			Status: status,
		})
	}

	return dps, nil
}
