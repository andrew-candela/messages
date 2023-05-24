package cmd

import (
	"fmt"
	"os"

	"github.com/andrew-candela/messages/messages"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var webserver_host string

func getMyIP(keyFile string, host string) {
	key, err := messages.ReadExistingKey(keyFile)
	if err != nil {
		fmt.Println("Could not read private key file:", keyFile)
		os.Exit(1)
	}
	client := messages.MakeClient()
	my_ip := messages.GetMyIp(host, key, client)
	fmt.Println(my_ip)
}

var IPCommand = &cobra.Command{
	Use: "get-ip",
	Run: func(cmd *cobra.Command, args []string) {
		viper.ReadInConfig()
		keyFile := viper.GetString("private_key_file")
		getMyIP(keyFile, webserver_host)
	},
}

func getConfig(keyFile string, host string) {
	key, err := messages.ReadExistingKey(keyFile)
	if err != nil {
		fmt.Println("Could not read private key file:", keyFile)
		os.Exit(1)
	}
	client := messages.MakeClient()
	config := messages.GetConfig(host, key, client)
	fmt.Println(config)
}

var webClientConfigCommand = &cobra.Command{
	Use: "get-config",
	Run: func(cmd *cobra.Command, args []string) {
		viper.ReadInConfig()
		keyFile := viper.GetString("private_key_file")
		getConfig(keyFile, webserver_host)
	},
}

func publish(keyFile string, host string) {
	key, err := messages.ReadExistingKey(keyFile)
	if err != nil {
		fmt.Println("Could not read private key file:", keyFile)
		os.Exit(1)
	}
	client := messages.MakeClient()

	config := messages.Publish(host, key, client)
	fmt.Println(config)
}

var webClientSubscribeCommand = &cobra.Command{
	Use: "get-config",
	Run: func(cmd *cobra.Command, args []string) {
		viper.ReadInConfig()
		keyFile := viper.GetString("private_key_file")
		getConfig(keyFile, webserver_host)
	},
}

func init() {
	IPCommand.Flags().StringVar(&webserver_host, "host", "localhost:8080", "The hostname and port of the webserver. Defaults to 'localhost:8080'.")
	webClientConfigCommand.Flags().StringVar(&webserver_host, "host", "localhost:8080", "The hostname and port of the webserver. Defaults to 'localhost:8080'.")
	rootCmd.AddCommand(IPCommand)
	rootCmd.AddCommand(webClientConfigCommand)
}
