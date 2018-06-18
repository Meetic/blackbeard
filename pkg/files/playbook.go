package files

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"text/template"

	"github.com/Meetic/blackbeard/pkg/playbook"
)

type playbooks struct {
	templatePath string
	defaultsPath string
}

func NewPlaybookRepository(templatePath, defaultsPath string) playbook.PlaybookRepository {
	return &playbooks{
		templatePath,
		defaultsPath,
	}

}

// GetTemplate returns the templates from the playbook
func (p *playbooks) GetTemplate() ([]playbook.ConfigTemplate, error) {

	//Get template list
	templates, _ := filepath.Glob(fmt.Sprintf("%s/*%s", p.templatePath, tplSuffix))

	if templates == nil {
		return nil, fmt.Errorf("no template files found in directory %s", p.templatePath)
	}

	var cfgTpl []playbook.ConfigTemplate

	for _, templ := range templates {

		tpl, err := template.ParseFiles(templ)
		if err != nil {
			return nil, err
		}

		//create config file from tpl by removing the .tpl extension
		ext := filepath.Ext(templ)
		_, configFile := filepath.Split(templ[0 : len(templ)-len(ext)])

		config := playbook.ConfigTemplate{
			Name:     configFile,
			Template: tpl,
		}

		cfgTpl = append(cfgTpl, config)
	}

	return cfgTpl, nil

}

// GetDefault reads the default inventory file and return an Inventory where namespace is set to "default"
func (p *playbooks) GetDefault() (playbook.Inventory, error) {

	defaults, err := ioutil.ReadFile(p.defaultsPath)

	if err != nil {
		return playbook.Inventory{}, playbook.NewErrorReadingDefaultsFile(err)
	}

	var inventory playbook.Inventory

	if err := json.Unmarshal(defaults, &inventory); err != nil {
		return playbook.Inventory{}, playbook.NewErrorReadingDefaultsFile(err)
	}

	return inventory, nil
}
