package cmd

import (
	"errors"
	"fmt"
	"log"

	"github.com/Meetic/blackbeard/pkg/files"
	"github.com/Meetic/blackbeard/pkg/kubectl"

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

	f := files.NewClient(templatePath, configPath, inventoryPath, defaultsPath)
	cli := kubectl.NewClient(configPath)

	//Reset inventory file
	if err := f.InventoryService().Reset(namespace); err != nil {
		return err
	}

	//Get inventory
	inv, err := f.InventoryService().Get(namespace)
	if err != nil {
		return err
	}

	//Apply inventory to configuration
	if err := f.ConfigService().Apply(inv); err != nil {
		return err
	}

	//Apply changes to Kubernetes
	if err = cli.NamespaceConfigurationService().Apply(inv.Namespace); err != nil {
		return err
	}

	fmt.Println("Namespace has been reset successfully")

	return nil
}
