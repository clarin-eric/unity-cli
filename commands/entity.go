package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func CreateEntityCommand(globalFlags *GlobalFlags) (*cobra.Command) {
	var identity_type string
	var identity_id int64
	var attribute_name string
	var old_group string
	var new_group string

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

	var UpdateEntityAttributeCmd = &cobra.Command{
		Use:   "update",
		Short: "Update entity attribute",
		Long:  `Update entity attribute.`,
		Run: func(cmd *cobra.Command, args []string) {
			//Get unity api client
			client, err := createUnityClient(globalFlags)
			if err != nil {
				fmt.Printf("Failed to initialize unity client. Error: %v\n", err)
				os.Exit(1)
			}

			//Get entity
			entity, err := client.Entity(identity_id, nil)
			if err != nil {
				fmt.Printf("Failed to fetch entity type=\"%v\" and value=\"%v\". Error: %v\n", identity_type, identity_id, err)
				os.Exit(1)
			}

			//Update attribute
			err = entity.UpdateAttributeGroupPath(attribute_name, old_group, new_group, client)
			if err != nil {
				fmt.Printf("Failed to update entities atttribute. Error: %v\n", err)
				os.Exit(1)
			}
		},
	}
	UpdateEntityAttributeCmd.Flags().Int64VarP(&identity_id, "id", "i", 1, "Identity id")
	UpdateEntityAttributeCmd.Flags().StringVarP(&attribute_name, "name", "N", "", "Attribute name")
	UpdateEntityAttributeCmd.Flags().StringVarP(&old_group, "old_group", "o", "", "Old group")
	UpdateEntityAttributeCmd.Flags().StringVarP(&new_group, "new_group", "n", "", "New group")

	var EntityAttributeCmd = &cobra.Command{
		Use:   "attributes",
		Short: "Manage entity attributes",
		Long:  `Manage entity attributes.`,
	}
	EntityAttributeCmd.AddCommand(UpdateEntityAttributeCmd)




	var EntityCmd = &cobra.Command{
		Use:   "entities",
		Short: "Entity management",
		Long:  `List, create, delete and query entities.`,
	}

	EntityCmd.AddCommand(ListEntityCmd, EntityAttributeCmd)

	return EntityCmd
}