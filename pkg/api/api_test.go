package api_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/Meetic/blackbeard/pkg/api"
	"github.com/Meetic/blackbeard/pkg/kubernetes"
	"github.com/Meetic/blackbeard/pkg/mock"
	"github.com/Meetic/blackbeard/pkg/resource"
)

var (
	kube = fake.NewSimpleClientset()

	blackbeard = api.NewApi(
		mock.NewInventoryRepository(),
		mock.NewConfigRepository(),
		mock.NewPlaybookRepository(),
		mock.NewNamespaceRepository(kube, false),
		kubernetes.NewPodRepository(kube),
		kubernetes.NewServiceRepository(kube, "kube.test"),
		kubernetes.NewClusterRepository(),
		kubernetes.NewJobRepository(kube),
	)
)

type clusterRepositoryMock struct{}

func (clusterRepositoryMock) GetVersion() (*resource.Version, error) {
	return &resource.Version{
		ServerVersion: struct {
			Major string `json:"major"`
			Minor string `json:"minor"`
		}{Major: "1", Minor: "2"},
		ClientVersion: struct {
			Major string `json:"major"`
			Minor string `json:"minor"`
		}{Major: "0", Minor: "9"},
	}, nil
}

func TestGetVersion(t *testing.T) {
	blackbeard = api.NewApi(
		mock.NewInventoryRepository(),
		mock.NewConfigRepository(),
		mock.NewPlaybookRepository(),
		mock.NewNamespaceRepository(kube, false),
		kubernetes.NewPodRepository(kube),
		kubernetes.NewServiceRepository(kube, "kube.test"),
		&clusterRepositoryMock{},
		kubernetes.NewJobRepository(kube),
	)

	version, err := blackbeard.GetVersion()

	assert.Nil(t, err)
	assert.Equal(t, version, &api.Version{Blackbeard: "dev", Kubernetes: "1.2", Kubectl: "0.9"})
}
