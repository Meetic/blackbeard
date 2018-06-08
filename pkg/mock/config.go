package mock

import (
	"github.com/Meetic/blackbeard/pkg/blackbeard"
)

type configRepository struct{}

// NewConfigRepository returns a new Mock ConfigRepository
func NewConfigRepository() blackbeard.ConfigRepository {
	return &configRepository{}
}

func (cr *configRepository) Save(namespace string, configs []blackbeard.Config) error {
	return nil
}

func (cr *configRepository) Delete(namespace string) error {
	return nil
}
