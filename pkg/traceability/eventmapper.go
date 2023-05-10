package traceability

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Axway/agent-sdk/pkg/agent"
	"github.com/Axway/agent-sdk/pkg/transaction"
	transutil "github.com/Axway/agent-sdk/pkg/transaction/util"
)

// EventMapper -
type EventMapper struct {
}

func (m *EventMapper) processMapping(gatewayTrafficLogEntry GwTrafficLogEntry) (*transaction.LogEvent, []transaction.LogEvent, error) {
	centralCfg := agent.GetCentralConfig()

	eventTimestamp, _ := time.Parse(time.RFC3339, gatewayTrafficLogEntry.EventTimestamp)
	eventTime := eventTimestamp.UTC().UnixNano() / int64(time.Millisecond)
	//eventTime := time.Now().UTC().Format(gatewayTrafficLogEntry.EventTimestamp)
	txID := gatewayTrafficLogEntry.Uuid
	txEventID := gatewayTrafficLogEntry.CorrelationId
	transInboundLogEventLeg, err := m.createTransactionEvent(eventTime, txID, gatewayTrafficLogEntry, txEventID, "Inbound")
	if err != nil {
		return nil, nil, err
	}

	transSummaryLogEvent, err := m.createSummaryEvent(eventTime, txID, gatewayTrafficLogEntry, centralCfg.GetTeamID())
	if err != nil {
		return nil, nil, err
	}

	return transSummaryLogEvent, []transaction.LogEvent{
		*transInboundLogEventLeg,
	}, nil
}

func (m *EventMapper) getTransactionEventStatus(code string) transaction.TxEventStatus {
	if code != "SUCCESS" {
		return transaction.TxEventStatusFail
	}
	return transaction.TxEventStatusPass
}

func (m *EventMapper) getHttpStatusCode(code string) int {
	if code == "SUCCESS" {
		return 200
	}
	return 500
}

func (m *EventMapper) getTransactionSummaryStatus(statusCode int) transaction.TxSummaryStatus {
	transSummaryStatus := transaction.TxSummaryStatusUnknown
	if statusCode >= http.StatusOK && statusCode < http.StatusBadRequest {
		transSummaryStatus = transaction.TxSummaryStatusSuccess
	} else if statusCode >= http.StatusBadRequest && statusCode < http.StatusInternalServerError {
		transSummaryStatus = transaction.TxSummaryStatusFailure
	} else if statusCode >= http.StatusInternalServerError && statusCode < http.StatusNetworkAuthenticationRequired {
		transSummaryStatus = transaction.TxSummaryStatusException
	}
	return transSummaryStatus
}

func (m *EventMapper) createTransactionEvent(eventTime int64, txID string, txDetails GwTrafficLogEntry, eventID, direction string) (*transaction.LogEvent, error) {

	httpStatus := m.getHttpStatusCode(txDetails.RequestStatus)

	host := txDetails.ServerId
	port := 443
	if strings.Index(host, ":") != -1 {
		uris := strings.Split(host, ":")
		host = uris[0]
		port, _ = strconv.Atoi(uris[1])

	}

	httpProtocolDetails, err := transaction.NewHTTPProtocolBuilder().
		SetURI(txDetails.OperationName).
		SetMethod(txDetails.NativeHttpMethod).
		SetStatus(httpStatus, http.StatusText(httpStatus)).
		SetHost(host).
		SetLocalAddress(host, port).
		Build()
	if err != nil {
		return nil, err
	}

	return transaction.NewTransactionEventBuilder().
		SetTimestamp(eventTime).
		SetTransactionID(txID).
		SetID(eventID).
		SetSource(txDetails.ServerId).
		SetDirection(direction).
		SetStatus(m.getTransactionEventStatus(txDetails.RequestStatus)).
		SetProtocolDetail(httpProtocolDetails).
		Build()
}

func (m *EventMapper) createSummaryEvent(eventTime int64, txID string, gatewayTrafficLogEntry GwTrafficLogEntry, teamID string) (*transaction.LogEvent, error) {
	statusCode := m.getHttpStatusCode(gatewayTrafficLogEntry.RequestStatus)
	method := gatewayTrafficLogEntry.NativeHttpMethod
	uri := gatewayTrafficLogEntry.OperationName
	host := gatewayTrafficLogEntry.ApplicationIp

	builder := transaction.NewTransactionSummaryBuilder().
		SetTimestamp(eventTime).
		SetTransactionID(txID).
		SetStatus(m.getTransactionSummaryStatus(statusCode), strconv.Itoa(statusCode)).
		SetTeam(teamID).
		SetDuration(gatewayTrafficLogEntry.TotalTime).
		SetEntryPoint("http", method, uri, host).
		// If the API is published to Central as unified catalog item/API service, se the Proxy details with the API definition
		// The Proxy.Name represents the name of the API
		// The Proxy.ID should be of format "remoteApiId_<ID Of the API on remote gateway>". Use transaction.FormatProxyID(<ID Of the API on remote gateway>) to get the formatted value.
		SetProxy(transutil.FormatProxyID(gatewayTrafficLogEntry.ApiName), gatewayTrafficLogEntry.ApiName, 0)

	if gatewayTrafficLogEntry.ApplicationName != "Unknown" && gatewayTrafficLogEntry.ApplicationId != "Unknown" {
		builder.SetApplication(gatewayTrafficLogEntry.ApplicationId, gatewayTrafficLogEntry.ApplicationName)
	}

	return builder.Build()

}
