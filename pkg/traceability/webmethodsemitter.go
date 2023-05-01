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

// WebmethodsEventEmitter - Gathers analytics data for publishing to Central.
type WebmethodsEventEmitter struct {
	eventChannel chan string
	logFile      string
}

// NewWebmethodsEventEmitter - Creates a client to poll for events.
func NewWebmethodsEventEmitter(logFile string, eventChannel chan string) *WebmethodsEventEmitter {
	me := &WebmethodsEventEmitter{
		eventChannel: eventChannel,
		logFile:      logFile,
	}
	return me
}

// Start retrieves analytics data from anypoint and sends them on the event channel for processing.
func (me *WebmethodsEventEmitter) Start() error {
	go me.tailFile()

	return nil

}

func (me WebmethodsEventEmitter) tailFile() {
	t, _ := tail.TailFile(me.logFile, tail.Config{Follow: true})
	for line := range t.Lines {
		me.eventChannel <- line.Text
	}
}
