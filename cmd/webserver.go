package cmd

import (
	"fmt"
	"os"

	"github.com/andrew-candela/messages/messages"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var webservice_address string

func startWebserver(recipConfig messages.Config, port string) {
	conf := messages.MakeConfigMap(recipConfig)
	messages.StartServer(webservice_address, conf)
}

var startWebserverCommand = &cobra.Command{
	Use: "webserver",
	Run: func(cmd *cobra.Command, args []string) {
		var recipConfig messages.Config
		viper.GetViper().ReadInConfig()
		err := viper.UnmarshalKey(group, &recipConfig)
		if err != nil {
			fmt.Print("error:", err)
			os.Exit(1)
		}
		startWebserver(recipConfig, listen_port)
	},
}

func init() {
	startWebserverCommand.Flags().StringVarP(&webservice_address, "address", "a", "localhost:8080", "The address that the webserver or messanger will listen on. Defaults to 'localhost:8080'.")
	rootCmd.AddCommand(startWebserverCommand)
}
