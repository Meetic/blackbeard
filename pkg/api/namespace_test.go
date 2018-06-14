package api_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApi_ListNamespaces(t *testing.T) {
	namespaces, err := blackbeard.ListNamespaces()

	assert.Nil(t, err)
	assert.NotNil(t, namespaces)
	assert.Equal(t, "test", namespaces[0].Name)
	assert.Equal(t, true, namespaces[0].Managed)
}
