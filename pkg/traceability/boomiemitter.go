package traceability

import (
	"github.com/hpcloud/tail"
)

const (
	healthCheckEndpoint = "ingestion"
	CacheKeyTimeStamp   = "LAST_RUN"
)

type Emitter interface {
	Start() error
}

// BoomiEventEmitter - Gathers analytics data for publishing to Central.
type BoomiEventEmitter struct {
	eventChannel chan string
	logFile      string
}

// NewBoomiEventEmitter - Creates a client to poll for events.
func NewBoomiEventEmitter(logFile string, eventChannel chan string) *BoomiEventEmitter {
	me := &BoomiEventEmitter{
		eventChannel: eventChannel,
		logFile:      logFile,
	}
	return me
}

// Start retrieves analytics data from anypoint and sends them on the event channel for processing.
func (me *BoomiEventEmitter) Start() error {
	go me.tailFile()

	return nil

}

func (me BoomiEventEmitter) tailFile() {
	t, _ := tail.TailFile(me.logFile, tail.Config{Follow: true})
	for line := range t.Lines {
		me.eventChannel <- line.Text
	}
}
