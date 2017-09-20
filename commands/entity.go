package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"unity-client/api"
	"os"
)

func CreateEntityCommand(globalFlags *GlobalFlags) (*cobra.Command) {
	var identity_type string
	var identity_id int64

	var ResolveCmd = &cobra.Command{
		Use:   "entity",
		Short: "Retrieve an entity",
		Long:  `Retrieve an entity.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			//Get unity api client
			client, err := api.GetUnityApiV1(globalFlags.Api_base, globalFlags.Verbose, globalFlags.Insecure, globalFlags.Username, globalFlags.Password)
			if err != nil {
				fmt.Printf("Failed to initialize unity client. Error: %v\n", err)
				os.Exit(1)
			}

			//Resolve
			var identity_type_ptr *string
			if identity_type == "" {
				identity_type_ptr = nil
			} else {
				identity_type_ptr = &identity_type
			}

			entity, err := client.Entity(identity_id, identity_type_ptr)
			if err != nil {
				fmt.Printf("Failed to resolve entity type=\"%v\" and value=\"%v\". Error: %v\n", identity_type, identity_id, err)
				os.Exit(1)
			}

			//Process response
			fmt.Printf("Response: %v\n", entity)
			return nil
		},
	}

	ResolveCmd.Flags().StringVarP(&identity_type, "type", "t", "", "Identity type")
	ResolveCmd.Flags().Int64VarP(&identity_id, "id", "i", 1, "Identity id")

	return ResolveCmd
}