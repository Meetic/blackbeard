package blackbeard

import (
	"bytes"
	"errors"
	"text/template"
	"time"
)

// Config represents a set of kubernetes configuration.
// Usually, Values are expected to be yaml.
type Config struct {
	Name   string
	Values string
}

// ConfigTemplate represents a set of kubernetes configuration template.
// Usually, Template is expected to be golang template of yaml.
type ConfigTemplate struct {
	Name     string
	Template *template.Template
}

// Release represents information related to an inventory release.
// An inventory may evolve with time. We want to keep trace of those evolution
// and we may inject data specific a release in the templates
type Release struct {
	Date string `json:"date"`
}

// InventoryRelease represents an inventory enriched with release data.
type InventoryRelease struct {
	Namespace string                 `json:"namespace"`
	Values    map[string]interface{} `json:"values"`
	Release   Release                `json:"release"`
}

// ConfigService define the way configuration are managed
type ConfigService interface {
	Generate(Inventory) error
	Delete(namespace string) error
}

// ConfigRepository represents a service that implements configs management
type ConfigRepository interface {
	GetTemplate() ([]ConfigTemplate, error)
	Save(namespace string, configs []Config) error
	Delete(namespace string) error
}

type configService struct {
	configs ConfigRepository
}

// NewConfigService creates a ConfigService
func NewConfigService(configs ConfigRepository) ConfigService {
	return &configService{
		configs,
	}
}

// Generate creates a set of kubernetes configurations by applying an InventoryRelease to
// Templates. It read each template, create an InventoryRelease for the given Inventory
// and apply it to the template in order to generate a set of kubernetes configurations.
func (cs *configService) Generate(inv Inventory) error {

	if inv.Namespace == "" {
		return errors.New("an namespace must be specified in the inventory")
	}

	tpls, err := cs.configs.GetTemplate()
	if err != nil {
		return err
	}

	invRelease := InventoryRelease{
		inv.Namespace,
		inv.Values,
		Release{
			Date: time.Now().Format(time.RFC850),
		},
	}

	var configs []Config

	for _, tpl := range tpls {

		confVal := bytes.Buffer{}

		tpl.Template.Execute(&confVal, invRelease)

		conf := Config{
			Name:   tpl.Name,
			Values: confVal.String(),
		}

		configs = append(configs, conf)
	}

	return cs.configs.Save(inv.Namespace, configs)

}

// Delete delete kubernetes configs for the given namespace.
func (cs *configService) Delete(namespace string) error {
	return cs.configs.Delete(namespace)
}
