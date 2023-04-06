package cmd

import (
	"github.com/andrew-candela/messages/messages"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	print, write bool
)

var keyCommand = &cobra.Command{
	Use: "keygen",
	Run: func(cmd *cobra.Command, args []string) {

		viper.GetViper().ReadInConfig()
		keyFile := viper.GetString("private_key_file")
		if write {
			newKey := messages.GenerateRandomKey()
			messages.WriteKeyToDisk(newKey, keyFile)
		}
		key, _ := messages.ReadExistingKey(keyFile)
		if print {
			messages.DisplayPublicKey(key)
		}
	},
}

func init() {
	keyCommand.Flags().BoolVar(&print, "print", false, "--print will print the public key")
	keyCommand.Flags().BoolVar(&write, "write", false, "Generates a random RSA private key and writes to configured keyfile")
	rootCmd.AddCommand(keyCommand)
}
