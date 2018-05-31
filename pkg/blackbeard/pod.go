package blackbeard

//Pods represent a list of pods.
type Pods []Pod

//Pod represent a Kubernetes pod.
type Pod struct {
	Name   string
	Status string
}

type PodRepository interface {
	List(string) (Pods, error)
}
