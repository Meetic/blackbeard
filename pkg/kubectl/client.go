package kubectl

import "github.com/Meetic/blackbeard/pkg/blackbeard"

//Client represents a kubectl client for namespaceService
type Client struct {
	configPath           string
	namespaceConfService NamespaceConfigurationService
}

//NewClient returns a new kubectl client
func NewClient(configPath string) *Client {
	c := &Client{
		configPath: configPath,
	}

	c.namespaceConfService.configPath = c.configPath

	return c
}

//NamespaceService returns the kubectl implementation of namespaceService
func (c *Client) NamespaceConfigurationService() blackbeard.NamespaceConfigurationService {
	return &c.namespaceConfService
}
