package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"clarin/unity-cli/api"
	"os"
	"time"
	"clarin/unity-cli/report"

)

type int64arr []int64
func (a int64arr) Len() int { return len(a) }
func (a int64arr) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a int64arr) Less(i, j int) bool { return a[i] < a[j] }

func CreateReportCommand(globalFlags *GlobalFlags) (*cobra.Command) {
	var (
		output string
		kind string
		fix_attributes bool
	)

	var TestReport = &cobra.Command{
		Use:   "test",
		Short: "Resolve an identity",
		Long:  `Resolve an identity based on type and value. This yields the same output as "entities ls".`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := createUnityClient(globalFlags)
			if err != nil {
				fmt.Printf("Failed to initialize unity client. Error: %v\n", err)
				return nil
			}

			attributes, err := client.GetAttributes(261, "/")
			if err != nil {
				fmt.Printf("%s\n", err)
				return nil
			}

			for _, a := range attributes {
				fmt.Printf("%50s %20s %20s\n", a.Name, a.GroupPath,  a.Visibility)
				for _, v := range a.Values {
					fmt.Printf("  %s\n", v)
				}
			}

			return nil
		},
	}

	var TestUpdateReport = &cobra.Command{
		Use:   "test-update",
		Short: "Resolve an identity",
		Long:  `Resolve an identity based on type and value. This yields the same output as "entities ls".`,
		Run: func(cmd *cobra.Command, args []string) {
			client, err := createUnityClient(globalFlags)
			if err != nil {
				fmt.Printf("Failed to initialize unity client. Error: %v\n", err)
			}

			fmt.Printf("Fetching list of entities\n")
			t1 := time.Now()
			entities, err := GetAllEntitiesForGroupWithClient(globalFlags, "/clarin", client)
			if err != nil {
				fmt.Printf("Failed to fecth entities. Error: %v\n", err)
				os.Exit(1)
			}
			d := time.Since(t1)

			fmt.Printf("Fetched list of %d entitites in %.3fms\n", len(entities), float64(d.Nanoseconds())/1000000.0)



			for _, entity := range entities {
				entity.RepairAttributeSet(client)
			}
		},
	}

	var ReportMissingAttributesCmd = &cobra.Command{
		Use:   "missing_attributes",
		Short: "Report missing attributes",
		Long:  `Report missing attributes.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			report := report.ReportEmptyAttributes{}

			t1 := time.Now()
			entities, err := GetAllEntitiesForGroup(globalFlags, "/clarin")
			if err != nil {
				fmt.Printf("Failed to fecth entities. Error: %v\n", err)
				os.Exit(1)
			}
			d := time.Since(t1)

			fmt.Printf("Generated in %dns\n", d.Nanoseconds())

			attribute_set := []string{"clarin-full-name", "member", "clarin-motivation"}
			report.Compute(entities, attribute_set)

			fmt.Printf("Json:\n%s\n", string(report.ToJson()))

			return nil
		},
	}
	ReportMissingAttributesCmd.Flags().BoolVar(&fix_attributes, "fix",  false,"Fix missing attribute values by setting default values")

	//
	// Root command
	//
	var ReportStatsCmd = &cobra.Command{
		Use:   "report",
		Short: "Report commands",
		Long:  `Report informaton.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			report := report.Report{}

			client, err := createUnityClient(globalFlags)
			if err != nil {
				fmt.Printf("Failed to initialize unity client. Error: %v\n", err)
			}

			t1 := time.Now()
			entities, err := GetAllEntitiesForGroupWithClient(globalFlags, "/", client)
			if err != nil {
				fmt.Printf("Failed to fecth entities. Error: %v\n", err)
				os.Exit(1)
			}

			for _, entity := range entities {
				client.GetEntityGroups(entity.Id)
			}
			d := time.Since(t1)

			fmt.Printf("Generated in %dns\n", d.Nanoseconds())

			report.Compute(entities)
			report.Write(kind, output)

			return nil
		},
	}

	//
	// Root command
	//
	var ReportCmd = &cobra.Command{
		Use:   "report",
		Short: "Report commands",
		Long:  `Report informaton.`,
	}

	ReportCmd.AddCommand(ReportMissingAttributesCmd, ReportStatsCmd, TestReport, TestUpdateReport)
	ReportCmd.PersistentFlags().StringVarP(&output, "output", "o", "pretty", "Output format. Supported values: pretty,json,csv,tsv,google")
	ReportCmd.PersistentFlags().StringVarP(&kind, "type", "t", "both", "Type of output. Supported values: anonymous,personal,both")

	return ReportCmd
}

func GetAllEntitiesForGroup(globalFlags *GlobalFlags, group_path string) ([]api.Entity, error) {
	return GetAllEntitiesForGroupWithClient(globalFlags, group_path, nil)
}

func GetAllEntitiesForGroupWithClient(globalFlags *GlobalFlags, group_path string, client *api.UnityApiV1) ([]api.Entity, error) {
	var err error
	empty := []api.Entity{}

	//Get unity api client
	if client == nil {
		client, err = createUnityClient(globalFlags)
		if err != nil {
			return empty, fmt.Errorf("Failed to initialize unity client. Error: %v\n", err)
		}
	}
	//Resolve group
	group, err := client.GetGroup(&group_path)
	if err != nil {
		return empty, fmt.Errorf("Failed to get group with path=\"%v\". Error: %v\n", group_path, err)
	}

	var entities []api.Entity
	for _, id := range group.Members {
		//Process response
		entity, err := client.Entity(id, nil)
		if err != nil {
			fmt.Printf("Failed to resolve entity with id=\"%v\". Error: %v\n", id, err)
		} else {
			entities = append(entities, *entity)
		}
	}

	return entities, nil
}