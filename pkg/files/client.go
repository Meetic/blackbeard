package files

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Meetic/blackbeard/pkg/playbook"
)

const (
	templateDir  = "templates"
	configDir    = "configs"
	inventoryDir = "inventories"
	defaultFile  = "defaults.json"
)

type Client struct {
	configs       playbook.ConfigRepository
	inventories   playbook.InventoryRepository
	playbooks     playbook.PlaybookRepository
	inventoryPath string
	configPath    string
}

func NewClient(wd string) (*Client, error) {
	if ok, _ := fileExists(wd); ok != true {
		return &Client{}, fmt.Errorf("Your specified working dir does not exit : %s", wd)
	}

	templatePath := filepath.Join(wd, templateDir)
	configPath := filepath.Join(wd, configDir)
	inventoryPath := filepath.Join(wd, inventoryDir)
	defaultsPath := filepath.Join(wd, defaultFile)

	if ok, _ := fileExists(templatePath); ok != true {
		return &Client{}, fmt.Errorf("A playbook must contains a `%s` dir. No one has been found.\n"+
			"Please check the playbook or change the working directory using the --dir option.", templateDir)
	}

	if ok, _ := fileExists(defaultsPath); ok != true {
		return &Client{}, fmt.Errorf("Your working directory must contains a `%s` file.\n"+
			"Please check the playbook or change the working directory using the --dir option.", defaultFile)
	}

	if ok, _ := fileExists(configPath); ok != true {
		if err := os.Mkdir(configPath, 0755); err != nil {
			return &Client{}, fmt.Errorf("Impossible to create the %s directory. Please check directory rights.", configDir)
		}
	}

	if ok, _ := fileExists(inventoryPath); ok != true {
		if err := os.Mkdir(inventoryPath, 0755); err != nil {
			return &Client{}, fmt.Errorf("Impossible to create the %s directory. Please check directory rights.", inventoryDir)
		}
	}

	return &Client{
		configs:       NewConfigRepository(configPath),
		inventories:   NewInventoryRepository(inventoryPath),
		playbooks:     NewPlaybookRepository(templatePath, defaultsPath),
		inventoryPath: inventoryPath,
		configPath:    configPath,
	}, nil
}

func (c *Client) Configs() playbook.ConfigRepository {
	return c.configs
}

func (c *Client) Inventories() playbook.InventoryRepository {
	return c.inventories
}

func (c *Client) Playbooks() playbook.PlaybookRepository {
	return c.playbooks
}

// InventoryPath returns the inventory path for the current playbook
func (c *Client) InventoryPath() string {
	return c.inventoryPath
}

// ConfigPath return the config path for the current playbook
func (c *Client) ConfigPath() string {
	return c.configPath
}

func fileExists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return true, err
	}

	return true, nil
}
