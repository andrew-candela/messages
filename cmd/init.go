package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/andrew-candela/messages/messages"
	"github.com/spf13/cobra"
)

func createUDPMConfigFile(example_config_path string, out_path string) {
	config_file_contents, err := os.ReadFile(example_config_path)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(out_path, config_file_contents, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func createUDPMConfig() (string, string) {
	home, err := os.UserHomeDir()
	udpm_path := filepath.Join(home, ".udpm")
	udpm_config := filepath.Join(udpm_path, "config")
	if err != nil {
		panic(err)
	}
	_ = os.MkdirAll(udpm_path, os.ModePerm)
	createUDPMConfigFile("sample_config.toml", udpm_config)
	return udpm_path, udpm_config

}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create the config UDPM needs",
	Long: `Creates the ~/.udpm/ directory, with a few files:
	config: a toml file
	udpm_id_rsa: a randomly generated RSA private key file in PKCS #1, ASN.1 DER form.`,
	Run: func(cmd *cobra.Command, args []string) {
		udpm_path, _ := createUDPMConfig()
		fmt.Println("Created udpm config dir:", udpm_path)
		messages.WriteKeyToDisk(messages.GenerateRandomKey(), filepath.Join(udpm_path, "udpm_id_rsa"))
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
