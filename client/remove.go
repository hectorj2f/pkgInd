package main

import (
	"fmt"

	"github.com/hectorj2f/pkgind/log"
	"github.com/spf13/cobra"
)

type removeCmdFlags struct {
	packageName string
	listenPort  int
	listenIP    string
	verbose     bool
}

func (f removeCmdFlags) Validate() {
	if f.packageName == "" {
		log.Logger().Fatal("package is required")
	}
}

var (
	removeCmd = &cobra.Command{
		Use:   "remove",
		Short: "remove pkgIndctl performs remove commands to the server",
		Long:  "remove pkgIndctl performs remove commands to the server",
		Run:   removeRun,
	}

	removeFlags = removeCmdFlags{}
)

func init() {
	removeCmd.Flags().StringVar(&removeFlags.packageName, "package", "", "package to remove")
	removeCmd.Flags().StringVar(&removeFlags.listenIP, "ip", listenerDefaultIP, "IP to listen")
	removeCmd.Flags().IntVar(&removeFlags.listenPort, "port", listenerDefaultPort, "port to listen")
	removeCmd.Flags().BoolVar(&removeFlags.verbose, "verbose", false, "verbose output")
}

func removeRun(cmd *cobra.Command, args []string) {
	removeFlags.Validate()
	client := &Client{
		listenIP:   queryFlags.listenIP,
		listenPort: queryFlags.listenPort,
	}
	resp, err := client.executeRemove(removeFlags.packageName)
	if err != nil {
		panic(fmt.Sprintf("unable to remove for the package '%s': %v", removeFlags.packageName, err))
	}
	fmt.Println(resp)

}
