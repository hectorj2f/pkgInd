package cmd

import (
	"github.com/spf13/cobra"
)

var (
	PkgIndCmd = &cobra.Command{
		Use:   "pkgInd",
		Short: "pkgInd is a package indexer...",
		Long:  "pkgInd is a package indexer...",
		Run:   pkgIndRun,
	}
)

func init() {
	PkgIndCmd.AddCommand(versionCmd)
	PkgIndCmd.AddCommand(startCmd)
}

func pkgIndRun(cmd *cobra.Command, args []string) {
	cmd.Help()
}
