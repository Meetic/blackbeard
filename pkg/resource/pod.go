package resource

import (
	"fmt"
	"strings"

	v1 "k8s.io/api/core/v1"
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
	List(namespace string) (Pods, error)
	Delete(namespace string, pod Pod) error
	//Get(string, string) (Pod, error)
}

type PodService interface {
	List(namespace string) (Pods, error)
	Find(namespace, deployment string) (Pod, error)
	Delete(namespace string, pod Pod) error
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

func (ps *podService) Find(namespace, deployment string) (Pod, error) {
	podList, err := ps.pods.List(namespace)

	if err != nil {
		return Pod{}, err
	}

	for _, pod := range podList {
		if strings.Contains(pod.Name, deployment) {
			return pod, nil
		}
	}

	return Pod{}, fmt.Errorf("no pod have been found for deployment name : %s", deployment)
}

func (ps *podService) Delete(namespace string, pod Pod) error {
	return ps.pods.Delete(namespace, pod)
}
