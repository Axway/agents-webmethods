package main

import (
	"fmt"
	"os"

	"git.ecd.axway.org/apigov/agents-webmethods/pkg/cmd/discovery"
)

func main() {
	if err := discovery.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
