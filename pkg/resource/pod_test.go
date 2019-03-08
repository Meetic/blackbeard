package resource_test

import (
	"testing"

	"github.com/Meetic/blackbeard/pkg/mock"
	"github.com/Meetic/blackbeard/pkg/resource"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
)

var (
	pods = resource.NewPodService(mock.NewPodRepository(kube))
)

func TestFindOk(t *testing.T) {
	p, err := pods.Find("foo", "te")

	assert.Contains(t, p.Name, "te")
	assert.Nil(t, err)
}

func TestFindNOk(t *testing.T) {
	_, err := pods.Find("foo", "bar")

	assert.NotNil(t, err)
}

func TestDeleteNOk(t *testing.T) {
	err := pods.Delete("foo", resource.Pod{"err", v1.PodRunning})
	assert.NotNil(t, err)
}
