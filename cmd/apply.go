package cmd

import (
	"errors"
	"fmt"
	"log"

	"github.com/spf13/cobra"
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
			log.Fatal(err)
		}

	},
}

func init() {
	RootCmd.AddCommand(applyCmd)
	applyCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "The namespace where to apply configuration")

}

func runApply(namespace string) error {

	if namespace == "" {
		return errors.New("you must specified a namespace using the --namespace flag")
	}

	api := newAPI()

	err := api.Apply(namespace, configPath)
	if err != nil {
		return err
	}

	fmt.Println("Playbook has been deployed!")

	return nil
}
