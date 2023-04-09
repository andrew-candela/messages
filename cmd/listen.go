package cmd

import (
	"fmt"
	"os"

	"github.com/andrew-candela/messages/messages"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const PORT = "1053"

var port string

func listen(keyFile string, groupName string, recipConfig messages.Config, port string) {

	key, err := messages.ReadExistingKey(keyFile)
	if err != nil {
		fmt.Println("Could not read keyfile:", err)
		return
	}
	targets := messages.MakeTargets(recipConfig)
	c := make(chan messages.Packet, 10)
	go messages.Listen(PORT, c, *key, *targets)
	messages.PrintUDPOutput(c)
}

func init() {
	listenCommand.Flags().StringVarP(&port, "port", "p", PORT, "Port number to listen on")
	rootCmd.AddCommand(listenCommand)
}

var listenCommand = &cobra.Command{
	Use: "listen",
	Run: func(cmd *cobra.Command, args []string) {
		var recipConfig messages.Config
		viper.GetViper().ReadInConfig()
		keyFile := viper.GetString("private_key_file")
		err := viper.UnmarshalKey(group, &recipConfig)
		if err != nil {
			fmt.Print("error:", err)
			os.Exit(1)
		}
		listen(keyFile, group, recipConfig, port)
	},
}
