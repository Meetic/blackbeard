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

	//Create config dir for a given namespace
	configDir := fmt.Sprintf("%s/%s", cs.configPath, inv.Namespace)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		os.Mkdir(configDir, os.ModePerm)
	}

	for _, templ := range templates {

		tpl, err := template.ParseFiles(templ)
		if err != nil {
			return err
		}

		//create config file from tpl by removing the .tpl extension
		ext := filepath.Ext(templ)
		_, configFile := filepath.Split(templ[0 : len(templ)-len(ext)])

		f, err := os.Create(configDir + "/" + configFile)
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
