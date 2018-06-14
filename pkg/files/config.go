package files

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Meetic/blackbeard/pkg/playbook"
)

const (
	tplSuffix = ".tpl"
)

type configs struct {
	configPath string
}

// NewConfigRepository returns a new ConfigRepository
// It takes as parameters the directory where configs are stored.
// Typically, the templates files for a given playbook are in a "templates" directory at the root of the playbook
// and configs are stored in a "configs" directory located at the root of the playbook
func NewConfigRepository(configPath string) playbook.ConfigRepository {
	return &configs{
		configPath: configPath,
	}
}

// Save writes kubernetes configs for a given namespace in files.
// files are named after the Config.Name value
func (cr *configs) Save(namespace string, configs []playbook.Config) error {

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

// Delete remove a config directory
// if the specified config dir does not exist, Delete return nil and does nothing.
func (cr *configs) Delete(namespace string) error {
	if !cr.exists(namespace) {
		return nil
	}
	return os.RemoveAll(filepath.Join(cr.configPath, namespace))
}

// exists return true if a config dir for the given namespace already exist.
// Else, it return false.
func (cr *configs) exists(namespace string) bool {
	if _, err := os.Stat(cr.path(namespace)); os.IsNotExist(err) {
		return false
	} else if err == nil {
		return true
	}
	return false
}

// path return the config dir path of a given namespace
func (cr *configs) path(namespace string) string {
	return filepath.Join(cr.configPath, namespace)
}
