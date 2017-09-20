package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"unity-client/api"
	"os"
)

func CreateResolveCommand(globalFlags *GlobalFlags) (*cobra.Command) {
	var identity_type string
	var identity_value string

	var ResolveCmd = &cobra.Command{
		Use:   "resolve",
		Short: "Resolve an identity",
		Long:  `Resolve an identity based on type and value.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			//Get unity api client
			client, err := api.GetUnityApiV1(globalFlags.Api_base, globalFlags.Verbose, globalFlags.Insecure, globalFlags.Username, globalFlags.Password)
			if err != nil {
				fmt.Printf("Failed to initialize unity client. Error: %v\n", err)
				os.Exit(1)
			}

			//Resolve
			entity, err := client.Resolve(identity_type, identity_value)
			if err != nil {
				fmt.Printf("Failed to resolve entity type=\"%v\" and value=\"%v\". Error: %v\n", identity_type, identity_value, err)
				os.Exit(1)
			}

			//Process response
			fmt.Printf("Response: %v\n", entity)
			return nil
		},
	}

	ResolveCmd.Flags().StringVarP(&identity_type, "type", "t", "email", "Identity type")
	ResolveCmd.Flags().StringVarP(&identity_value, "value", "V", "john doe", "Identity value")

	return ResolveCmd
}