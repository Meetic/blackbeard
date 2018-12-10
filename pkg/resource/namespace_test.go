package resource_test

import (
	"testing"

	"github.com/Meetic/blackbeard/pkg/mock"
	"github.com/Meetic/blackbeard/pkg/resource"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/fake"
)

var (
	kube       = fake.NewSimpleClientset()
	namespaces = resource.NewNamespaceService(mock.NewNamespaceRepository(kube, false), mock.NewPodRepository(kube))
)

func TestGetStatusOk(t *testing.T) {
	status, err := namespaces.GetStatus("test")

	assert.Nil(t, err)
	assert.Equal(t, 100, status)
}

func TestGetStatusUncomplete(t *testing.T) {
	status, err := namespaces.GetStatus("testko")

	assert.Nil(t, err)
	assert.Equal(t, 0, status)
}

func TestNamespaceCreate(t *testing.T) {
	err := namespaces.Create("mynamespace")

	assert.Nil(t, err)
}

func TestNamespaceCreateError(t *testing.T) {
	namespaces = resource.NewNamespaceService(mock.NewNamespaceRepository(kube, true), mock.NewPodRepository(kube))

	err := namespaces.Create("foobar")

	assert.Equal(t, resource.ErrorCreateNamespace{Msg: "namespace foobar already exist"}, err)
}
