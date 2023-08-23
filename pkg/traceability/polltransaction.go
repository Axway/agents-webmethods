package traceability

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Axway/agent-sdk/pkg/util/log"
	"io"
	"time"

	"github.com/Axway/agent-sdk/pkg/cache"
	"github.com/Axway/agents-webmethods/pkg/config"
	"github.com/Axway/agents-webmethods/pkg/webmethods"

	"github.com/Axway/agent-sdk/pkg/jobs"

	"github.com/sirupsen/logrus"

	hc "github.com/Axway/agent-sdk/pkg/util/healthcheck"
)

const (
	CacheKeyTimeStamp = "LAST_RUN"
	dateFormat        = "2006-01-02 15:04:05"
)

type TraceCache struct {
}

type Emitter interface {
	Start() error
	OnConfigChange(gatewayCfg *config.AgentConfig)
}

// WebmethodsEventEmitter - Gathers analytics data for publishing to Central.
type WebmethodsEventEmitter struct {
	client           webmethods.Client
	eventChannel     chan WebmethodsEvent
	pollInterval     time.Duration
	cache            cache.Cache
	cachePath        string
	timezoneLocation time.Location
	analyticsDelay   time.Duration
}

// WebmethodsEmitterJob wraps an Emitter and implements the Job interface so that it can be executed by the sdk.
type WebmethodsEmitterJob struct {
	Emitter
	consecutiveErrors int
	jobID             string
	pollInterval      time.Duration
	client            webmethods.Client
}

// NewWebmethodsEventEmitter - Creates a client to poll for events.
func NewWebmethodsEventEmitter(agentConfig config.AgentConfig, eventChannel chan WebmethodsEvent, client webmethods.Client, timezoneLocation time.Location) *WebmethodsEventEmitter {
	we := &WebmethodsEventEmitter{
		eventChannel:     eventChannel,
		client:           client,
		pollInterval:     agentConfig.WebMethodConfig.PollInterval,
		timezoneLocation: timezoneLocation,
		analyticsDelay:   agentConfig.WebMethodConfig.AnalyticsDelay,
	}
	we.cachePath = formatCachePath(agentConfig.WebMethodConfig.CachePath)
	we.cache = cache.Load(we.cachePath)
	return we
}

func (we *WebmethodsEventEmitter) unzipFile(body []byte) ([]byte, error) {
	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		logrus.WithError(err).Error("Error extracting zip file")
		return nil, err
	}
	var unzippedFileBytes []byte
	for _, zipFile := range zipReader.File {
		logrus.Infof("Reading file: %s", zipFile.Name)
		unzippedFileBytes, err = readZipFile(zipFile)
		if err != nil {
			continue
		}
	}
	return unzippedFileBytes, nil

}

func readZipFile(zf *zip.File) ([]byte, error) {
	f, err := zf.Open()
	if err != nil {
		return nil, err
	}
	defer func(f io.ReadCloser) {
		err := f.Close()
		if err != nil {
			fmt.Printf("Error closing zip file %s", err)
		}
	}(f)
	return io.ReadAll(f)
}

// Start retrieves analytics data from webmethods and sends them on the event channel for processing.
func (we *WebmethodsEventEmitter) Start() error {
	strStartTime, strEndTime := we.getLastRun()
	logrus.Infof("Start time : %s End time :%s", strStartTime, strEndTime)
	data, err := we.client.GetTransactionsWindow(strStartTime, strEndTime)
	eventsBytes, err := we.unzipFile(data)

	if err != nil {
		logrus.WithError(err).Error("failed to unzip transactions data")
		return err
	}

	reader := bytes.NewReader(eventsBytes)
	scanner := bufio.NewScanner(reader)
	events := make([]WebmethodsEvent, 0)
	for scanner.Scan() {
		lineItem := scanner.Bytes()
		event := &WebmethodsEvent{}
		err = json.Unmarshal(lineItem, event)
		events = append(events, *event)
	}
	logrus.Infof("Total number of events retrieved  from Webmethods : %d", len(events))
	if err != nil {
		logrus.WithFields(logrus.Fields{"strStartTime": strStartTime}).Warn("Unable to Parse Last Time")
	}
	for _, event := range events {

		if err != nil {
			log.Warnf("failed to marshal event: %s", err.Error())
		}
		we.eventChannel <- event
	}
	we.saveLastRun(strEndTime)
	return nil
}

func (we *WebmethodsEventEmitter) getLastRun() (string, string) {
	tStamp, _ := we.cache.Get(CacheKeyTimeStamp)
	//tNow := fmt.Sprintf("%d-%02d-%d %d:%d:%d", now.Year(), int(now.Month()), now.Day(), now.Hour(), now.Minute(), now.Second())
	//tNow := now.Format(dateFormat)
	var tNow string
	now := time.Now()
	now = now.In(&we.timezoneLocation)
	if tStamp == nil {
		tStamp = now.Add(-we.analyticsDelay * 2).Format(dateFormat)
		tNow = now.Add(-we.analyticsDelay).Format(dateFormat)
	} else {
		tNow = now.Add(-we.analyticsDelay).Format(dateFormat)
	}
	return tStamp.(string), tNow
}
func (we *WebmethodsEventEmitter) saveLastRun(lastTime string) {
	err := we.cache.Set(CacheKeyTimeStamp, lastTime)
	if err != nil {
		log.Error("Failed to set value to cache")
	}
	err = we.cache.Save(we.cachePath)
	if err != nil {
		log.Error("Failed to save value to cache")
	}
}

// OnConfigChange passes the new config to the client to handle config changes
// since the MuleEventEmitter only has cache config value references and should not be changed
func (we *WebmethodsEventEmitter) OnConfigChange(gatewayCfg *config.AgentConfig) {
	err := we.client.OnConfigChange(gatewayCfg.WebMethodConfig)
	if err != nil {
		return
	}
}

// NewMuleEventEmitterJob creates a struct that implements the Emitter and Job interfaces.
func NewMuleEventEmitterJob(
	emitter Emitter,
	pollInterval time.Duration,
	client webmethods.Client,
) (*WebmethodsEmitterJob, error) {

	return &WebmethodsEmitterJob{
		Emitter:      emitter,
		pollInterval: pollInterval,
		client:       client,
	}, nil
}

// Start registers the job with the sdk.
func (w *WebmethodsEmitterJob) Start() error {
	jobID, err := jobs.RegisterIntervalJob(w, w.pollInterval)
	logrus.Trace("starting job")
	w.jobID = jobID
	return err
}

// OnConfigChange updates the MuleEventEmitterJob with any config changes, and calls OnConfigChange on the Emitter
func (w *WebmethodsEmitterJob) OnConfigChange(gatewayCfg *config.AgentConfig) {
	w.pollInterval = gatewayCfg.WebMethodConfig.PollInterval
	w.Emitter.OnConfigChange(gatewayCfg)
}

// Execute called by the sdk on each interval.
func (w *WebmethodsEmitterJob) Execute() error {
	return w.Emitter.Start()
}

// Status Performs a health check for this job before it is executed.
func (w *WebmethodsEmitterJob) Status() error {
	maxRetry := 3
	status := w.client.Healthcheck("health")
	if status.Result == hc.OK {
		w.consecutiveErrors = 0
	} else {
		w.consecutiveErrors++
	}
	if w.consecutiveErrors >= maxRetry {
		// If the job fails 3 times return an error
		return fmt.Errorf("failed to start the Traceability agent %d times in a row", maxRetry)
	}
	return nil
}

// Ready determines if the job is ready to run.
func (w *WebmethodsEmitterJob) Ready() bool {
	status := w.client.Healthcheck("health")
	if status.Result == hc.OK {
		return true
	}
	return false
}

func formatCachePath(path string) string {
	return fmt.Sprintf("%s/webmethods.cache", path)
}
