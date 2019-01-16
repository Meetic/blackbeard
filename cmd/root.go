package cmd

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/Meetic/blackbeard/pkg/api"
	"github.com/Meetic/blackbeard/pkg/files"
	"github.com/Meetic/blackbeard/pkg/kubernetes"
	"github.com/sirupsen/logrus"
)

var (
	cfgFile           string
	dir               string
	kubectlConfigPath string
	v                 string
	namespace         string
	cors              bool
	wait              bool
	timeout           time.Duration
	port              int
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "blackbeard",
	Short: "Blackbeard is a tool that let you create and manage multiple version of the same stack using Kubernetes and namespace",
	Long: `Blackbeard let you apply a bunch of configuration files template into different namespaces using some provided values.

Blackbeard is made to be executed using a directory containing configuration files and directories called a Playbook.

Using blackbeard and a Playbook, you can easily create a namespace by using the "create" command.
This command will generate an inventory file containing the default configuration for the namespace you are creating.

Feel free to update this inventory file manually.

Then Blackbeard configure your namespace using a auto-generated Kubernetes config using the specified inventory file.
This action can be done using the "apply" command.
	`,
}

func NewBlackbeardCommand() *cobra.Command {

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if err := setUpLogs(os.Stdout, v); err != nil {
			return err
		}
		return nil
	}

	initConfig()

	rootCmd.AddCommand(NewServeCommand())
	rootCmd.AddCommand(NewApplyCommand())
	rootCmd.AddCommand(NewCreateCommand())
	rootCmd.AddCommand(NewDeleteCommand())
	rootCmd.AddCommand(NewGetCommand())
	rootCmd.AddCommand(NewResetCommand())

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.blackbeard.yaml)")
	rootCmd.PersistentFlags().StringVar(&dir, "dir", "", "Use the specified dir as root path to execute commands. Default is the current dir.")
	rootCmd.PersistentFlags().StringVar(&kubectlConfigPath, "kube-config-path", kubernetes.KubeConfigDefaultPath(), "kubectl config file")
	rootCmd.PersistentFlags().StringVarP(&v, "verbosity", "v", logrus.WarnLevel.String(), "Log level (debug, info, warn, error, fatal, panic")

	viper.BindPFlag("dir", rootCmd.PersistentFlags().Lookup("dir"))

	return rootCmd

}

func addCommonNamespaceCommandFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&namespace, "namespace", "n", "", "The namespace where to apply configuration")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".blackbeard") // name of config file (without extension)
	viper.AddConfigPath("$HOME")       // adding home directory as first search path
	viper.AutomaticEnv()               // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	if dir == "" {
		//Define current working dir as default value
		currentDir, err := os.Getwd()
		if err != nil {
			log.Fatal("Error when getting the working dir : ", err)
		}
		viper.SetDefault("working-dir", currentDir)
		dir = viper.GetString("working-dir")
	}

}

func askForConfirmation(message string, reader io.Reader) bool {

	r := bufio.NewReader(reader)

	for {
		fmt.Printf("%s [y/n]: ", message)

		response, err := r.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else {
			return false
		}
	}
}

func newKubernetesClient() *kubernetes.Client {
	kube, err := kubernetes.NewClient(kubectlConfigPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	return kube
}

func newFileClient(dir string) *files.Client {
	f, err := files.NewClient(dir)
	if err != nil {
		log.Fatal(err.Error())
	}

	return f

}

func newAPI(files *files.Client, kube *kubernetes.Client) api.Api {
	return api.NewApi(
		files.Inventories(),
		files.Configs(),
		files.Playbooks(),
		kube.Namespaces(),
		kube.Pods(),
		kube.Services(),
		kube.Cluster(),
	)
}

func setUpLogs(out io.Writer, level string) error {

	logrus.SetOutput(out)
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}
	logrus.SetLevel(lvl)
	return nil
}
