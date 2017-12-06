package cmd

import (
	"github.com/Meetic/blackbeard/pkg/files"
	"github.com/Meetic/blackbeard/pkg/kubecli"

	"github.com/Meetic/blackbeard/pkg/http"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Launch the blackbeard server",
	Long: `This command run a web server that expose a REST API.
This API let the client use all the features provided by Blackbeard such as create a namespace and apply a change in a inventory.`,
	Run: func(cmd *cobra.Command, args []string) {
		runServe(inventoryPath)
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringVarP(&inventoryPath, "inventory-dir", "d", "", "The directory where to store the inventory config")

}

func runServe(inventoryPath string) {

	files := files.NewClient(templatePath, configPath, inventoryPath, defaultsPath)

	cli := kubecli.NewClient(configPath)

	h := http.NewHandler(files, cli)
	s := http.NewServer(h)
	s.Serve()

}
