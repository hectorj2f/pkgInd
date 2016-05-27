package client

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Show pkgIndctl version",
		Long:  "Show pkgIndctl version",
		Run:   versionRun,
	}

	projectVersion string
	projectBuild   string
)

func versionRun(cmd *cobra.Command, args []string) {
	fmt.Printf("pkgIndctl version %s (build %s)\n", projectVersion, projectBuild)
}
