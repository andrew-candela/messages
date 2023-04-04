package cmd

import (
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

func makeTargets(conf config) *[]messages.RecipientDetails {
	var recips []messages.RecipientDetails
	for _, user_conf := range conf.Users {
		recips = append(recips, messages.RecipientDetails{
			DestinationHostPort: user_conf.Host,
			PublicKey:           messages.ParsePublicKey(user_conf.Key),
		})
	}
	return &recips
}

func produce(keyFile string, group_name string, recip_config config) {
	var host string
	var user string
	flag.StringVar(&host, "host", "10.0.0.186:1053", "Provide the HOST:PORT to write to")
	flag.StringVar(&user, "user", os.Getenv("USER"), "What do you call yourself?")
	flag.Parse()
	fmt.Println(recip_config)
	targets := makeTargets(recip_config)
	messages.ProduceMessages(*targets, user)
}

func init() {
	rootCmd.AddCommand(produceCommand)
}

var produceCommand = &cobra.Command{
	Use: "produce",
	Run: func(cmd *cobra.Command, args []string) {
		var recipConfig config
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
