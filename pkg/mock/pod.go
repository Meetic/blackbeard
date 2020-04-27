package mock

import (
	"github.com/stretchr/testify/mock"

	"github.com/Meetic/blackbeard/pkg/resource"
)

type PodRepository struct {
	mock.Mock
}

// GetPods of all the pods in a given namespace.
// This method returns a Pods slice containing the pod name and the pod status (pod status phase).
func (m* PodRepository) List(n string) (resource.Pods, error) {
	args := m.Called(n)
	return args.Get(0).(resource.Pods), args.Error(1)
}
