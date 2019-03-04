package cmd

import (
	"errors"
	"os"

	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// deleteCmd represents the create command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a namespace",
	Long:  `This command delete a namespace and all the associated resources`,

	Run: func(cmd *cobra.Command, args []string) {
		err := runDelete(namespace)
		if err != nil {
			logrus.Fatal(err.Error())
		}
	},
}

func NewDeleteCommand() *cobra.Command {
	addCommonNamespaceCommandFlags(deleteCmd)

	return deleteCmd
}

func runDelete(namespace string) error {
	if namespace == "" {
		return errors.New("you must specified a namespace using the --namespace flag")
	}

	if !askForConfirmation(fmt.Sprintf("You are about to delete the inventory %s and all its associated files. Are you sure?", namespace), os.Stdin) {
		return nil
	}

	api := newAPI(newFileClient(dir), newKubernetesClient())

	err := api.Delete(namespace, false)
	if err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"namespace": namespace,
	}).Info("namespace deleted")

	return nil
}
