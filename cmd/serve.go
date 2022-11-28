package cmd

import (
	"github.com/spf13/cobra"

	"github.com/Meetic/blackbeard/pkg/http"
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

func NewServeCommand() *cobra.Command {
	serveCmd.Flags().BoolVar(&cors, "cors", false, "Enable cors")
	serveCmd.Flags().IntVar(&port, "port", 8080, "Use a specific port")

	return serveCmd
}

func runServe() {
	files := newFileClient(playbookDir)

	api := newAPI(files, newKubernetesClient())

	go api.WatchNamespaceDeleted()

	h := http.NewHandler(api, files.ConfigPath(), cors)
	s := http.NewServer(h)

	// start http web server
	s.Serve(port)
}
