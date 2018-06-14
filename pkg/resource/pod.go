package resource

// Pods represent a list of pods.
type Pods []Pod

// Pod represent a Kubernetes pod.
// The status is the pod phase status. It could be :
// * running
// * pending
// etc...
type Pod struct {
	Name   string
	Status string
}

// PodRepository represents the way Pods are managed
type PodRepository interface {
	List(string) (Pods, error)
}
