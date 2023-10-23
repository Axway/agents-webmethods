package errors

import "github.com/Axway/agent-sdk/pkg/util/errors"

// Error definitions
var (
	ErrConfigFile = errors.New(3520, "could not find the 'traceability_agent' section in the configuration file")
)
