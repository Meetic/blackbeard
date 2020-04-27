package mock

import (
	"github.com/stretchr/testify/mock"

	"github.com/Meetic/blackbeard/pkg/resource"
)

type StatefulsetRepository struct {
	mock.Mock
}

func (m *StatefulsetRepository) List(namespace string) (resource.Statefulsets, error) {
	args := m.Called(namespace)
	return args.Get(0).(resource.Statefulsets), args.Error(1)
}
