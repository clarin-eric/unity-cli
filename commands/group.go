package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"unity-client/api"
	"os"
)

func CreateGroupCommand(globalFlags *GlobalFlags) (*cobra.Command) {
	var group_path string
	var recursive bool

	//
	// List
	//
	var GroupListCmd = &cobra.Command{
		Use:   "list",
		Short: "List group details",
		Long:  `List subgroups and members of this group.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			//Get unity api client
			client, err := api.GetUnityApiV1(globalFlags.Api_base, globalFlags.Verbose, globalFlags.Insecure, globalFlags.Username, globalFlags.Password)
			if err != nil {
				fmt.Printf("Failed to initialize unity client. Error: %v\n", err)
				os.Exit(1)
			}

			//Resolve
			group, err := client.GetGroup(&group_path)
			if err != nil {
				fmt.Printf("Failed to get group with path=\"%v\". Error: %v\n", group_path, err)
				os.Exit(1)
			}

			//Process response
			group.Print()

			return nil
		},
	}
	GroupListCmd.Flags().StringVarP(&group_path, "path", "P", "/", "Group path")

	//
	// Create
	//
	var GroupCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new group",
		Long:  `Create a new group.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			//Get unity api client
			client, err := api.GetUnityApiV1(globalFlags.Api_base, globalFlags.Verbose, globalFlags.Insecure, globalFlags.Username, globalFlags.Password)
			if err != nil {
				fmt.Printf("Failed to initialize unity client. Error: %v\n", err)
				os.Exit(1)
			}

			//Resolve
			err = client.CreateGroup(&group_path)
			if err != nil {
				fmt.Printf("Failed to create group with path=\"%v\". Error: %v\n", group_path, err)
				os.Exit(1)
			}

			fmt.Printf("Succesfully created group with path %s\n", group_path)

			return nil
		},
	}
	GroupCreateCmd.Flags().StringVarP(&group_path, "path", "P", "/", "Group path")

	//
	// Delete
	//
	var GroupDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a group",
		Long:  `Delete a group.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			//Get unity api client
			client, err := api.GetUnityApiV1(globalFlags.Api_base, globalFlags.Verbose, globalFlags.Insecure, globalFlags.Username, globalFlags.Password)
			if err != nil {
				fmt.Printf("Failed to initialize unity client. Error: %v\n", err)
				os.Exit(1)
			}

			//Resolve
			err = client.DeleteGroup(&group_path, &recursive)
			if err != nil {
				fmt.Printf("Failed to create group with path=\"%v\". Error: %v\n", group_path, err)
				os.Exit(1)
			}

			fmt.Printf("Succesfully removed group with path %s\n", group_path)

			return nil
		},
	}
	GroupDeleteCmd.Flags().StringVarP(&group_path, "path", "P", "/", "Group path")
	GroupDeleteCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "Enforce recursive removal")


	//
	// Container
	//
	var GroupCmd = &cobra.Command{
		Use:   "groups",
		Short: "Group management",
		Long:  `List, create, delete and query groups.`,
	}
	GroupCmd.AddCommand(GroupListCmd)
	GroupCmd.AddCommand(GroupCreateCmd)
	GroupCmd.AddCommand(GroupDeleteCmd)
	return GroupCmd
}