package kubernetes

import (
	"log"
	"net/url"
	"strings"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// NewClient return a new kubernetes client
func NewClient(configFile string) kubernetes.Interface {

	config, _ := clientcmd.BuildConfigFromFlags("", configFile)

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	return clientSet
}

// GetKubernetesHost return the kubernetes cluster domain name used in the ~/.kube/config file
func GetKubernetesHost(configFile string) string {

	config, _ := clientcmd.BuildConfigFromFlags("", configFile)

	u, err := url.Parse(config.Host)
	if err != nil {
		log.Fatal(err)
	}

	return strings.Split(u.Host, ":")[0]
}
