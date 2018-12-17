package api

import (
	"fmt"
	"time"

	"github.com/Meetic/blackbeard/pkg/resource"
)

const (
	tickerDuration = 2 * time.Second
)

//Namespace represents a kubernetes namespace enrich with informations from the playbook.
type Namespace struct {
	//Name is the namespace name
	Name string
	//Phase is the namespace status phase. It could be "active" or "terminating"
	Phase string
	//Status is the namespace status. It is a percentage of runnning pods vs all pods in the namespace.
	Status int
	//Managed is true if the namespace as an associated inventory on the current playbook. False if not.
	Managed bool
}

// ListNamespaces returns a list of Namespace.
// For each kubernetes namespace, it checks if an associated inventory exists.
func (api *api) ListNamespaces() ([]Namespace, error) {
	nsList, err := api.namespaces.List()
	if err != nil {
		return nil, err
	}

	var namespaces []Namespace

	for _, ns := range nsList {

		namespace := Namespace{
			Name:    ns.Name,
			Phase:   ns.Phase,
			Status:  ns.Status,
			Managed: false,
		}

		if api.inventories.Exists(ns.Name) {
			namespace.Managed = true
		}

		namespaces = append(namespaces, namespace)
	}

	return namespaces, nil

}

type progress interface {
	Set(int) error
}

// WaitForNamespaceReady wait until all pods in the specified namespace are ready.
// And error is returned if the timeout is reach.
func (api *api) WaitForNamespaceReady(namespace string, timeout time.Duration, bar progress) error {

	ticker := time.NewTicker(tickerDuration)
	timerCh := time.NewTimer(timeout).C
	doneCh := make(chan bool)

	go func(bar progress, ns resource.NamespaceService, namespace string) {
		for range ticker.C {
			status, err := ns.GetStatus(namespace)
			if err != nil {
				ticker.Stop()
			}
			bar.Set(status.Status)
			if status.Status == 100 {
				doneCh <- true
			}
		}
	}(bar, api.namespaces, namespace)

	for {
		select {
		case <-timerCh:
			ticker.Stop()
			return fmt.Errorf("time out : Some pods are not yet ready")
		case <-doneCh:
			return nil
		}
	}
}
