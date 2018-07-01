package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var getNamespaceCmd = &cobra.Command{
	Use:   "namespaces",
	Short: "todo",
	Long:  `todo`,
	Run: func(cmd *cobra.Command, args []string) {
		err := runGetNamespaces()
		if err != nil {
			log.Fatal(err.Error())
		}

	},
}

func init() {
	getCmd.AddCommand(getNamespaceCmd)
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
