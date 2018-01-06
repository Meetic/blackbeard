package files

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/Meetic/blackbeard/pkg/blackbeard"
)

//ConfigService is used to managed kubernetes namespace configuration
type ConfigService struct {
	templatePath string
	configPath   string
}

//Ensure that ConfigService implements the interface
var _ blackbeard.ConfigService = (*ConfigService)(nil)

//Apply generate a list of kubernetes config file from a given inventory
func (cs *ConfigService) Apply(inv blackbeard.Inventory) error {

	//Get template list
	templates, _ := filepath.Glob(cs.templatePath + "/*.tpl")

	if templates == nil {
		return fmt.Errorf("no template files found in directory %s", cs.templatePath)
	}

	//Create config dir for a given namespace
	configDir := filepath.Join(cs.configPath, inv.Namespace)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if e := os.Mkdir(configDir, os.ModePerm); e != nil {
			return fmt.Errorf("the configs dir '%s' could not be created : %s", configDir, e.Error())
		}
	}

	for _, templ := range templates {

		tpl, err := template.ParseFiles(templ)
		if err != nil {
			return err
		}

		//create config file from tpl by removing the .tpl extension
		ext := filepath.Ext(templ)
		_, configFile := filepath.Split(templ[0 : len(templ)-len(ext)])

		f, err := os.Create(filepath.Join(configDir, configFile))
		if err != nil {
			return err
		}
		defer f.Close()

		//execute the template
		err = tpl.Execute(f, inv)
		if err != nil {
			return err
		}

	}
	return nil
}

//Delete remove a config directory
//if the specified config dir does not exist, Delete return nil and does nothing.
func (cs *ConfigService) Delete(namespace string) error {
	if !cs.exists(namespace) {
		return nil
	}
	return os.RemoveAll(filepath.Join(cs.configPath, namespace))
}

//exists return true if a config dir for the given namespace already exist.
//Else, it return false.
func (cs *ConfigService) exists(namespace string) bool {
	if _, err := os.Stat(cs.path(namespace)); os.IsNotExist(err) {
		return false
	} else if err == nil {
		return true
	}
	return false
}

//Path return the config dir path of a given namespace
func (cs *ConfigService) path(namespace string) string {
	return filepath.Join(cs.configPath, namespace)
}
