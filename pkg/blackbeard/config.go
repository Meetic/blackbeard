package blackbeard

import (
	"bytes"
	"errors"
	"text/template"
	"time"
)

type Config struct {
	Name   string
	Values string
}

type ConfigTemplate struct {
	Name     string
	Template *template.Template
}

type Release struct {
	Date string `json:"date"`
}

type InventoryRelease struct {
	Namespace string                 `json:"namespace"`
	Values    map[string]interface{} `json:"values"`
	Release   Release                `json:"release"`
}

//ConfigService define the way configuration should be managed
type ConfigService interface {
	Generate(Inventory) error
	Delete(string) error
}

type ConfigRepository interface {
	GetTemplate() ([]ConfigTemplate, error)
	Save(namespace string, configs []Config) error
	Delete(string) error
}

type configService struct {
	configs ConfigRepository
}

func NewConfigService(configs ConfigRepository) ConfigService {
	return &configService{
		configs,
	}
}

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

func (cs *configService) Delete(namespace string) error {
	return cs.configs.Delete(namespace)
}
