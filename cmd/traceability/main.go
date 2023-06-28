package main

import (
	"fmt"
	"os"

	_ "github.com/Axway/agent-sdk/pkg/traceability"
	agentCmd "github.com/Axway/agents-webmethods/pkg/cmd/traceability"
)

func main() {
	if err := agentCmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
