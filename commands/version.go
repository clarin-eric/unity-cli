package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

func CreateVersionCommand(globalFlags *GlobalFlags) (*cobra.Command) {
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of unity-cli",
		Long:  `All software has versions. This is unity-cli's.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			printCheckerVersion()
			return nil
		},
	}
	return versionCmd
}

func printCheckerVersion() {
	fmt.Printf("unity-cli v%s by CLARIN ERIC\n", "1.0.0-beta")
}
