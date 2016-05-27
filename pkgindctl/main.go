package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	listenerDefaultIP   = "127.0.0.1"
	listenerDefaultPort = 8080
)

var (
	pkgIndCtlCmd = &cobra.Command{
		Use:   "pkgIndctl",
		Short: "pkgIndctl is a package indexer...",
		Long:  "pkgIndctl is a package indexer...",
		Run:   pkgIndCtlRun,
	}
)

func pkgIndCtlRun(cmd *cobra.Command, args []string) {
	cmd.Help()
}

func main() {
	pkgIndCtlCmd.AddCommand(indexCmd)
	pkgIndCtlCmd.AddCommand(queryCmd)
	pkgIndCtlCmd.AddCommand(removeCmd)
	pkgIndCtlCmd.AddCommand(versionCmd)

	if err := pkgIndCtlCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
