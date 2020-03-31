package kubernetes

import (
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

func (c *jobRepository) Delete(namespace, resourceName string) error {
	if err := c.kubernetes.BatchV1().Jobs(namespace).Delete(resourceName, &metav1.DeleteOptions{}); err != nil {
		return err
	}
	return nil
}
