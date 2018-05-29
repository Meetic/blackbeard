package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"bufio"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	templateDir  = "templates"
	configDir    = "configs"
	inventoryDir = "inventories"
	defaultFile  = "defaults.json"
)

var cfgFile string
var dir string
var templatePath string
var configPath string
var inventoryPath string
var defaultsPath string
var namespace string
var cors bool
var kubeConfigPath string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
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

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.blackbeard.yaml)")
	RootCmd.PersistentFlags().StringVar(&dir, "dir", "", "Use the specified dir as root path to execute commands. Default is the current dir.")

	viper.BindPFlag("dir", RootCmd.PersistentFlags().Lookup("dir"))
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

	//Define current working dir as default value
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Error when getting the working dir : ", err)
	}
	viper.SetDefault("working-dir", currentDir)

	var wd string

	//dir is the value of the flag --dir
	if dir != "" {
		wd = dir
	} else {
		wd = viper.GetString("working-dir")
	}

	if ok, _ := fileExists(wd); ok != true {
		log.Fatalf("Your specified working dir does not exit : %s", wd)
	}

	templatePath = filepath.Join(wd, templateDir)
	configPath = filepath.Join(wd, configDir)
	inventoryPath = filepath.Join(wd, inventoryDir)
	defaultsPath = filepath.Join(wd, defaultFile)
	kubeConfigPath = filepath.Join(homeDir(), ".kube", "config")

	checkWorkingDir()
}

func checkWorkingDir() {

	if ok, _ := fileExists(templatePath); ok != true {
		log.Fatalf("A playbook must contains a `%s` dir. No one has been found.\nPlease check the playbook or change the working directory using the --dir option.", templateDir)
	}

	if ok, _ := fileExists(defaultsPath); ok != true {
		log.Fatalf("Your working directory must contains a `%s` file.\nPlease check the playbook or change the working directory using the --dir option.", defaultFile)
	}

	if ok, _ := fileExists(configPath); ok != true {
		if err := os.Mkdir(configPath, 0755); err != nil {
			log.Fatalf("Impossible to create the %s directory. Please check directory rights.", configDir)
		}
	}

	if ok, _ := fileExists(inventoryPath); ok != true {
		if err := os.Mkdir(inventoryPath, 0755); err != nil {
			log.Fatalf("Impossible to create the %s directory. Please check directory rights.", inventoryDir)
		}
	}
}

func fileExists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return true, err
	}

	return true, nil
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func askForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}
