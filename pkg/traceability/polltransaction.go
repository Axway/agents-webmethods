package traceability

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/Axway/agent-sdk/pkg/cache"
	"github.com/Axway/agents-webmethods/pkg/config"
	"github.com/Axway/agents-webmethods/pkg/webmethods"

	"github.com/Axway/agent-sdk/pkg/jobs"

	"github.com/sirupsen/logrus"

	hc "github.com/Axway/agent-sdk/pkg/util/healthcheck"
	"github.com/Axway/agent-sdk/pkg/util/log"
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

type healthChecker func(name, endpoint string, check hc.CheckStatus) (string, error)

// WebmethodsEventEmitter - Gathers analytics data for publishing to Central.
type WebmethodsEventEmitter struct {
	client       webmethods.Client
	eventChannel chan WebmethodsEvent
	pollInterval time.Duration
	cache        cache.Cache
	cachePath    string
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
func NewWebmethodsEventEmitter(cachePath string, pollInterval time.Duration, eventChannel chan WebmethodsEvent, client webmethods.Client) *WebmethodsEventEmitter {
	me := &WebmethodsEventEmitter{
		eventChannel: eventChannel,
		client:       client,
		pollInterval: pollInterval,
	}
	me.cachePath = formatCachePath(cachePath)
	me.cache = cache.Load(me.cachePath)
	return me
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
	defer f.Close()
	return ioutil.ReadAll(f)
}

// Start retrieves analytics data from anypoint and sends them on the event channel for processing.
func (we *WebmethodsEventEmitter) Start() error {
	strStartTime, strEndTime := we.getLastRun()
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
	lastTime, err := time.Parse(dateFormat, strEndTime)
	if err != nil {
		logrus.WithFields(logrus.Fields{"strStartTime": strStartTime}).Warn("Unable to Parse Last Time")
	}
	for _, event := range events {

		if err != nil {
			log.Warnf("failed to marshal event: %s", err.Error())
		}
		we.eventChannel <- event
	}
	// Add 1 second to the last time stamp if we found records from this pull.
	// This will prevent duplicate records from being retrieved
	if len(events) > 0 {
		we.saveLastRun(lastTime.Add(time.Second * 1).Format(dateFormat))
	}
	return nil

}
func (we *WebmethodsEventEmitter) getLastRun() (string, string) {
	tStamp, _ := we.cache.Get(CacheKeyTimeStamp)
	//tNow := fmt.Sprintf("%d-%02d-%d %d:%d:%d", now.Year(), int(now.Month()), now.Day(), now.Hour(), now.Minute(), now.Second())
	//tNow := now.Format(dateFormat)
	var tNow string
	if tStamp == nil {
		now := time.Now()
		endTime := now.Add(we.pollInterval).Format(dateFormat)
		we.saveLastRun(endTime)
		tStamp = now.Format(dateFormat)
		tNow = endTime
	} else {
		lastTime, err := time.Parse(dateFormat, tStamp.(string))
		if err != nil {
			logrus.WithFields(logrus.Fields{"cacheLastTime": lastTime}).Warn("Unable to Parse Last Time")
		}
		tNow = lastTime.Add(we.pollInterval).Format(dateFormat)
	}
	return tStamp.(string), tNow
}
func (we *WebmethodsEventEmitter) saveLastRun(lastTime string) {
	we.cache.Set(CacheKeyTimeStamp, lastTime)
	we.cache.Save(we.cachePath)
}

// OnConfigChange passes the new config to the client to handle config changes
// since the MuleEventEmitter only has cache config value references and should not be changed
func (we *WebmethodsEventEmitter) OnConfigChange(gatewayCfg *config.AgentConfig) {
	we.client.OnConfigChange(gatewayCfg.WebMethodConfig)
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
func (m *WebmethodsEmitterJob) Start() error {
	jobID, err := jobs.RegisterIntervalJob(m, m.pollInterval)
	m.jobID = jobID
	return err
}

// OnConfigChange updates the MuleEventEmitterJob with any config changes, and calls OnConfigChange on the Emitter
func (m *WebmethodsEmitterJob) OnConfigChange(gatewayCfg *config.AgentConfig) {
	m.pollInterval = gatewayCfg.WebMethodConfig.PollInterval
	m.Emitter.OnConfigChange(gatewayCfg)
}

// Execute called by the sdk on each interval.
func (m *WebmethodsEmitterJob) Execute() error {
	return m.Emitter.Start()
}

// Status Performs a health check for this job before it is executed.
func (m *WebmethodsEmitterJob) Status() error {
	max := 3
	status := m.client.Healthcheck("health")
	if status.Result == hc.OK {
		m.consecutiveErrors = 0
	} else {
		m.consecutiveErrors++
	}
	if m.consecutiveErrors >= max {
		// If the job fails 3 times return an error
		return fmt.Errorf("failed to start the Traceability agent %d times in a row", max)
	}
	return nil
}

// Ready determines if the job is ready to run.
func (m *WebmethodsEmitterJob) Ready() bool {
	status := m.client.Healthcheck("health")
	if status.Result == hc.OK {
		return true
	}
	return false
}

func formatCachePath(path string) string {
	return fmt.Sprintf("%s/webmethods.cache", path)
}
