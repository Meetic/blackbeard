package cmd

import (
	"github.com/Meetic/blackbeard/pkg/websocket"

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

func init() {
	RootCmd.AddCommand(serveCmd)
	serveCmd.Flags().BoolVar(&cors, "cors", false, "Enable cors")
	serveCmd.Flags().IntVar(&port, "port", 8080, "Use a specific port")
}

func runServe() {

	files := newFileClient(dir)

	api := newAPI(files, newKubernetesClient())

	wh := websocket.NewHandler(api)

	h := http.NewHandler(api, wh, files.ConfigPath(), cors)
	s := http.NewServer(h)

	// start http web server
	s.Serve(port)
}
