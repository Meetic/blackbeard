// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	Short: "Launch the blackbeard webserver",
	Long:  ``,
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
