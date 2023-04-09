package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configFile string
	group      string
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
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is 'sample_config.toml')")
	rootCmd.PersistentFlags().StringVarP(&group, "group", "g", "", "Group Name to listen or write to")
}

func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigFile("sample_config.toml")
	}
}
