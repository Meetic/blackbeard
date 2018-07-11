package resource

import (
	"k8s.io/api/core/v1"
)

// Pods represent a list of pods.
type Pods []Pod

// Pod represent a Kubernetes pod.
// The status is the pod phase status. It could be :
// * running
// * pending
// etc...
type Pod struct {
	Name   string
	Status v1.PodPhase
}

type podService struct {
	pods PodRepository
}

// PodRepository represents the way Pods are managed
type PodRepository interface {
	List(string) (Pods, error)
	//Get(string, string) (Pod, error)
}

type PodService interface {
	List(string) (Pods, error)
}

//NewPodService returns a new PodService
func NewPodService(pods PodRepository) PodService {
	return &podService{
		pods: pods,
	}
}

// List returns the list of pods in a kubernetes namespace with their associated status.
func (ps *podService) List(namespace string) (Pods, error) {
	return ps.pods.List(namespace)
}
