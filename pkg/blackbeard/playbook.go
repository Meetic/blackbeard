package blackbeard

import "text/template"

// ConfigTemplate represents a set of kubernetes configuration template.
// Usually, Template is expected to be golang template of yaml.
type ConfigTemplate struct {
	Name     string
	Template *template.Template
}

type PlaybookService interface {
	GetDefault() (Inventory, error)
	GetTemplate() ([]ConfigTemplate, error)
}

type PlaybookRepository interface {
	GetDefault() (Inventory, error)
	GetTemplate() ([]ConfigTemplate, error)
}

type playbookService struct {
	playbooks PlaybookRepository
}

func NewPlaybookService(playbooks PlaybookRepository) PlaybookService {
	return &playbookService{
		playbooks: playbooks,
	}
}

// GetTemplate returns the templates of a playbook
func (ps *playbookService) GetTemplate() ([]ConfigTemplate, error) {
	return ps.playbooks.GetTemplate()
}

// GetDefault returns the default inventory of a playbook
func (ps *playbookService) GetDefault() (Inventory, error) {
	return ps.playbooks.GetDefault()
}
