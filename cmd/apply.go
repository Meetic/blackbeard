package cmd

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gosuri/uiprogress"
	"github.com/spf13/cobra"
)

const (
	defaultTimeout = 5 * time.Minute
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply a given inventory to the associated namespace",
	Long: `This command will update the configuration files for the given namespace using the inventory file
and apply the changes to the Kubernetes namespace.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		err := runApply(namespace)
		if err != nil {
			log.Fatal(err)
		}

	},
}

func NewApplyCommand() *cobra.Command {
	addCommonNamespaceCommandFlags(applyCmd)
	applyCmd.Flags().BoolVar(&wait, "wait", false, "wait until all pods are running")
	applyCmd.Flags().DurationVarP(&timeout, "timeout", "t", defaultTimeout, "The max time to wait for pods to be all running.")

	return applyCmd
}

func runApply(namespace string) error {

	if namespace == "" {
		return errors.New("you must specified a namespace using the --namespace flag")
	}

	files := newFileClient(dir)
	api := newAPI(files, newKubernetesClient())

	err := api.Apply(namespace, files.ConfigPath())
	if err != nil {
		return err
	}

	fmt.Println("Playbook has been deployed.")

	if wait {

		fmt.Println("Waiting for namespace to be ready...")
		//init progress bar
		uiprogress.Start()
		bar := uiprogress.AddBar(100).AppendCompleted().PrependElapsed()

		if err := api.WaitForNamespaceReady(namespace, timeout, bar); err != nil {
			return err
		}

		fmt.Println("Namespace is ready.")

	}

	return nil
}
