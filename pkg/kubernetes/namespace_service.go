package kubernetes

import (
	"github.com/Meetic/blackbeard/pkg/blackbeard"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type NamespaceService struct {
	client kubernetes.Interface
}

//Ensure that ResourceService implements the interface
var _ blackbeard.NamespaceService = (*NamespaceService)(nil)

//Create create a namespace
func (rs *NamespaceService) Create(namespace string) error {
	_, err := rs.client.CoreV1().Namespaces().Create(&v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespace}})
	return err
}

//Delete delete a given namespace
func (rs *NamespaceService) Delete(namespace string) error {
	err := rs.client.CoreV1().Namespaces().Delete(namespace, &metav1.DeleteOptions{})
	switch t := err.(type) {
	case *errors.StatusError:
		return nil
	case *errors.UnexpectedObjectError:
		return nil
	default:
		return t
	}
}
