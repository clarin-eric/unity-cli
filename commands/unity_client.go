package commands

import (
	"github.com/spf13/cobra"
	"fmt"
)

type GlobalFlags struct {
	Verbose bool
	Api_base string
	Insecure bool
	Username string
	Password string
}

var UnityCliCmd = &cobra.Command{
	Use:   "unity-client",
	Short: "CLI client for the unity idm REST interface",
	Long: `Command line interface (CLI) client for the unity idm REST interface.`,
}

func Execute() {
	var globalFlags GlobalFlags

	//Add all global flags
	UnityCliCmd.PersistentFlags().BoolVarP(&globalFlags.Verbose, "verbose", "v", false, "Run in verbose mode")
	UnityCliCmd.PersistentFlags().BoolVarP(&globalFlags.Insecure, "insecure", "k", false, "Don't verify certificates")
	UnityCliCmd.PersistentFlags().StringVarP(&globalFlags.Api_base, "base", "b", "https://localhost:2443", "Base url for the unity-idm rest endpoint")
	UnityCliCmd.PersistentFlags().StringVarP(&globalFlags.Username, "user", "u", "admin", "Specify username")
	UnityCliCmd.PersistentFlags().StringVarP(&globalFlags.Password, "pass", "p", "Admin12345", "Specify password")

	//Add subcommands
	UnityCliCmd.AddCommand(CreateVersionCommand(&globalFlags))
	UnityCliCmd.AddCommand(CreateResolveCommand(&globalFlags))
	UnityCliCmd.AddCommand(CreateEntityCommand(&globalFlags))
	//Process all arguments
	UnityCliCmd.Execute()
}

func init() {
	fmt.Sprintf("Init")
}