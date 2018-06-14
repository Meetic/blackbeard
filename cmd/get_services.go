package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var getServicesCmd = &cobra.Command{
	Use:   "services",
	Short: "Show informations about exposed services from a given namespace.",
	Long: `This command display informations from a given namespace such as the list of exposed services
or the url where you can join services throw ingress.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := runGetServices()
		if err != nil {
			log.Fatal(err.Error())
		}

	},
}

func init() {
	getCmd.AddCommand(getServicesCmd)
	getServicesCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "The namespace from which get info")
}

func runGetServices() error {

	if namespace == "" {
		return errors.New("you must specified a namespace using the --namespace flag")
	}

	api := newAPI()

	//get exposed services (NodePort)
	services, err := api.ListExposedServices(namespace)
	if err != nil {
		return errors.New(fmt.Sprintf("an error occurend when getting information about services : %v", err))
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	fmt.Fprintln(w, "Service Name\tAddress\tPort\tExposed Port\t")
	for _, svc := range services {
		for _, p := range svc.Ports {
			fmt.Fprintf(w, fmt.Sprintf("%s\t%s\t%d\t%d\t\n", svc.Name, svc.Addr, p.Port, p.ExposedPort))
		}
	}
	fmt.Fprintln(w)
	w.Flush()

	return nil

}
