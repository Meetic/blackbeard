package cmd

import (
	"errors"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// resetCmd represents the reset command
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset a namespace based on the template files and the default inventory.",
	Long:  `This command will override the inventory and the config files for the given namespace and apply the changes into Kubernetes.`,

	Run: func(cmd *cobra.Command, args []string) {
		err := runReset(namespace)
		if err != nil {
			log.Fatal(err.Error())
		}
	},
}

func init() {
	RootCmd.AddCommand(resetCmd)
	resetCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "The namespace where to apply configuration")
}

func runReset(namespace string) error {

	if namespace == "" {
		return errors.New("you must specified a namespace using the --namespace flag")
	}

	api := newAPI()

	//Reset inventory file
	err := api.Reset(namespace, configPath)
	if err != nil {
		return err
	}

	fmt.Println("Namespace has been reset successfully")

	return nil
}
