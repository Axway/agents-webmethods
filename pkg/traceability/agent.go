package traceability

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/elastic/beats/v7/filebeat/beater"
	fbcfg "github.com/elastic/beats/v7/filebeat/config"
	v2 "github.com/elastic/beats/v7/filebeat/input/v2"
	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/common"
	"github.com/elastic/beats/v7/libbeat/logp"

	"github.com/Axway/agent-sdk/pkg/traceability"
	localerrors "github.com/Axway/agents-webmethods/pkg/errors"
)

func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	err := validateInput(cfg)
	if err != nil {
		return nil, localerrors.ErrInvalidInputConfig.FormatError(err.Error())
	}

	agentCfg := GetConfig()
	if agentCfg == nil {
		return nil, localerrors.ErrConfigFile
	}
	eventProcessor := NewEventProcessor(agentCfg, traceability.GetMaxRetries())
	traceability.SetOutputEventProcessor(eventProcessor)

	factory := func(beat.Info, *logp.Logger, beater.StateStore) []v2.Plugin {
		return []v2.Plugin{}
	}

	// Initialize the filebeat to read events
	creator := beater.New(factory)
	return creator(b, cfg)
}

func validateInput(cfg *common.Config) error {
	filebeatConfig := fbcfg.DefaultConfig
	err := cfg.Unpack(&filebeatConfig)
	if err != nil {
		return err
	}

	if len(filebeatConfig.Inputs) == 0 {
		return errors.New("no inputs configured")
	}

	inputsEnabled := 0
	for _, input := range filebeatConfig.Inputs {
		inputConfig := struct {
			Enabled bool     `config:"enabled"`
			Paths   []string `config:"paths"`
		}{}
		input.Unpack(&inputConfig)
		if inputConfig.Enabled {
			inputsEnabled++
			err = validateInputPaths(inputConfig.Paths)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func validateInputPaths(paths []string) error {
	foundPath := false
	for _, path := range paths {
		path = strings.TrimSpace(path)
		if path != "" {
			parentDir := filepath.Dir(path)
			fileInfo, err := os.Stat(parentDir)
			if err != nil {
				return err
			}
			if !fileInfo.IsDir() {
				return errors.New("invalid path " + path)
			}
			foundPath = true
		}
	}
	if !foundPath {
		return errors.New("no paths were defined for input processing")
	}
	return nil
}
