package main

import (
	"fmt"
	"os"

	"github.com/hectorj2f/pkgInd/cmd"
)

func main() {
	if err := cmd.PkgIndCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
