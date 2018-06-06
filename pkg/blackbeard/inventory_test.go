package blackbeard_test

import (
	"testing"

	"github.com/Meetic/blackbeard/pkg/blackbeard"
	"github.com/Meetic/blackbeard/pkg/mock"
	"github.com/stretchr/testify/assert"
)

var (
	inventories = blackbeard.NewInventoryService(mock.NewInventoryRepository())
)

func TestCreateOK(t *testing.T) {
	inv, err := inventories.Create("test1")

	assert.Equal(t, inv.Namespace, "test1")
	assert.Nil(t, err)
}

func TestCreateEmptyNamespace(t *testing.T) {
	_, err := inventories.Create("")

	assert.Error(t, err)
}

func TestGetOK(t *testing.T) {
	inv, _ := inventories.Get("test")
	assert.Equal(t, inv.Namespace, "test")
}

func TestGetEmptyNamespace(t *testing.T) {
	_, err := inventories.Get("")

	assert.Error(t, err)
}

func TestGetDefaultOk(t *testing.T) {
	inv, err := inventories.GetDefault()

	assert.Equal(t, inv.Namespace, "default")
	assert.Nil(t, err)
}

func TestListOk(t *testing.T) {
	_, err := inventories.List()

	assert.Nil(t, err)
}

func TestUpdateOk(t *testing.T) {
	def, _ := inventories.GetDefault()

	assert.Nil(t, inventories.Update("test", def))
}

func TestResetOk(t *testing.T) {
	inv, err := inventories.Reset("test")

	assert.Equal(t, "test", inv.Namespace)
	assert.Nil(t, err)
}
