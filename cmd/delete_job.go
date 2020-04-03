package cmd

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var deleteJobCmd = &cobra.Command{
	Use:   "job [NAME]",
	Short: "Delete a job object from a namespace",
	Long: `Delete a job object that started a pod.

It won't remove any configuration in the inventory. Reapplying the inventory will redeploy it.

Kubernetes will also remove the pod whatever the status.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := runDeleteJob(args[0])
		if err != nil {
			logrus.Fatal(err.Error())
		}
	},
}

func NewDeleteJobCommand() *cobra.Command {
	addCommonNamespaceCommandFlags(deleteJobCmd)
	return deleteJobCmd
}

func runDeleteJob(resource string) error {
	if namespace == "" {
		// should set the namespace to default namespace value set in the kube/config
		return errors.New("you must specified a namespace using the --namespace flag")
	}

	api := newAPI(newFileClient(dir), newKubernetesClient())
	err := api.DeleteResource(namespace, resource)
	if err != nil {
		return errors.New(fmt.Sprintf("an error occurend when removing the job : %v", err))
	}

	logrus.WithFields(logrus.Fields{
		"namespace": namespace,
		"job":       resource,
	}).Info("job deleted")

	return nil
}
