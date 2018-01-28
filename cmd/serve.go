package cmd

import (
	"github.com/Meetic/blackbeard/pkg/files"
	"github.com/Meetic/blackbeard/pkg/kubectl"
	"github.com/Meetic/blackbeard/pkg/kubernetes"
	"github.com/Meetic/blackbeard/pkg/websocket"

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
		runServe()
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)
}

func runServe() {

	f := files.NewClient(templatePath, configPath, inventoryPath, defaultsPath)
	cli := kubectl.NewClient(configPath)
	kube := kubernetes.NewClient(kubeConfigPath)
	ws := websocket.NewHandler(kube, f)

	h := http.NewHandler(f, cli, kube, ws)
	s := http.NewServer(h)
	s.Serve()

}
