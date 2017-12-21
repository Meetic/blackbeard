package cmd

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/Meetic/blackbeard/pkg/kubernetes"
	"github.com/spf13/cobra"
)

var getServicesCmd = &cobra.Command{
	Use:   "services",
	Short: "Show informations about exposed services from a given namespace.",
	Long: `This command display informations from a given namespace such as the list of exposed services
or the url where you can join services throw ingress.`,
	Run: func(cmd *cobra.Command, args []string) {
		runGetServices()
	},
}

func init() {
	getCmd.AddCommand(getServicesCmd)
	getServicesCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "The namespace from which get info")
}

func runGetServices() {
	kube := kubernetes.NewClient(kubeConfigPath)

	//get exposed services (NodePort)
	services, err := kube.ResourceService().GetExposedServices(namespace)
	if err != nil {
		log.Fatalf("an error occurend when getting information about services : %v", err)
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

}
