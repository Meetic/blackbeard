package kubernetes

import (
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/Meetic/blackbeard/pkg/resource"
)

const (
	configDir  = ".kube"
	configFile = "config"
)

type Client struct {
	kubernetes kubernetes.Interface
	namespaces resource.NamespaceRepository
	pods       resource.PodRepository
	services   resource.ServiceRepository
	cluster    resource.ClusterRepository
}

// NewClient return a new kubernetes client
func NewClient(configFilePath string) (*Client, error) {

	config, err := clientcmd.BuildConfigFromFlags("", configFilePath)
	if err != nil {
		return &Client{}, err
	}

	config.QPS = float32(250)
	config.Burst = 500

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return &Client{}, err
	}

	return &Client{
		kubernetes: clientSet,
		namespaces: NewNamespaceRepository(clientSet),
		pods:       NewPodRepository(clientSet),
		services:   NewServiceRepository(clientSet, GetKubernetesHost(configFilePath)),
		cluster:    NewClusterRepository(),
	}, nil
}

func (c *Client) Namespaces() resource.NamespaceRepository {
	return c.namespaces
}

func (c *Client) Pods() resource.PodRepository {
	return c.pods
}

func (c *Client) Services() resource.ServiceRepository {
	return c.services
}

func (c *Client) Cluster() resource.ClusterRepository {
	return c.cluster
}

// KubeConfigDefaultPath return the kubernetes default config path
func KubeConfigDefaultPath() string {
	return filepath.Join(homeDir(), configDir, configFile)
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

// GetKubernetesHost return the kubernetes cluster domain name used in the ~/.kube/config file
// The returned host takes the form : mydomainname.com
// Notice : this is just the host, without any schema or port.
func GetKubernetesHost(configFilePath string) string {

	config, _ := clientcmd.BuildConfigFromFlags("", configFilePath)

	u, err := url.Parse(config.Host)
	if err != nil {
		log.Fatal(err)
	}

	return strings.Split(u.Host, ":")[0]
}
