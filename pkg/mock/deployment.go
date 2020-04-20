package mock

import (
	"k8s.io/client-go/kubernetes"

	"github.com/Meetic/blackbeard/pkg/resource"
)

type deploymentRepository struct {
	kubernetes.Interface
}

func NewDeploymentRepository(kubernetes kubernetes.Interface) resource.DeploymentRepository {
	return &deploymentRepository{kubernetes}
}

func (deploymentRepository) List(namespace string) (resource.Deployments, error) {
	if namespace == "testko" {
		dps := resource.Deployments{
			{
				Name:   "test",
				Status: resource.DeploymentNotReady,
			},
		}

		return dps, nil
	}

	dps := resource.Deployments{
		{
			Name:   "test",
			Status: resource.DeploymentReady,
		},
	}

	return dps, nil
}
