package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"

	"github.com/Meetic/blackbeard/pkg/blackbeard"
	"github.com/Meetic/blackbeard/pkg/files"
	"github.com/Meetic/blackbeard/pkg/kubecli"

	"github.com/spf13/cobra"
)

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
		err := runCreate(namespace)
		if err != nil {
			log.Fatal(err.Error())
		}
	},
}

func init() {
	RootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "The namespace where to apply configuration")
}

func runCreate(namespace string) error {

	if namespace == "" {
		return errors.New("you must specified a namespace for the testing env using the --namespace flag")
	}

	f := files.NewClient(templatePath, configPath, inventoryPath, defaultsPath)

	inv, err := f.InventoryService().Create(namespace)

	if err != nil {
		log.Println(err.Error())
		log.Println("continue")
	}

	err = f.ConfigService().Apply(inv)
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
