package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"clarin/unity-cli/api"
	"os"
	"time"
	"clarin/unity-cli/report"

	"strings"
)

type int64arr []int64
func (a int64arr) Len() int { return len(a) }
func (a int64arr) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a int64arr) Less(i, j int) bool { return a[i] < a[j] }

func CreateReportCommand(globalFlags *GlobalFlags) (*cobra.Command) {
	var (
		outputs string
		kind string
	)

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

			fmt.Printf("Generated missing attributes in %dns\n", d.Nanoseconds())

			attribute_set := []string{"clarin-full-name", "member", "clarin-motivation"}
			report.Compute(entities, attribute_set)

			//for _, o := range SplitOutputs(outputs) {
			ExportJson("report_missing_attributes.json", report)
			//}

			return nil
		},
	}
	ReportMissingAttributesCmd.Flags().StringVarP(&outputs, "outputs", "o", "json", "List of output formats (comma separated). Supported values: json")

	//
	// Root command
	//
	var ReportStatsCmd = &cobra.Command{
		Use:   "entities",
		Short: "Report entities",
		Long:  `Report entity informaton.`,
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

			fmt.Printf("Generated entity report in %dns\n", d.Nanoseconds())

			report.Compute(entities)

			for _, o := range SplitOutputs(outputs) {
				if o == "tsv" || o == "csv" {
					separator := GetSeparator(o)
					if kind == "both" {
						ExportSeparated(fmt.Sprintf("%s_report.%s", "anonymous", o), report.GetReportAsArray("anonymous"), separator)
						ExportSeparated(fmt.Sprintf("%s_report.%s", "personal", o), report.GetReportAsArray("personal"), separator)
					} else {
						ExportSeparated(fmt.Sprintf("%s_report.%s", kind, o), report.GetReportAsArray(kind), separator)
					}
				} else if o == "json" {
					if kind == "both" {
						ExportJson(fmt.Sprintf("%s_report.json", "anonymous"), report.GetReport("anonymous"))
						ExportJson(fmt.Sprintf("%s_report.json", "personal"), report.GetReport("personal"))
					} else {
						ExportJson(fmt.Sprintf("%s_report.json", kind), report.GetReport(kind))
					}
				} else {
					fmt.Printf("Unsupported output format: %s\n", o)
				}
			}
			return nil
		},
	}
	ReportStatsCmd.Flags().StringVarP(&outputs, "outputs", "o", "json", "List of output formats (comma separated). Supported values: json,csv,tsv")
	ReportStatsCmd.Flags().StringVarP(&kind, "type", "t", "both", "Type of output. Supported values: anonymous,personal,both")

	//
	// Root command
	//
	var ReportCmd = &cobra.Command{
		Use:   "report",
		Short: "Report commands",
		Long:  `Report informaton.`,
	}
	ReportCmd.AddCommand(ReportMissingAttributesCmd, ReportStatsCmd)

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

func GetSeparator(output string) (rune) {
	separator := ','
	if output == "tsv" {
		separator = '\t'
	}
	return separator
}

func SplitOutputs(outputs string) ([]string) {
	return strings.Split(outputs, ",")
}