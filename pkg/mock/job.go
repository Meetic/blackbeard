package mock

import (
	"github.com/stretchr/testify/mock"

	"github.com/Meetic/blackbeard/pkg/resource"
)

type JobRepository struct {
	mock.Mock
}

func (m *JobRepository) List(namespace string) (resource.Jobs, error) {
	args := m.Called(namespace)
	return args.Get(0).(resource.Jobs), args.Error(1)
}

func (m *JobRepository) Delete(namespace, resourceName string) error {
	args := m.Called(namespace, resourceName)
	return args.Error(0)
}
