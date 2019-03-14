package main

import (
	"os"

	"github.com/S1mplee/eventstorebeat/cmd"

	_ "github.com/S1mplee/eventstorebeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
