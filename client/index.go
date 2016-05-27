package client

import (
	"fmt"
	"strings"

	"github.com/hectorj2f/pkgInd/log"
	"github.com/spf13/cobra"
)

type indexCmdFlags struct {
	packageName         string
	packageDependencies string
	listenPort          int
	listenIP            string
	verbose             bool
}

func (f indexCmdFlags) Validate() {
	if f.packageName == "" {
		log.Logger().Fatal("package is required")
	}
}

var (
	indexCmd = &cobra.Command{
		Use:   "index",
		Short: "index pkgIndctl performs index commands to the server",
		Long:  "index pkgIndctl performs index commands to the server",
		Run:   indexRun,
	}

	indexFlags = indexCmdFlags{}

	dependencies []string
)

func init() {
	indexCmd.Flags().StringVar(&indexFlags.packageName, "package", "", "package name")
	indexCmd.Flags().StringVar(&indexFlags.packageDependencies, "dependencies", "", "package depencies")
	indexCmd.Flags().StringVar(&indexFlags.listenIP, "ip", listenerDefaultIP, "IP to listen")
	indexCmd.Flags().IntVar(&indexFlags.listenPort, "port", listenerDefaultPort, "port to listen")
	indexCmd.Flags().BoolVar(&indexFlags.verbose, "verbose", false, "verbose output")
}

func indexRun(cmd *cobra.Command, args []string) {
	dependencies = strings.Split(indexFlags.packageDependencies, ",")

	indexFlags.Validate()
	client := &Client{
		listenIP:   indexFlags.listenIP,
		listenPort: indexFlags.listenPort,
	}

	resp, err := client.executeIndex(indexFlags.packageName, dependencies)
	if err != nil {
		panic(fmt.Sprintf("unable to index the package '%s': %v", indexFlags.packageName, err))
	}
	fmt.Println(resp)
}
