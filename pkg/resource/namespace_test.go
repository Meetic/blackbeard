package resource_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/Meetic/blackbeard/pkg/mock"
	"github.com/Meetic/blackbeard/pkg/resource"
)

var (
	kube       = fake.NewSimpleClientset()
	namespaces = resource.NewNamespaceService(
		mock.NewNamespaceRepository(kube, false),
		mock.NewPodRepository(kube),
		mock.NewDeploymentRepository(kube),
		mock.NewStatefulsetRepository(kube),
	)
)

func TestGetStatusOk(t *testing.T) {
	status, err := namespaces.GetStatus("test")

	assert.Nil(t, err)
	assert.Equal(t, 100, status.Status)
}

func TestGetStatusUncomplete(t *testing.T) {
	status, err := namespaces.GetStatus("testko")

	assert.Nil(t, err)
	assert.Equal(t, 0, status.Status)
}

func TestNamespaceCreate(t *testing.T) {
	err := namespaces.Create("mynamespace")

	assert.Nil(t, err)
}

func TestNamespaceCreateError(t *testing.T) {
	namespaces = resource.NewNamespaceService(
		mock.NewNamespaceRepository(kube, true),
		mock.NewPodRepository(kube),
		mock.NewDeploymentRepository(kube),
		mock.NewStatefulsetRepository(kube),
	)

	err := namespaces.Create("foobar")

	assert.Equal(t, resource.ErrorCreateNamespace{Msg: "namespace foobar already exist"}, err)
}

func TestAddListener(t *testing.T) {
	namespaces.AddListener("foobar")
	ch := namespaces.Events("foobar")

	assert.NotNil(t, ch)
}

func TestEmit(t *testing.T) {
	namespaces.AddListener("foobar")
	event := resource.NamespaceEvent{Type: "ADDED", Namespace: "namespace"}

	namespaces.Emit(event)
	ch := namespaces.Events("foobar")

	assert.Equal(t, event, <-ch)
}

func TestRemoveListener(t *testing.T) {
	namespaces.AddListener("foobar")
	ch := namespaces.Events("foobar")

	assert.NotNil(t, ch)

	namespaces.RemoveListener("foobar")
	chNil := namespaces.Events("foobar")

	assert.Nil(t, chNil)
}

func TestDelete(t *testing.T) {
	err := namespaces.Delete("foobar")

	assert.Nil(t, err)
}

func TestApplyConfig(t *testing.T) {
	err := namespaces.ApplyConfig("foobar", "config")

	assert.Nil(t, err)
}

func TestList(t *testing.T) {
	namespaces, err := namespaces.List()

	expectedNamespaces := []resource.Namespace{
		{
			Name:   "test",
			Phase:  "Active",
			Status: 100,
		},
	}

	assert.Nil(t, err)
	assert.Equal(t, expectedNamespaces, namespaces)
}
