package cmd

import (
	"errors"
	"log"

	"fmt"

	"github.com/Meetic/blackbeard/pkg/files"
	"github.com/Meetic/blackbeard/pkg/kubernetes"
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
			log.Fatal(err.Error())
		}
	},
}

func init() {
	RootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "The namespace where to apply configuration")
}

func runDelete(namespace string) error {
	if namespace == "" {
		return errors.New("you must specified a namespace using the --namespace flag")
	}

	askForConfirmation(fmt.Sprintf("You are about to delete the inventory %s and all its associated files. Are you sure?", namespace))

	kube := kubernetes.NewClient(kubeConfigPath)
	if err := kube.NamespaceService().Delete(namespace); err != nil {
		return err
	}

	f := files.NewClient(templatePath, configPath, inventoryPath, defaultsPath)

	if err := f.InventoryService().Delete(namespace); err != nil {
		return err
	}

	if err := f.ConfigService().Delete(namespace); err != nil {
		return err
	}

	fmt.Printf("Namespace %s has been sucessfully deleted.\n", namespace)

	return nil
}
