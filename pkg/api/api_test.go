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
		mock.NewPodRepository(kube),
		kubernetes.NewServiceRepository(kube, "kube.test"),
		kubernetes.NewClusterRepository(),
	)
)

type clusterRepositoryMock struct{}

func (clusterRepositoryMock) GetVersion() (*resource.Version, error) {
	return &resource.Version{
		ClientVersion: struct {
			Major string `json:"major"`
			Minor string `json:"minor"`
		}{Major: "1", Minor: "2"},
		ServerVersion: struct {
			Major string `json:"major"`
			Minor string `json:"minor"`
		}{Major: "0", Minor: "9"},
	}, nil
}

func TestGetVersion(t *testing.T) {
	b := api.NewApi(
		mock.NewInventoryRepository(),
		mock.NewConfigRepository(),
		mock.NewPlaybookRepository(),
		mock.NewNamespaceRepository(kube, false),
		kubernetes.NewPodRepository(kube),
		kubernetes.NewServiceRepository(kube, "kube.test"),
		&clusterRepositoryMock{},
	)

	version, err := b.GetVersion()

	assert.Nil(t, err)
	assert.Equal(t, version, &api.Version{Blackbeard: "dev", Kubernetes: "1.2", Kubectl: "0.9"})
}
