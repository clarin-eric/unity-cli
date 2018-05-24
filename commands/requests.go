package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"clarin/unity-cli/api"
	"os"
	"encoding/csv"
	"log"
	"time"
	"strings"
)

func CreateRequestsCommand(globalFlags *GlobalFlags) (*cobra.Command) {
	var identity_type string
	var identity_id int64

	var ListEntityCmd = &cobra.Command{
		Use:   "list",
		Short: "List request details",
		Long:  `List request details.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			//Get unity api client
			client, err := createUnityClient(globalFlags)
			if err != nil {
				fmt.Printf("Failed to initialize unity client. Error: %v\n", err)
				os.Exit(1)
			}

			requests, err := client.GetRegistrationRequests()
			if err != nil {
				fmt.Printf("Failed to get account requests\". Error: %v\n", err)
				os.Exit(1)
			}

			output_format := "csv"

			//Process response
			switch output_format {
			case "csv":
				outputToCsv(requests)
			case "tsv":
				outputToTsv(requests)
			case "stdout":
				outputToTsv(requests)
			default:
				fmt.Printf("Unkown output format: %s\n", output_format)
			}

			return nil
		},
	}

	ListEntityCmd.Flags().StringVarP(&identity_type, "type", "t", "", "Identity type")
	ListEntityCmd.Flags().Int64VarP(&identity_id, "id", "i", 1, "Identity id")


	var RequestsCmd = &cobra.Command{
		Use:   "requests",
		Short: "Request management",
		Long:  `List, create, delete and query entities.`,
	}

	RequestsCmd.AddCommand(ListEntityCmd)

	return RequestsCmd
}

func outputToStdout(requests []api.RegistrationRequest) {
	fmt.Printf("Found %d requests\n", len(requests))
	for _, request := range requests {
		fmt.Printf("Request:\n")
		fmt.Printf("  Id:        %s\n", request.RequestId)
		fmt.Printf("  Timestamp: %d\n", request.Timestamp)
		fmt.Printf("  Status:    %s\n", request.Status)
		fmt.Printf("  Identities:\n")
		for _, identity := range request.Content.Identities {
			entity_id := ""
			if identity.EntityId > 0 {
				entity_id = fmt.Sprintf("%d", identity.EntityId)
			}
			fmt.Printf("    Entity id: %s\n", entity_id)
			fmt.Printf("    Type:      %s\n", identity.TypeId)
			fmt.Printf("    Value:     %s\n", identity.Value)
		}
		fmt.Printf("  Attributes:\n")
		for _, attribute := range request.Content.Attributes {
			fmt.Printf("    Name: %s\n", attribute.Name)
			fmt.Printf("    Values:\n")
			for _, value := range attribute.Values {
				fmt.Printf("      - %s\n", value)
			}
		}
	}
}

func addToSlice(data []string, value string) ([]string) {
	exists := false
	for _, v := range data {
		if v == value {
			exists = true
		}
	}

	if !exists {
		data = append(data, value)
	}

	return data
}

func sliceToString(values []string) (string) {
	result := ""
	if len(values) > 0 {
		result = strings.Replace(values[0], "\n", " ", -1)
		for i := 1; i < len(values); i++ {
			value := strings.Replace(values[i], "\n", " ", -1)
			result = fmt.Sprintf("%s %s", result, value)
		}
	}
	return result
}

func getAttributeByName(name string, request api.RegistrationRequest) (string) {
	for _, attribute := range request.Content.Attributes {
		if attribute != nil && attribute.Name == name {
			return sliceToString(attribute.Values)
		}
	}
	return ""
}

func outputToCsv(requests []api.RegistrationRequest) {
	outputToCsvWithCustomSeparator(requests, ',')
}

func outputToTsv(requests []api.RegistrationRequest) {
	outputToCsvWithCustomSeparator(requests, '\t')
}

func outputToCsvWithCustomSeparator(requests []api.RegistrationRequest, separator rune) {
	var data [][]string

	var atrribute_names []string
	for _, request := range requests {
		for _, attribute := range request.Content.Attributes {
			if attribute != nil {
				atrribute_names = addToSlice(atrribute_names, attribute.Name)
			}
		}
	}

	//Add header
	header := []string{"Request Id", "Timestamp", "Status", "Email Identity"}
	for _, a := range atrribute_names {
		header = append(header, a)
	}
	data = append(data, header)

	//Add data rows
	for _, request := range requests {
		email_identity := ""
		for _, identity := range request.Content.Identities {
			if identity.TypeId == "email" {
				email_identity= request.Content.Identities[0].Value
			}
		}

		row := []string{
			request.RequestId,
			time.Unix(request.Timestamp/1000, 0).String(),
			request.Status,
			email_identity,
		}

		for _, a := range atrribute_names {
			row = append(row, getAttributeByName(a, request))
		}

		data = append(data, row)
	}

	//Write csv to stdout
	writer := csv.NewWriter(os.Stdout)
	writer.Comma = separator
	defer writer.Flush()

	for _, value := range data {
		err := writer.Write(value)
		checkError("Cannot write to file", err)
	}
}

func checkError(message string, err error) {
	if err != nil {
		log.Printf(message, err)
	}
}