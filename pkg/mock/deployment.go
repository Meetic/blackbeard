package mock

import (
	"github.com/stretchr/testify/mock"

	"github.com/Meetic/blackbeard/pkg/resource"
)

type DeploymentRepository struct {
	mock.Mock
}

func (m *DeploymentRepository) List(namespace string) (resource.Deployments, error) {
	args := m.Called(namespace)
	return args.Get(0).(resource.Deployments), args.Error(1)
}
