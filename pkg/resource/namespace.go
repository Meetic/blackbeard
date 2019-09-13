package resource

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
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
	WatchNamespaces() error
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
	Type       string `json:"type"`
	Namespace  string `json:"namespace"`
	Phase      string `json:"phase"`
	Status     int    `json:"status"`
	PodsStatus Pods   `json:"pods_status"`
}

type EventEmitter func(event NamespaceEvent)

// NewNamespaceService creates a new NamespaceService
func NewNamespaceService(namespaces NamespaceRepository, pods PodRepository) NamespaceService {

	ns := &namespaceService{
		namespaces: namespaces,
		pods:       pods,
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
func (ns *namespaceService) WatchNamespaces() error {
	go ns.watchStatus()

	go func() {
		for {
			ns.namespaces.WatchPhase(ns.Emit)
			logrus.WithFields(logrus.Fields{
				"component": "watcher",
			}).Debug("Phase watcher restarted due to closed http connection")
		}
	}()

	return nil
}

func (ns *namespaceService) watchStatus() error {
	ticker := time.NewTicker(10 * time.Second)

	defer ticker.Stop()

	var lastEvents []NamespaceEvent

	for range ticker.C {
		namespaces, err := ns.List()
		if err != nil {
			return err
		}

		var events []NamespaceEvent

		for _, n := range namespaces {
			pods, err := ns.pods.List(n.Name)
			if err != nil {
				return err
			}

			events = append(events, NamespaceEvent{
				Type:       getEventType(n.Status),
				Namespace:  n.Name,
				Phase:      n.Phase,
				Status:     n.Status,
				PodsStatus: pods,
			})
		}

		returnedEvents := compareEvents(events, lastEvents)
		lastEvents = events

		for _, e := range returnedEvents {
			ns.Emit(e)
		}
	}

	return nil
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

// List returns a slice of Namespace from the kubernetes package and enrich each of the
// returned namespace with its status.
func (ns *namespaceService) List() ([]Namespace, error) {
	namespaces, err := ns.namespaces.List()
	if err != nil {
		return nil, err
	}

	for i, namespace := range namespaces {
		status, err := ns.GetStatus(namespace.Name)
		if err != nil {
			return nil, err
		}

		namespaces[i].Status = status.Status
	}

	return namespaces, nil
}

// GetStatus returns the status of an inventory
// The status is an int that represents the percentage of pods in a "running" state inside the given namespace
func (ns *namespaceService) GetStatus(namespace string) (*NamespaceStatus, error) {

	// get namespace state
	n, err := ns.namespaces.Get(namespace)
	if err != nil {
		return &NamespaceStatus{0, ""}, err
	}

	if n.Phase == "Terminating" {
		return &NamespaceStatus{0, n.Phase}, nil
	}

	//  get pod's namespace
	pods, err := ns.pods.List(namespace)
	if err != nil {
		return &NamespaceStatus{0, ""}, err
	}

	totalPods := len(pods)

	if totalPods == 0 {
		return &NamespaceStatus{0, n.Phase}, nil
	}

	var i int

	for _, pod := range pods {
		if pod.Status == v1.PodRunning {
			i++
		}
	}

	status := i * 100 / totalPods

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
