package kubecli

import "github.com/Meetic/blackbeard/pkg/blackbeard"

//Client represents a file client for configService and inventoryService
type Client struct {
	configPath       string
	namespaceService NamespaceService
}

//NewClient return a new file client
func NewClient(configPath string) *Client {
	c := &Client{
		configPath: configPath,
	}

	c.namespaceService.configPath = c.configPath

	return c
}

//NamespaceService return the kubectl implementation of namespaceService
func (c *Client) NamespaceService() blackbeard.NamespaceService { return &c.namespaceService }
