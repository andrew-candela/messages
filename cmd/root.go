package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	user, mode, group, listen_port string
)

const (
	UDP_MODE  = "udp"
	HTTP_MODE = "http"
)

var rootCmd = &cobra.Command{
	Use:   "messenger",
	Short: "messenger: send messages directly to your friends.",
	Long:  "messenger: Entrypoint to the messaging tool that sends messages to your pals.",
}

func Execute() {
	rootCmd.Execute()
}
func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&group, "group", "g", "", "Group Name to listen or write to")
	rootCmd.PersistentFlags().StringVarP(&user, "user", "u", os.Getenv("USER"), "Your display name in messages")
}

func initConfig() {
	viper.SetConfigName("config")      // name of config file (without extension)
	viper.SetConfigType("toml")        // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/udpm")   // path to look for the config file in
	viper.AddConfigPath("$HOME/.udpm") // call multiple times to add many search paths
	viper.AddConfigPath(".")           // optionally look for config in the working directory
}
