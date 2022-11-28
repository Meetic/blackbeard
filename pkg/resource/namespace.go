package resource

import (
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type Namespace struct {
	Name   string
	Phase  string
	Status int
}

// NamespaceService defined the way namespace are managed.
type NamespaceService interface {
	Create(namespace string) error
	ApplyConfig(namespace string, configPath string) error
	Delete(namespace string) error
	GetStatus(namespace string) (*NamespaceStatus, error)
	List() ([]Namespace, error)
	Watch(events chan NamespaceEvent)
}

// NamespaceRepository defined the way namespace area actually managed.
type NamespaceRepository interface {
	Create(namespace string) error
	Get(namespace string) (*Namespace, error)
	ApplyConfig(namespace string, configPath string) error
	Delete(namespace string) error
	List() ([]Namespace, error)
	Watch(events chan<- NamespaceEvent) error
}

type namespaceService struct {
	namespaces   NamespaceRepository
	pods         PodRepository
	deployments  DeploymentRepository
	statefulsets StatefulsetRepository
	jobs         JobRepository
}

// NamespaceStatus represent namespace with percentage of pods running and status phase (Active or Terminating)
type NamespaceStatus struct {
	Status int    `json:"status"`
	Phase  string `json:"phase"`
}

type NamespaceEvent struct {
	Namespace string
	Type      string
}

// NewNamespaceService creates a new NamespaceService
func NewNamespaceService(
	namespaces NamespaceRepository,
	pods PodRepository,
	deployments DeploymentRepository,
	statefulsets StatefulsetRepository,
	jobs JobRepository,
) NamespaceService {

	ns := &namespaceService{
		namespaces:   namespaces,
		pods:         pods,
		deployments:  deployments,
		statefulsets: statefulsets,
		jobs:         jobs,
	}

	return ns
}

// Create creates a kubernetes namespace
func (ns *namespaceService) Create(n string) error {
	err := ns.namespaces.Create(n)

	if err != nil {
		return ErrorCreateNamespace{err.Error()}
	}

	return nil
}

// ApplyConfig apply kubernetes configurations to the given namespace.
// Warning : For now, this method takes a configPath as parameter. This parameter is the directory containing configs in a playbook
// This may change since the NamespaceService should not be aware that configs are stored in files.
func (ns *namespaceService) ApplyConfig(namespace, configPath string) error {
	return ns.namespaces.ApplyConfig(namespace, configPath)
}

// Delete deletes a kubernetes namespace
func (ns *namespaceService) Delete(namespace string) error {
	return ns.namespaces.Delete(namespace)
}

// List returns a slice of namespace from the kubernetes package and enrich each of the
// returned namespace with their status.
func (ns *namespaceService) List() ([]Namespace, error) {
	namespaces, err := ns.namespaces.List()
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup

	for i := range namespaces {
		wg.Add(1)

		go func(index int) {
			status, err := ns.GetStatus(namespaces[index].Name)

			if err != nil {
				namespaces[index].Status = 0
			}

			namespaces[index].Status = status.Status
			wg.Done()
		}(i)
	}

	wg.Wait()

	return namespaces, nil
}

// GetStatus returns the status of an inventory
// The status is an int that represents the percentage of pods in a "running" state inside the given namespace
func (ns *namespaceService) GetStatus(namespace string) (*NamespaceStatus, error) {

	// get namespace state
	n, err := ns.namespaces.Get(namespace)
	if err != nil {
		return nil, fmt.Errorf("namespace get status: %v", err)
	}

	if n.Phase == "Terminating" {
		return &NamespaceStatus{0, n.Phase}, nil
	}

	dps, errDps := ns.deployments.List(namespace)
	sfs, errSfs := ns.statefulsets.List(namespace)
	jbs, errJbs := ns.jobs.List(namespace)

	if errDps != nil || errSfs != nil || errJbs != nil {
		return &NamespaceStatus{0, ""}, fmt.Errorf("namespace get status: list deployments, statefulsets or jobs: %v", err)
	}

	totalApps := len(dps) + len(sfs) + len(jbs)

	if totalApps == 0 {
		return &NamespaceStatus{0, n.Phase}, nil
	}

	var i int

	for _, dp := range dps {
		if dp.Status == DeploymentReady {
			i++
		}
	}

	for _, sf := range sfs {
		if sf.Status == StatefulsetReady {
			i++
		}
	}

	for _, job := range jbs {
		if job.Status == JobReady {
			i++
		}
	}

	status := i * 100 / totalApps

	return &NamespaceStatus{status, n.Phase}, nil
}

func (ns *namespaceService) Watch(events chan NamespaceEvent) {
	ticker := time.NewTicker(5 * time.Second)
	defer close(events)

	for range ticker.C {
		err := ns.namespaces.Watch(events)
		if err != nil {
			ticker.Stop()
		}

		logrus.
			WithFields(logrus.Fields{"component": "watcher"}).
			Debug("watch namespace restarted")
	}

	logrus.
		WithFields(logrus.Fields{"component": "watcher"}).
		Error("watch namespace stopped due to error")
}

// ErrorCreateNamespace represents an error due to a namespace creation failure on kubernetes cluster
type ErrorCreateNamespace struct {
	Msg string
}

// Error returns the error message
func (err ErrorCreateNamespace) Error() string {
	return err.Msg
}
