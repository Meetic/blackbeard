package cmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/Meetic/blackbeard/pkg/resource"
	"github.com/gosuri/uiprogress"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	defaultTimeout = 5 * time.Minute
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply a given inventory to the associated namespace",
	Long: `This command will update the configuration files for the given namespace using the inventory file
and apply the changes to the Kubernetes namespace.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		err := runApply(namespace)
		if err != nil {
			logrus.Fatal(err)
		}

	},
}

//NewApplyCommand returns the apply cobra command
func NewApplyCommand() *cobra.Command {
	addCommonNamespaceCommandFlags(applyCmd)
	applyCmd.Flags().BoolVar(&wait, "wait", false, "wait until all pods are running")
	applyCmd.Flags().DurationVarP(&timeout, "timeout", "t", defaultTimeout, "The max time to wait for pods to be all running.")

	return applyCmd
}

func runApply(namespace string) error {

	if namespace == "" {
		return errors.New("you must specified a namespace using the --namespace flag")
	}

	files := newFileClient(dir)
	api := newAPI(files, newKubernetesClient())

	err := api.Apply(namespace, files.ConfigPath())
	if err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"namespace": namespace,
	}).Info("Playbook has been deployed")

	if wait {
		logrus.WithFields(logrus.Fields{
			"namespace": namespace,
		}).Info("Waiting for namespace to be ready...")

		if err := waitForNamespaceReady(api.Namespaces(), namespace); err != nil {
			return err
		}
	}

	return nil
}

func waitForNamespaceReady(ns resource.NamespaceService, namespace string) error {

	//Starting watching namespace
	go func() {
		err := ns.WatchNamespaces()
		if err != nil {
			logrus.Errorf("Error while trying to watch namespace : %s", err.Error())
		}
	}()

	//register for event listening
	ns.AddListener("cli_progress_bar")

	//initiate progress bar
	uiprogress.Start()
	bar := uiprogress.AddBar(100).AppendCompleted().PrependElapsed()

	//Init timer
	timerCh := time.NewTimer(timeout).C
	doneCh := make(chan bool)

	go func() {
		for {
			select {
			case e := <-ns.Events("cli_progress_bar"):
				if e.Type == resource.NamespaceStatusUpdate && e.Namespace == namespace {

					logrus.WithFields(logrus.Fields{
						"namespace": e.Namespace,
						"type":      e.Type,
						"Phase":     e.Phase,
						"Status":    e.Status,
					}).Debug("Event received for namespace")

					if err := bar.Set(e.Status); err != nil {
						logrus.Debugf("Error when incrementing progress bar : %s", err.Error())
					}

					if e.Status == 100 {
						doneCh <- true
					}
				}
			}
		}
	}()

	for {
		select {
		case <-timerCh:
			return fmt.Errorf("time out : Some pods are not yet ready")
		case <-doneCh:
			logrus.WithFields(logrus.Fields{
				"namespace": namespace,
			}).Info("Namespace is ready")
			return nil
		}
	}
}
