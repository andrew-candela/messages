package cmd

import (
	"crypto/rsa"
	"flag"
	"fmt"
	"os"

	"github.com/andrew-candela/messages/messages"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type config struct {
	Name  string
	Users []deliver_config
}

type deliver_config struct {
	Host string
	Key  string
}

func getKey(keyFile string) *rsa.PublicKey {
	k, _ := messages.ReadExistingKey(keyFile)
	return &k.PublicKey
}

func makeTargets(k *rsa.PublicKey) *[]messages.RecipientDetails {
	return &[]messages.RecipientDetails{
		{
			DestinationHostPort: "10.0.0.176:1053",
			PublicKey:           *k,
		},
		{
			DestinationHostPort: "aim.andrewcandela.com:1053",
			PublicKey:           *k,
		},
	}
}

func produce(keyFile string, group_name string, recip_config []config) {
	var host string
	var user string
	flag.StringVar(&host, "host", "10.0.0.186:1053", "Provide the HOST:PORT to write to")
	flag.StringVar(&user, "user", os.Getenv("USER"), "What do you call yourself?")
	flag.Parse()
	k := getKey(keyFile)
	fmt.Println(recip_config)
	targets := makeTargets(k)
	messages.ProduceMessages(*targets, user)
}

func init() {
	rootCmd.AddCommand(produceCommand)
}

var produceCommand = &cobra.Command{
	Use: "produce",
	Run: func(cmd *cobra.Command, args []string) {
		var recipConfig []config
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