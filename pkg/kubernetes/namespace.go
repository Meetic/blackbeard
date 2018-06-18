package kubernetes

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/Meetic/blackbeard/pkg/resource"
	"k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	timeout = 60 * time.Second
)

type namespaceRepository struct {
	kubernetes kubernetes.Interface
}

// NewNamespaceRepository returns a new NamespaceRepository.
// The parameter is a go-client Kubernetes client
func NewNamespaceRepository(kubernetes kubernetes.Interface) resource.NamespaceRepository {
	return &namespaceRepository{
		kubernetes: kubernetes,
	}
}

// Create creates a namespace
func (ns *namespaceRepository) Create(namespace string) error {
	_, err := ns.kubernetes.CoreV1().Namespaces().Create(&v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespace}})
	return err
}

// Delete deletes a given namespace
func (ns *namespaceRepository) Delete(namespace string) error {
	err := ns.kubernetes.CoreV1().Namespaces().Delete(namespace, &metav1.DeleteOptions{})
	switch t := err.(type) {
	case *kerr.StatusError:
		return nil
	case *kerr.UnexpectedObjectError:
		return nil
	default:
		return t
	}
}

// List returns a slice of Namespace.
// Name is the namespace name from Kubernetes.
// Phase is the status phase.
// List returns an error if the namespace list could not be get from Kubernetes cluster.
func (ns *namespaceRepository) List() ([]resource.Namespace, error) {
	nsList, err := ns.kubernetes.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var namespaces []resource.Namespace
	for _, ns := range nsList.Items {
		namespaces = append(namespaces, resource.Namespace{
			Name:  ns.GetName(),
			Phase: string(ns.Status.Phase),
		})
	}

	return namespaces, nil
}

// ApplyConfig loads configuration files into kubernetes
func (ns *namespaceRepository) ApplyConfig(namespace, configPath string) error {

	err := execute(fmt.Sprintf("kubectl apply -f %s -n %s", filepath.Join(configPath, namespace), namespace), timeout)
	if err != nil {
		return fmt.Errorf("the namespace could not be configured : %v", err)
	}

	return nil
}

func execute(c string, t time.Duration) error {

	cmd := exec.Command("/bin/sh", "-c", c)

	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd")
		return err
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	//Start process. Exit code 127 if process fail to start.
	if err := cmd.Start(); err != nil {
		fmt.Fprintln(os.Stderr, "Error stating Cmd")
		return err
	}
	var timer *time.Timer
	if t > 0 {
		timer = time.NewTimer(t)
		var err error
		go func(timer *time.Timer, cmd *exec.Cmd) {
			for range timer.C {
				e := cmd.Process.Kill()
				if e != nil {
					err = errors.New("the command has timeout but the process could not be killed")
				} else {
					err = errors.New("the command timed out")
				}
			}
		}(timer, cmd)
	}

	err = cmd.Wait()

	if t > 0 {
		timer.Stop()
	}

	if err != nil {
		return errors.New("the command did not succeed")
	}

	return nil
}