package kubernetes

import (
	"fmt"

	"k8s.io/api/batch/v1"

	"github.com/Meetic/blackbeard/pkg/resource"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type jobRepository struct {
	kubernetes kubernetes.Interface
}

// NewJobRepository returns a new JobRepository.
// The parameter is a go-client kubernetes client.
func NewJobRepository(kubernetes kubernetes.Interface) resource.JobRepository {
	return &jobRepository{
		kubernetes: kubernetes,
	}
}

func (c *jobRepository) List(namespace string) (resource.Jobs, error) {
	jl, err := c.kubernetes.BatchV1().Jobs(namespace).List(metav1.ListOptions{})

	if err != nil {
		return nil, fmt.Errorf("unable to list jobs: %v", err)
	}

	jobs := make(resource.Jobs, 0)

	for _, job := range jl.Items {
		status := resource.JobNotReady

		if len(job.Status.Conditions) > 0 && job.Status.Conditions[len(job.Status.Conditions)].Type == v1.JobComplete {
			status = resource.JobReady
		}

		jobs = append(jobs, resource.Job{
			Name:   job.Name,
			Status: status,
		})
	}

	return jobs, nil
}

func (c *jobRepository) Delete(namespace, resourceName string) error {
	pp := metav1.DeletePropagationBackground
	if err := c.kubernetes.BatchV1().Jobs(namespace).Delete(resourceName, &metav1.DeleteOptions{PropagationPolicy: &pp}); err != nil {
		return err
	}
	return nil
}
