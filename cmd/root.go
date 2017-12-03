package cmd

import (
	"fmt"
	"log"
	"os"

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

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "blackbeard",
	Short: "blackbeard is a tool that let you create a manage any testing environment using Kubernetes.",
	Long: `blackbeard is a tool that let you create a manage any testing environment using Kubernetes.

Kugen let you first create en Kubernetes namespace using the "create" command.
This command wil generate a inventory file containning the default config for the testing env you are creating.
Feel free to update this inventory file.

Then blackbeard configure your testing env using a auto-generated Kubernetes config using the specified inventory file.
This action can be done using the "apply" command.
	`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
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

	fmt.Println(dir)

	if dir != "" {
		wd = dir
	} else {
		wd = viper.GetString("working-dir")
	}

	if ok, _ := fileExists(wd); ok != true {
		panic(fmt.Sprintf("Your specified working dir does not exit : %s", wd))
	}

	templatePath = fmt.Sprintf("%s/%s/", wd, templateDir)
	configPath = fmt.Sprintf("%s/%s/", wd, configDir)
	inventoryPath = fmt.Sprintf("%s/%s/", wd, inventoryDir)
	defaultsPath = fmt.Sprintf("%s/%s", wd, defaultFile)

	checkWorkingDir()
}

func checkWorkingDir() {
	if ok, _ := fileExists(templatePath); ok != true {
		panic(fmt.Sprintf("Your working directory must contains a `%s` dir.", templateDir))
	}

	if ok, _ := fileExists(configPath); ok != true {
		panic(fmt.Sprintf("Your working directory must contains a `%s` dir.", configDir))
	}

	if ok, _ := fileExists(inventoryPath); ok != true {
		panic(fmt.Sprintf("Your working directory must contains a `%s` dir.", inventoryDir))
	}

	if ok, _ := fileExists(defaultsPath); ok != true {
		panic(fmt.Sprintf("Your working directory must contains a `%s` file.", defaultFile))
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
