package mock

import (
	"github.com/Meetic/blackbeard/pkg/playbook"
)

type configRepository struct{}

// NewConfigRepository returns a new Mock ConfigRepository
func NewConfigRepository() playbook.ConfigRepository {
	return &configRepository{}
}

func (cr *configRepository) Save(namespace string, configs []playbook.Config) error {
	return nil
}

func (cr *configRepository) Delete(namespace string) error {
	return nil
}
