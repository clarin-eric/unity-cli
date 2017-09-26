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
	var group_path string

	var ResolveIdentityCmd = &cobra.Command{
		Use:   "identity",
		Short: "Resolve an identity",
		Long:  `Resolve an identity based on type and value. This yields the same output as "entities ls".`,
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
			entity.Print()

			return nil
		},
	}
	ResolveIdentityCmd.Flags().StringVarP(&identity_type, "type", "t", "email", "Identity type")
	ResolveIdentityCmd.Flags().StringVarP(&identity_value, "value", "V", "john doe", "Identity value")

	var ResolveGroupCmd = &cobra.Command{
		Use:   "group",
		Short: "Resolve entities in a group",
		Long:  `Resolve all entities in a group.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			//Get unity api client
			client, err := api.GetUnityApiV1(globalFlags.Api_base, globalFlags.Verbose, globalFlags.Insecure, globalFlags.Username, globalFlags.Password)
			if err != nil {
				fmt.Printf("Failed to initialize unity client. Error: %v\n", err)
				os.Exit(1)
			}

			//Resolve group
			group, err := client.GetGroup(&group_path)
			if err != nil {
				fmt.Printf("Failed to get group with path=\"%v\". Error: %v\n", group_path, err)
				os.Exit(1)
			}

			for _, id := range group.Members {
				//Process response
				entity, err := client.Entity(id, nil)
				if err != nil {
					fmt.Printf("Failed to resolve entity with id=\"%v\". Error: %v\n", id, err)
				} else {
					entity.Print()
				}
			}

			return nil
		},
	}
	ResolveGroupCmd.Flags().StringVarP(&group_path, "path", "P", "/", "Group path")

	//
	// Root command
	//
	var ResolveCmd = &cobra.Command{
		Use:   "resolve",
		Short: "Resolution commands",
		Long:  `Resolve informaton.`,
	}
	ResolveCmd.AddCommand(ResolveIdentityCmd)
	ResolveCmd.AddCommand(ResolveGroupCmd)
	return ResolveCmd
}