package files

import "github.com/Meetic/blackbeard/pkg/blackbeard"

//Client represents a file client for configService and inventoryService
type Client struct {
	templatePath     string
	configPath       string
	inventoryPath    string
	defaultsPaths    string
	configService    ConfigService
	inventoryService InventoryService
}

//NewClient return a new file client
func NewClient(templatePath, configPath, inventoryPath, defaultsPaths string) *Client {
	c := &Client{
		templatePath:  templatePath,
		configPath:    configPath,
		inventoryPath: inventoryPath,
		defaultsPaths: defaultsPaths,
	}

	c.configService.configPath = c.configPath
	c.configService.templatePath = c.templatePath
	c.inventoryService.inventoryPath = c.inventoryPath
	c.inventoryService.defaultsPath = c.defaultsPaths

	return c
}

//InventoryService return the file implementation of inventoryService
func (c *Client) InventoryService() blackbeard.InventoryService { return &c.inventoryService }

//ConfigService return the file implementation of configService
func (c *Client) ConfigService() blackbeard.ConfigService { return &c.configService }
