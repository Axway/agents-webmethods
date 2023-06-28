package errors

import "github.com/Axway/agent-sdk/pkg/util/errors"

// Error definitions
var (
	ErrGatewayConfig = errors.Newf(3500, "error apigateway.getHeaders is set to true and apigateway.%s not set in config")

	ErrV7HealthcheckHost = errors.Newf(3501, "%s Failed. Error communicating with API Gateway: %s. Check API Gateway host configuration values")
	ErrV7HealthcheckAuth = errors.Newf(3502, "%s Failed. Error sending request to API Gateway. HTTP response code %v. Check API Gateway authentication configuration values")

	// Event Processing
	ErrEventNoMsg         = errors.Newf(3510, "the log event had no message field: %s")
	ErrEventMsgStructure  = errors.Newf(3511, "could not parse the log event: %s")
	ErrTrxnDataGet        = errors.New(3512, "could not retrieve the transaction data from API Gateway")
	ErrTrxnDataProcess    = errors.New(3513, "could not process the transaction data")
	ErrTrxnHeaders        = errors.Newf(3514, "could not process the transaction headers: %s")
	ErrProtocolStructure  = errors.Newf(3515, "could not parse the %s transaction details: %s")
	ErrCreateCondorEvent  = errors.New(3516, "error creating the Amplify Visibility event")
	ErrNotHTTPService     = errors.New(3517, "the event processor can only handle http services")
	ErrWrongProcessor     = errors.Newf(3518, "the message looks to be from the %s log but open traffic process is %s, check OPENTRAFFIC_LOG_INPUT")
	ErrInvalidInputConfig = errors.Newf(3519, "input configuration is not valid: %s")
	ErrConfigFile         = errors.New(3520, "could not find the 'traceability_agent' section in the configuration file")

	// API Gateway Communication
	ErrAPIGWRequest  = errors.New(3530, "error encountered sending a request to API Gateway")
	ErrAPIGWResponse = errors.Newf(3531, "unexpected HTTP response code %v in response from API Gateway")
)
