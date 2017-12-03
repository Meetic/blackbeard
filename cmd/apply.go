package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/Meetic/blackbeard/pkg/files"
	"github.com/Meetic/blackbeard/pkg/kubecli"

	"github.com/spf13/cobra"
)

var namespace string

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply a given inventory to the testing env targeted by the inventory.",
	Long: `Apply a given inventory to the testing env targeted by the inventory.
This command will update the Kubernetes files and apply this configuration to the testing env.
	`,
	Run: func(cmd *cobra.Command, args []string) {

		err := runApply(namespace)

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(-1)
		}

	},
}

func init() {
	RootCmd.AddCommand(applyCmd)
	applyCmd.Flags().StringVarP(&namespace, "username", "u", "", "The username for the environment")
	applyCmd.Flags().StringVarP(&usr, "namespace", "n", "", "Same as username")

}

func runApply(namespace string) error {

	if namespace == "" {
		return errors.New("You must specified a username")
	}

	files := files.NewClient(templatePath, configPath, inventoryPath, defaultsPath)

	inv, err := files.InventoryService().Get(namespace)

	if err != nil {
		return err
	}

	err = files.ConfigService().Apply(inv)
	if err != nil {
		return err
	}

	cli := kubecli.NewClient(configPath)

	err = cli.NamespaceService().Apply(inv)
	if err != nil {
		return err
	}

	fmt.Println("Playbook has been deployed!")

	return nil
}
