package kubernetes

import (
	"log"
	"net/url"
	"strings"

	"github.com/Meetic/blackbeard/pkg/blackbeard"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	running = "running"
)

type Client struct {
	configFile      string
	resourceService ResourceService
}

//Ensure that ResourceService implements the interface
var _ blackbeard.KubernetesClient = (*Client)(nil)

//NewClient return a new kubernetes client
func NewClient(configFile string) *Client {
	c := &Client{
		configFile: configFile,
	}

	config, _ := clientcmd.BuildConfigFromFlags("", c.configFile)

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	c.resourceService.client = clientSet
	u, err := url.Parse(config.Host)

	if err != nil {
		log.Fatal(err)
	}

	c.resourceService.host = strings.Split(u.Host, ":")[0]
	return c
}

//ResourceService returns the kubernetes resource service
func (c *Client) ResourceService() blackbeard.ResourceService { return &c.resourceService }
