package resource_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/Meetic/blackbeard/pkg/mock"
	"github.com/Meetic/blackbeard/pkg/resource"
)

var (
	podRepository         = new(mock.PodRepository)
	deploymentRepository  = new(mock.DeploymentRepository)
	statefulsetRepository = new(mock.StatefulsetRepository)
	jobRepository         = new(mock.JobRepository)
)

var (
	kube       = fake.NewSimpleClientset()
	namespaces = resource.NewNamespaceService(
		mock.NewNamespaceRepository(kube, false),
		podRepository,
		deploymentRepository,
		statefulsetRepository,
		jobRepository,
	)
)

func TestGetStatusOk(t *testing.T) {
	deploymentRepository.
		On("List", "test").
		Return(resource.Deployments{{Name: "app", Status: resource.DeploymentReady}}, nil)

	statefulsetRepository.
		On("List", "test").
		Return(resource.Statefulsets{{Name: "app", Status: resource.StatefulsetReady}}, nil)

	jobRepository.
		On("List", "test").
		Return(resource.Jobs{{Name: "app", Status: resource.JobReady}}, nil)

	status, err := namespaces.GetStatus("test")

	deploymentRepository.AssertExpectations(t)
	statefulsetRepository.AssertExpectations(t)
	jobRepository.AssertExpectations(t)

	assert.Nil(t, err)
	assert.Equal(t, 100, status.Status)
}

func TestGetStatusIncomplete(t *testing.T) {
	deploymentRepository.
		On("List", "testko").
		Return(resource.Deployments{{Name: "app", Status: resource.DeploymentReady}}, nil)

	statefulsetRepository.
		On("List", "testko").
		Return(resource.Statefulsets{}, errors.New("some error"))

	jobRepository.
		On("List", "testko").
		Return(resource.Jobs{{Name: "app", Status: resource.JobReady}}, nil)

	status, err := namespaces.GetStatus("testko")

	deploymentRepository.AssertExpectations(t)
	statefulsetRepository.AssertExpectations(t)
	jobRepository.AssertExpectations(t)

	assert.NotNil(t, err)
	assert.Equal(t, 0, status.Status)
}

func TestNamespaceCreate(t *testing.T) {
	err := namespaces.Create("mynamespace")

	assert.Nil(t, err)
}

func TestNamespaceCreateError(t *testing.T) {
	namespaces = resource.NewNamespaceService(
		mock.NewNamespaceRepository(kube, true),
		podRepository,
		deploymentRepository,
		statefulsetRepository,
		jobRepository,
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
