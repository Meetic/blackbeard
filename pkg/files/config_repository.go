package files

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/Meetic/blackbeard/pkg/blackbeard"
)

const (
	tplSuffix = ".tpl"
)

type configRepository struct {
	templatePath string
	configPath   string
}

func NewConfigRepository(templatePath, configPath string) blackbeard.ConfigRepository {
	return &configRepository{
		templatePath: templatePath,
		configPath:   configPath,
	}
}

func (cr *configRepository) GetTemplate() ([]blackbeard.ConfigTemplate, error) {

	//Get template list
	templates, _ := filepath.Glob(fmt.Sprintf("%s/*%s", cr.templatePath, tplSuffix))

	if templates == nil {
		return nil, fmt.Errorf("no template files found in directory %s", cr.templatePath)
	}

	var cfgTpl []blackbeard.ConfigTemplate

	for _, templ := range templates {

		tpl, err := template.ParseFiles(templ)
		if err != nil {
			return nil, err
		}

		//create config file from tpl by removing the .tpl extension
		ext := filepath.Ext(templ)
		_, configFile := filepath.Split(templ[0 : len(templ)-len(ext)])

		config := blackbeard.ConfigTemplate{
			Name:     configFile,
			Template: tpl,
		}

		cfgTpl = append(cfgTpl, config)
	}

	return cfgTpl, nil

}

func (cr *configRepository) Save(namespace string, configs []blackbeard.Config) error {

	//Create config dir for a given namespace
	configDir := filepath.Join(cr.configPath, namespace)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if e := os.Mkdir(configDir, os.ModePerm); e != nil {
			return fmt.Errorf("the configs dir '%s' could not be created : %s", configDir, e.Error())
		}
	}

	for _, config := range configs {
		f, err := os.Create(filepath.Join(configDir, config.Name))
		if err != nil {
			return err
		}

		_, errorWrite := f.Write([]byte(config.Values))
		if errorWrite != nil {
			return errorWrite
		}

		f.Close()
	}

	return nil
}

//Delete remove a config directory
//if the specified config dir does not exist, Delete return nil and does nothing.
func (cr *configRepository) Delete(namespace string) error {
	if !cr.exists(namespace) {
		return nil
	}
	return os.RemoveAll(filepath.Join(cr.configPath, namespace))
}

//exists return true if a config dir for the given namespace already exist.
//Else, it return false.
func (cr *configRepository) exists(namespace string) bool {
	if _, err := os.Stat(cr.path(namespace)); os.IsNotExist(err) {
		return false
	} else if err == nil {
		return true
	}
	return false
}

//Path return the config dir path of a given namespace
func (cr *configRepository) path(namespace string) string {
	return filepath.Join(cr.configPath, namespace)
}
