package cmd

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Show informations about a given namespace.",
	Long: `This command display informations from a given namespace such as the list of exposed services
or the url where you can join services throw ingress.`,
	Run: func(cmd *cobra.Command, args []string) {
		runGet()
	},
}

func NewGetCommand() *cobra.Command {
	addCommonNamespaceCommandFlags(getCmd)

	getCmd.AddCommand(NewGetNamespacesCommand())
	getCmd.AddCommand(NewGetServicesCommand())

	return getCmd
}

func runGet() {
	tpl := template.Must(template.New("getCmd").Parse(`
Using the get command without any sub-command makes no sens. Please use one of the following sub-command :
{{range .}}
- {{.}}
{{end}}
`))

	data := []string{"get services", "get namespaces"}

	contents := bytes.Buffer{}
	if err := tpl.Execute(&contents, data); err != nil {
		logrus.Fatalf("error when executing template : %v", err)
	}

	fmt.Println(contents.String())

}
