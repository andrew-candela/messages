package cmd

import (
	"fmt"
	"os"

	"github.com/andrew-candela/messages/messages"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func produce(keyFile string, group_name string, recip_config messages.Config) {

	if mode == HTTP_MODE {
		fmt.Println("HTTP mode is not supported yet")
		return
	}
	fmt.Println(recip_config)
	targets := messages.MakeTargets(recip_config)
	key, err := messages.ReadExistingKey(keyFile)
	if err != nil {
		fmt.Println("Could not read private key file:", keyFile)
		os.Exit(1)
	}
	messages.ProduceMessages(*targets, user, key)
}

func init() {
	rootCmd.AddCommand(produceCommand)
}

var produceCommand = &cobra.Command{
	Use: "produce",
	Run: func(cmd *cobra.Command, args []string) {
		var recipConfig messages.Config
		viper.ReadInConfig()
		keyFile := viper.GetString("private_key_file")
		err := viper.UnmarshalKey(group, &recipConfig)
		if err != nil {
			fmt.Print("error:", err)
			os.Exit(1)
		}
		produce(keyFile, group, recipConfig)
	},
}
