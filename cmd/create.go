package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/Meetic/blackbeard/pkg/blackbeard"
	"github.com/Meetic/blackbeard/pkg/files"
	"github.com/Meetic/blackbeard/pkg/kubecli"

	"github.com/spf13/cobra"
)

var usr string

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a testing environnement for a given user.",
	Long: `Create a testing environnement for a given user.
This command will generate an inventory file calle {{username}}_inventory.json
This file contains all the necessary parameters used to build a comple Kubernetes configuration.
Feel free to edit this file before applying changes.
`,

	Run: func(cmd *cobra.Command, args []string) {
		err := runCreate(usr, inventoryPath)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(-1)
		}
	},
}

func init() {
	RootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&usr, "username", "u", "", "The username for the environement")
	createCmd.Flags().StringVarP(&usr, "namespace", "n", "", "Same as username")
}

func runCreate(usr, inventoryPath string) error {

	if usr == "" {
		return errors.New("You must specified a username for the testing env using the --username flag")
	}

	files := files.NewClient(templatePath, configPath, inventoryPath, defaultsPath)

	inv, err := files.InventoryService().Create(usr)

	if err != nil {
		log.Println(err.Error())
		log.Println("continue")
	}

	err = files.ConfigService().Apply(inv)
	if err != nil {
		return err
	}

	cli := kubecli.NewClient(configPath)

	err = cli.NamespaceService().Create(inv)

	if err != nil {
		return err
	}

	tpl := template.Must(template.New("config").Parse(`Recette env for user {{.Inv.Namespace}} has been created !

	A inventory file has been generated : {{.File}}
	Feel free to edit this file to match your desired testing env configuration.
	`))

	message := bytes.Buffer{}
	if err := tpl.Execute(&message, struct {
		File string
		Inv  blackbeard.Inventory
	}{
		File: inventoryPath + "/" + inv.Namespace + "_inventory.json",
		Inv:  inv,
	}); err != nil {
		return err
	}

	fmt.Println(message.String())

	return nil
}
