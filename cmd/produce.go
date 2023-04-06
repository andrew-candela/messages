package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/andrew-candela/messages/messages"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func produce(keyFile string, group_name string, recip_config messages.Config) {
	var host string
	var user string
	flag.StringVar(&host, "host", "10.0.0.186:1053", "Provide the HOST:PORT to write to")
	flag.StringVar(&user, "user", os.Getenv("USER"), "What do you call yourself?")
	flag.Parse()
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
		group := args[0]
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
