package cmd

import (
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var getNamespacesCmd = &cobra.Command{
	Use:   "namespaces",
	Short: "Show informations about kubernetes namespaces.",
	Long: `Show informations about kubernetes namespaces such as names, status (percentage of pods in a running status),
managed or not with the current playbook, etc.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := runGetNamespaces()
		if err != nil {
			logrus.Fatal(err.Error())
		}

	},
}

func NewGetNamespacesCommand() *cobra.Command {
	return getNamespacesCmd
}

func runGetNamespaces() error {

	api := newAPI(newFileClient(dir), newKubernetesClient())

	namespaces, err := api.ListNamespaces()
	if err != nil {
		return errors.New(fmt.Sprintf("an error occurend when getting information about namespaces : %v", err))
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	fmt.Fprintln(w, "Namespace\tPhase\tStatus\tManaged\t")
	for _, namespace := range namespaces {
		fmt.Fprint(w, fmt.Sprintf("%s\t%s\t%d%%\t%t\t\n", namespace.Name, namespace.Phase, namespace.Status, namespace.Managed))
	}
	fmt.Fprintln(w)
	w.Flush()

	return nil

}
