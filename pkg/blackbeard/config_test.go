package blackbeard_test

import (
	"testing"

	"github.com/Meetic/blackbeard/pkg/blackbeard"
	"github.com/Meetic/blackbeard/pkg/mock"
	"github.com/stretchr/testify/assert"
)

func TestGenerateOk(t *testing.T) {

	configs := blackbeard.NewConfigService(mock.NewConfigRepository())
	inventories := mock.NewInventoryRepository()

	inv, _ := inventories.Get("test1")

	assert.Nil(t, configs.Generate(inv))
}

func TestGenerateEmptyNamespace(t *testing.T) {

	configs := blackbeard.NewConfigService(mock.NewConfigRepository())
	inventories := mock.NewInventoryRepository()

	inv, _ := inventories.Get("")

	assert.Error(t, configs.Generate(inv))
}

func TestDeleteOk(t *testing.T) {
	configs := blackbeard.NewConfigService(mock.NewConfigRepository())
	assert.Nil(t, configs.Delete("test"))
}
