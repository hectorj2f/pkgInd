package cmd

import (
	"flag"
	"fmt"

	"github.com/hectorj2f/pkgind/log"
	"github.com/hectorj2f/pkgind/server"
	"github.com/hectorj2f/pkgind/utils"
	"github.com/spf13/cobra"
)

const (
	listenerDefaultIP   = "0.0.0.0"
	listenerDefaultPort = 8080

	maxPortNumber = 65535
)

type startCmdFlags struct {
	listenIP   string
	listenPort int
	verbose    bool
}

func (f startCmdFlags) Validate() {
	if f.listenPort <= 0 || f.listenPort > 65535 {
		log.Logger().Fatalf("wrong port range. --port '%d' has to be greater than '0' and lower than '%d'", f.listenPort, maxPortNumber)
	}

	if f.listenIP == "" {
		log.Logger().Fatal("listenIP is not a correct IP address")
	} else if !utils.ValidIP4(f.listenIP) {
		log.Logger().Fatal("listenIP is not a valid IP address")
	}
}

var (
	startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start pkgInd server",
		Long:  "Start pkgInd server",
		Run:   startRun,
	}

	startFlags = startCmdFlags{}

	argvNrThreads = flag.Int("n", 2, "number of threads")
)

func init() {
	startCmd.Flags().StringVar(&startFlags.listenIP, "ip", listenerDefaultIP, "IP to listen")
	startCmd.Flags().IntVar(&startFlags.listenPort, "port", listenerDefaultPort, "port to listen")
	startCmd.Flags().BoolVar(&startFlags.verbose, "verbose", false, "verbose output")
}

func startRun(cmd *cobra.Command, args []string) {
	startFlags.Validate()
	pkgIndServer := server.NewServer()

	if err := pkgIndServer.Start(fmt.Sprintf("%s:%d", startFlags.listenIP, startFlags.listenPort), "tcp"); err != nil {
		log.Logger().Fatalf("unable to start server: %v", err)
	}
}
