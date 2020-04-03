package cmd

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [command]",
	Short: "Delete an object",
	Long: `Delete resources by namespace or names.

Deletetion of a namespace will delete the namespace and remove all his attached object including the intentory attached to it. While removing an object will only supress it form the namespace but keep everything else.`,
	Run: func(cmd *cobra.Command, args []string) {
		runDelete()
	},
}

func NewDeleteCommand() *cobra.Command {
	deleteCmd.AddCommand(NewDeleteJobCommand())
	deleteCmd.AddCommand(NewDeleteNamespaceCommand())

	return deleteCmd
}

func runDelete() {
	tpl := template.Must(template.New("deleteCmd").Parse(`
Using the get command without any sub-command makes no sens. Please use one of the following sub-command :
{{range . -}}
- {{.}}
{{end -}}
`))

	data := []string{"delete namespace", "delete job"}

	contents := bytes.Buffer{}
	if err := tpl.Execute(&contents, data); err != nil {
		logrus.Fatalf("error while executing template : %v", err)
	}

	fmt.Println(contents.String())
}
