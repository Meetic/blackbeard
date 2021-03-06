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
	Events(listener string) chan NamespaceEvent
	AddListener(name string)
	RemoveListener(name string) error
	Emit(event NamespaceEvent)
	WatchNamespaces()
}

// NamespaceRepository defined the way namespace area actually managed.
type NamespaceRepository interface {
	Create(namespace string) error
	Get(namespace string) (*Namespace, error)
	ApplyConfig(namespace string, configPath string) error
	Delete(namespace string) error
	List() ([]Namespace, error)
	WatchPhase(emit EventEmitter) error
}

type namespaceService struct {
	namespaces      NamespaceRepository
	pods            PodRepository
	deployments     DeploymentRepository
	statefulsets    StatefulsetRepository
	jobs            JobRepository
	namespaceEvents map[string]chan NamespaceEvent
}

// NamespaceStatus represent namespace with percentage of pods running and status phase (Active or Terminating)
type NamespaceStatus struct {
	Status int    `json:"status"`
	Phase  string `json:"phase"`
}

const (
	EventStatusUpdate string = "STATUS.UPDATE"
	EventStatusReady  string = "STATUS.READY"
)

// NamespaceEvent represent a namespace event happened on kubernetes cluster
type NamespaceEvent struct {
	Type      string `json:"type"`
	Namespace string `json:"namespace"`
	Phase     string `json:"phase"`
	Status    int    `json:"status"`
}

type EventEmitter func(event NamespaceEvent)

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

// AddListener register a new listener
func (ns *namespaceService) AddListener(name string) {
	if ns.namespaceEvents == nil {
		ns.namespaceEvents = make(map[string]chan NamespaceEvent)
	}

	ns.namespaceEvents[name] = make(chan NamespaceEvent)
}

// RemoveListener close channel and remove listener
func (ns *namespaceService) RemoveListener(name string) error {
	if listener, ok := ns.namespaceEvents[name]; ok {
		close(listener)
		delete(ns.namespaceEvents, name)
		return nil
	}

	return fmt.Errorf("listener does not exist")
}

// Emit publish event to all registered listeners
func (ns *namespaceService) Emit(event NamespaceEvent) {
	logrus.WithFields(logrus.Fields{
		"component": "emmiter",
		"event":     event.Type,
		"namespace": event.Namespace,
	}).Debugf("new status : %d | new phase %s", event.Status, event.Phase)

	for _, ch := range ns.namespaceEvents {
		go func(handler chan NamespaceEvent) {
			handler <- event
		}(ch)
	}
}

// WatchNamespaces create a watcher on namespace events from kubernetes cluster and send result over received channel
func (ns *namespaceService) WatchNamespaces() {
	var wg sync.WaitGroup
	wg.Add(2)

	logrus.WithFields(logrus.Fields{"component": "watcher"}).Debug("Starting watch status")

	go ns.watchStatus()

	logrus.WithFields(logrus.Fields{"component": "watcher"}).Debug("Starting watch phase")

	go func() {
		for {
			ns.namespaces.WatchPhase(ns.Emit)
			logrus.WithFields(logrus.Fields{
				"component": "watcher",
			}).Debug("Phase watcher restarted due to closed http connection")
		}
	}()

	wg.Wait()
}

func (ns *namespaceService) watchStatus() {
	ticker := time.NewTicker(5 * time.Second)

	defer ticker.Stop()

	var lastEvents []NamespaceEvent

	for range ticker.C {
		logrus.WithFields(logrus.Fields{"component": "watcher"}).Debug("Watch status tick")

		namespaces, err := ns.List()

		if err != nil {
			continue
		}

		var events []NamespaceEvent

		for _, n := range namespaces {
			events = append(events, NamespaceEvent{
				Type:      getEventType(n.Status),
				Namespace: n.Name,
				Phase:     n.Phase,
				Status:    n.Status,
			})
		}

		returnedEvents := compareEvents(events, lastEvents)
		lastEvents = events

		logrus.WithFields(logrus.Fields{"component": "watcher"}).Debugf("Namespace status events to send %d", len(returnedEvents))

		for _, e := range returnedEvents {
			ns.Emit(e)
		}
	}
}

func getEventType(status int) string {
	if status == 100 {
		return EventStatusReady
	}

	return EventStatusUpdate
}

func compareEvents(now []NamespaceEvent, before []NamespaceEvent) []NamespaceEvent {
	var diff []NamespaceEvent

	for _, s1 := range now {
		found := false
		statusDiff := true
		for _, s2 := range before {
			if s1.Namespace == s2.Namespace {
				found = true
				if s1.Status == s2.Status {
					statusDiff = false
				}
				break
			}
		}

		// not found before or status diff
		if !found || statusDiff {
			diff = append(diff, s1)
		}
	}

	return diff
}

// Create creates a kubernetes namespace
func (ns *namespaceService) Create(n string) error {
	err := ns.namespaces.Create(n)

	if err != nil {
		return ErrorCreateNamespace{err.Error()}
	}

	return nil
}

// Events
func (ns *namespaceService) Events(listener string) chan NamespaceEvent {
	if ch, ok := ns.namespaceEvents[listener]; ok {
		return ch
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

// ErrorCreateNamespace represents an error due to a namespace creation failure on kubernetes cluster
type ErrorCreateNamespace struct {
	Msg string
}

// Error returns the error message
func (err ErrorCreateNamespace) Error() string {
	return err.Msg
}
