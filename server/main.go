package main

import (
	"fmt"
	"github.com/warmans/rsk-search/cmd"
	"os"
)

func main() {
	if err := cmd.RootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
