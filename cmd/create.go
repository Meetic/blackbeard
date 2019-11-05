package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"path/filepath"

	"github.com/Meetic/blackbeard/pkg/playbook"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a namespace and generated a dedicated inventory.",
	Long: `This command will generate an inventory file called {{namespace}}_inventory.json

This file contains all the parameters needed to build a complete Kubernetes configuration.
Feel free to edit this file before applying changes.
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := runCreate(namespace)
		if err != nil {
			logrus.Fatal(err.Error())
		}
	},
}

func NewCreateCommand() *cobra.Command {
	addCommonNamespaceCommandFlags(createCmd)
	return createCmd
}

func runCreate(namespace string) error {

	if namespace == "" {
		return errors.New("you must specified a namespace using the --namespace flag")
	}

	files := newFileClient(playbookDir)

	api := newAPI(files, newKubernetesClient())

	inv, err := api.Create(namespace)
	if err != nil {
		return err
	}

	tpl := template.Must(template.New("config").Parse(`Namespace for user {{.Inv.Namespace}} has been created !

	A inventory file has been generated : {{.File}}
	Feel free to edit this file to match your desired testing env configuration.
	`))

	message := bytes.Buffer{}
	if err := tpl.Execute(&message, struct {
		File string
		Inv  playbook.Inventory
	}{
		File: filepath.Join(files.InventoryPath(), inv.Namespace+"_inventory.json"),
		Inv:  inv,
	}); err != nil {
		return err
	}

	fmt.Println(message.String())

	return nil
}
