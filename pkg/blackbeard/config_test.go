package blackbeard_test

import (
	"testing"

	"github.com/Meetic/blackbeard/pkg/blackbeard"
	"github.com/Meetic/blackbeard/pkg/mock"
	"github.com/stretchr/testify/assert"
)

var configs = blackbeard.NewConfigService(mock.NewConfigRepository(),
	blackbeard.NewPlaybookService(mock.NewPlaybookRepository()))

func TestGenerateOk(t *testing.T) {
	inventories := mock.NewInventoryRepository()

	inv, _ := inventories.Get("test1")

	assert.Nil(t, configs.Generate(inv))
}

func TestGenerateEmptyNamespace(t *testing.T) {

	inventories := mock.NewInventoryRepository()

	inv, _ := inventories.Get("")

	assert.Error(t, configs.Generate(inv))
}

func TestDeleteOk(t *testing.T) {
	assert.Nil(t, configs.Delete("test"))
}
