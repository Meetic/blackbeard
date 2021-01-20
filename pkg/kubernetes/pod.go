package kubernetes

import (
	"context"
	"github.com/Meetic/blackbeard/pkg/resource"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type podRepository struct {
	kubernetes kubernetes.Interface
}

// NewPodRepository returns a new PodRepository.
// The parameter is a go-client kubernetes client.
func NewPodRepository(kubernetes kubernetes.Interface) resource.PodRepository {
	return &podRepository{
		kubernetes: kubernetes,
	}
}

// GetPods of all the pods in a given namespace.
// This method returns a Pods slice containing the pod name and the pod status (pod status phase).
func (pr *podRepository) List(n string) (resource.Pods, error) {
	// get all pods except job or cron jobs in a succeeded state
	podsList, err := pr.kubernetes.CoreV1().Pods(n).List(
		context.Background(),
		metav1.ListOptions{FieldSelector: "status.phase!=Succeeded"},
	)

	if err != nil {
		return nil, err
	}

	var pods resource.Pods

	for _, pod := range podsList.Items {

		pods = append(pods, resource.Pod{
			Name:   pod.ObjectMeta.Name,
			Status: pod.Status.Phase,
		})
	}

	return pods, nil
}
