package cmd

import (
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
	serveCmd.Flags().BoolVar(&cors, "cors", false, "Enable cors")
}

func runServe() {

	api := newAPI()

	wh := websocket.NewHandler(api.Namespaces(), api.Inventories())

	h := http.NewHandler(api, wh, configPath, cors)
	s := http.NewServer(h)
	s.Serve()

}
