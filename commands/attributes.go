package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func CreateAttributesCommand(globalFlags *GlobalFlags) (*cobra.Command) {
	var identity_type string
	var identity_id int64

	var ListEntityCmd = &cobra.Command{
		Use:   "list",
		Short: "List entity details",
		Long:  `List an entity.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			//Get unity api client
			client, err := createUnityClient(globalFlags)
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

			//group_path := "/clarin"
			entity, err := client.Entity(identity_id, identity_type_ptr)
			if err != nil {
				fmt.Printf("Failed to resolve entity type=\"%v\" and value=\"%v\". Error: %v\n", identity_type, identity_id, err)
				os.Exit(1)
			}

			//Process response
			entity.Print()
			return nil
		},
	}

	ListEntityCmd.Flags().StringVarP(&identity_type, "type", "t", "", "Identity type")
	ListEntityCmd.Flags().Int64VarP(&identity_id, "id", "i", 1, "Identity id")

	var EntityCmd = &cobra.Command{
		Use:   "attributes",
		Short: "Attribute management",
		Long:  `List, create, delete, query and manage entity attributes.`,
	}

	EntityCmd.AddCommand(ListEntityCmd)

	return EntityCmd
}
