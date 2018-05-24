package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"time"
	"clarin/unity-client/api"
	"encoding/json"
	"io/ioutil"
	"clarin/unity-client/report"
)

func CreateRepairCommand(globalFlags *GlobalFlags) (*cobra.Command) {
	var (
		output string
	)

	var repairEntities = &cobra.Command{
		Use:   "entities",
		Short: "Repair all entities",
		Long:  `Enumerate all entities and repair the attribute set for any missing values.`,
		Run: func(cmd *cobra.Command, args []string) {
			/*
			 * Retrieve list of entities
			 */
			fmt.Printf("Fetching list of entities\n")
			t1 := time.Now()
			client, entities, err := CreateClientAndGetEntities(globalFlags)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				os.Exit(1)
			}
			d := time.Since(t1)
			fmt.Printf("Fetched list of %d entitites in %.3fms\n", len(entities), float64(d.Nanoseconds())/1000000.0)

			/*
			 * Repair entities if needed
			 */
			fmt.Printf("Analyzing and repairing entities\n")
			repaired := 0
			t1 = time.Now()
			repaired_entities := []*api.AttributeRepairResult{}
			for _, entity := range entities {
				r := entity.RepairAttributeSet(client)
				if r.Updated {
					repaired++
				}
				repaired_entities =
					append(repaired_entities, r)
			}
			d = time.Since(t1)
			fmt.Printf("Analyzed %d and repaired %d entitites in %.3fms\n", len(entities), repaired, float64(d.Nanoseconds())/1000000.0)

			/*
			 * Export results
			 */

			filename := "repair.json"
			fmt.Printf("Writing result to: %s\n", filename)

			json_bytes, err := json.Marshal(repaired_entities)
			if err != nil {
				fmt.Printf("Failed to marshall to json: %s1\n", err)
				return
			}

			r := report.Report{}
			file, err := r.GetFile(filename)
			if err != nil {
				fmt.Printf("Failed to create file. Error: %s\n", err)
				return
			}
			err = ioutil.WriteFile(file.Name(), json_bytes, 0644)
			if err != nil {
				fmt.Printf("Failed to write json data to file. Error: %s\n", err)
				return
			}
		},
	}

	//
	// Root command
	//
	var ReportCmd = &cobra.Command{
		Use:   "repair",
		Short: "Repair commands",
		Long:  `Various commands to make repairs to the unity database.`,
	}

	ReportCmd.AddCommand(repairEntities)
	ReportCmd.PersistentFlags().StringVarP(&output, "output", "o", "pretty", "Output format. Supported values: pretty,json,csv,tsv,google")

	return ReportCmd
}

func CreateClientAndGetEntities(globalFlags *GlobalFlags) (*api.UnityApiV1, []api.Entity, error) {
	client, err := createUnityClient(globalFlags)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to initialize unity client. Error: %v\n", err)
	}

	entities, err := GetAllEntitiesForGroupWithClient(globalFlags, "/clarin", client)
	if err != nil {
		return client, nil, fmt.Errorf("Failed to fecth entities. Error: %v\n", err)
	}

	return client, entities, nil
}