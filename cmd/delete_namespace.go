package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// deleteCmd represents the create command
var deleteNamespaceCmd = &cobra.Command{
	Use:   "namespace [NAME]",
	Short: "Delete a namespace",
	Long:  `This command delete a namespace and all the associated resources`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := runDeleteNamespace(args[0])
		if err != nil {
			logrus.Fatal(err.Error())
		}
	},
}

func NewDeleteNamespaceCommand() *cobra.Command {
	return deleteNamespaceCmd
}

func runDeleteNamespace(namespace string) error {
	if !askForConfirmation(fmt.Sprintf("You are about to delete the inventory %s and all its associated files. Are you sure?", namespace), os.Stdin) {
		return nil
	}

	api := newAPI(newFileClient(playbookDir), newKubernetesClient())

	err := api.Delete(namespace, false)
	if err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"namespace": namespace,
	}).Info("namespace deleted")

	return nil
}
