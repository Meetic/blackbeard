package kubectl

import "github.com/Meetic/blackbeard/pkg/blackbeard"

//Client represents a kubectl client for namespaceService
type Client struct {
	configPath       string
	namespaceService NamespaceService
}

//NewClient returns a new kubectl client
func NewClient(configPath string) *Client {
	c := &Client{
		configPath: configPath,
	}

	c.namespaceService.configPath = c.configPath

	return c
}

//NamespaceService returns the kubectl implementation of namespaceService
func (c *Client) NamespaceService() blackbeard.NamespaceService { return &c.namespaceService }
