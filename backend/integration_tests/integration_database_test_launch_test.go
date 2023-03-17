package integration_tests

import (
	"data-graph-backend/pkg/dbConnector"
	"data-graph-backend/pkg/logging"
	"data-graph-backend/pkg/properties"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func Test_db_not_available_on_launch(t *testing.T) {
	var (
		_, b, _, _ = runtime.Caller(0)
		configPath = filepath.Dir(b) + "/../config/config_test.yaml"
	)

	config, err := properties.GetConfig(configPath)
	if err != nil {
		t.Errorf("Can't get config at the path: " + configPath)
	}

	logging.Init(config.ProgramSettings.LogPath)
	logger := logging.GetLogger()
	logger.Info("Starting backend for DataGraph")

	expected := "A connection attempt failed because the connected party did not properly respond after a period of time, or established connection failed because connected host has failed to respond"
	_, err = dbConnector.NewConnection(config, logger)
	if err == nil {
		t.Errorf("No error when failed connection (or database is okay)")
		return
	}
	if strings.Contains(err.Error(), expected) {
		t.Errorf("expected error to be %s \ngot %s", expected, err.Error())
	}
	return
}
