package main

import (
	"fmt"

	"github.com/hectorj2f/pkgind/log"
	"github.com/spf13/cobra"
)

type queryCmdFlags struct {
	packageName string
	verbose     bool
	listenIP    string
	listenPort  int
}

func (f queryCmdFlags) Validate() {
	if f.packageName == "" {
		log.Logger().Fatal("package is required")
	}
}

var (
	queryCmd = &cobra.Command{
		Use:   "query",
		Short: "Query pkgIndctl performs query commands to the server",
		Long:  "Query pkgIndctl performs query commands to the server",
		Run:   queryRun,
	}

	queryFlags = queryCmdFlags{}
)

func init() {
	queryCmd.Flags().StringVar(&queryFlags.packageName, "package", "", "package to ask")
	queryCmd.Flags().StringVar(&queryFlags.listenIP, "ip", listenerDefaultIP, "IP to listen")
	queryCmd.Flags().IntVar(&queryFlags.listenPort, "port", listenerDefaultPort, "port to listen")
	queryCmd.Flags().BoolVar(&queryFlags.verbose, "verbose", false, "verbose output")
}

func queryRun(cmd *cobra.Command, args []string) {
	queryFlags.Validate()
	client := &Client{
		listenIP:    queryFlags.listenIP,
		listenPort:  queryFlags.listenPort,
		enableDebug: queryFlags.verbose,
	}
	resp, err := client.executeQuery(queryFlags.packageName)
	if err != nil {
		panic(fmt.Sprintf("unable to query for the package '%s': %v", queryFlags.packageName, err))
	}
	fmt.Println(resp)
}
