package api_test

import (
	"github.com/Meetic/blackbeard/pkg/api"
	"github.com/Meetic/blackbeard/pkg/kubernetes"
	"github.com/Meetic/blackbeard/pkg/mock"
	"k8s.io/client-go/kubernetes/fake"
)

var (
	kube = fake.NewSimpleClientset()

	blackbeard = api.NewApi(
		mock.NewInventoryRepository(),
		mock.NewConfigRepository(),
		mock.NewPlaybookRepository(),
		mock.NewNamespaceRepository(kube, false),
		kubernetes.NewPodRepository(kube),
		kubernetes.NewServiceRepository(kube, "kube.test"))
)
