package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Show pkgInd version",
		Long:  "Show pkgInd version",
		Run:   versionRun,
	}

	projectVersion string
	projectBuild   string
)

func versionRun(cmd *cobra.Command, args []string) {
	fmt.Printf("pkgInd version %s (build %s)\n", projectVersion, projectBuild)
}
