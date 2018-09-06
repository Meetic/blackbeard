package cmd

import (
	"errors"
	"log"

	"fmt"

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

	if !askForConfirmation(fmt.Sprintf("You are about to delete the inventory %s and all its associated files. Are you sure?", namespace)) {
	    return nil
	}

	api := newAPI(newFileClient(dir), newKubernetesClient())

	err := api.Delete(namespace)
	if err != nil {
		return err
	}

	fmt.Printf("Namespace %s has been sucessfully deleted.\n", namespace)

	return nil
}
