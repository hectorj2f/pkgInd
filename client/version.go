package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Show pkgindctl version",
		Long:  "Show pkgindctl version",
		Run:   versionRun,
	}

	projectVersion string
	projectBuild   string
)

func versionRun(cmd *cobra.Command, args []string) {
	fmt.Printf("pkgindctl version %s (build %s)\n", projectVersion, projectBuild)
}
