package properties

import (
	"errors"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

var (
	_, b, _, _  = runtime.Caller(0)
	projectPath = filepath.Dir(filepath.Dir(filepath.Dir(b)))
)

func Test_GetConfig(t *testing.T) {
	testTable := []struct {
		name          string
		path          string
		expectedError error
	}{
		{
			name:          "Correct config path",
			path:          projectPath + "/config/config_example.yaml",
			expectedError: nil,
		},
		{
			name:          "Incorrect path",
			path:          "./configuration.yaml",
			expectedError: errors.New("GetConfig error: open ./configuration.yaml:"),
		},
		{
			name:          "Nil path",
			path:          "",
			expectedError: errors.New("GetConfig error: path is nil"),
		},
	}
	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			config, err := GetConfig(test.path)
			if err != nil && test.expectedError == nil {
				t.Errorf("Actual error = %s\nExpected error = nil", err.Error())
				return
			}
			if err == nil && test.expectedError != nil {
				t.Errorf("Actual error = nil\nExpected error = %s", test.expectedError)
				return
			}
			if test.expectedError == nil {
				if config == nil {
					t.Errorf("Config is nil\n Error: %s", err.Error())
					return
				}
			}
			if err != nil && test.expectedError != nil {
				if !strings.Contains(err.Error(), test.expectedError.Error()) {
					t.Errorf("Actual error = %s\nExpected error = %s", err.Error(), test.expectedError)
					return
				}
			}
		})
	}
}
