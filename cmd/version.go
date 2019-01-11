package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print blackbeard version",
	Long:  "This command will print blackbeard version.",

	Run: func(cmd *cobra.Command, args []string) {
		runVersion()
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}

func runVersion() {
	fmt.Println(fmt.Sprintf("blackbeard version %s", Version))
}
