package cmd

import (
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// killCmd represents the kill command
var killCmd = &cobra.Command{
	Use:   "kill",
	Short: "kill a pod",
	Long:  `This command find the corresponding pod and kill it`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		errs := runKill(namespace, args)
		if len(errs) > 0 {
			for _, err := range errs {
				logrus.Error(err.Error())
			}
			logrus.Fatal("some pod could not be killed")
		}
	},
}

func NewKillCommand() *cobra.Command {
	addCommonNamespaceCommandFlags(killCmd)

	return killCmd
}

func runKill(namespace string, deployments []string) []error {
	if namespace == "" {
		return []error{errors.New("you must specified a namespace using the --namespace flag")}
	}

	api := newAPI(newFileClient(dir), newKubernetesClient())

	errs := api.Kill(namespace, deployments)
	if len(errs) > 0 {
		return errs
	}

	logrus.WithFields(logrus.Fields{
		"namespace":   namespace,
		"deployments": deployments,
	}).Info("pod are being killed")

	return nil
}
