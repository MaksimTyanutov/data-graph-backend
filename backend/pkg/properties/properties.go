package properties

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

var CompanyIdShift int

type Config struct {
	DbSettings struct {
		DbName     string `yaml:"dbName"`
		DbPort     string `yaml:"dbPort"`
		DbHost     string `yaml:"dbHost"`
		DbUsername string `yaml:"dbUsername"`
		DbPassword string `yaml:"dbPassword"`
	} `yaml:"DBSettings"`
	ProgramSettings struct {
		Host    string `yaml:"host"`
		Port    string `yaml:"port"`
		LogPath string `yaml:"logPath"`
	} `yaml:"ProgramSettings"`
}

func GetConfig(path string) (*Config, error) {
	if path == "" {
		return nil, errors.New("GetConfig error: path is nil")
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.New("GetConfig error: " + err.Error())
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Panic("Something wrong with config: " + err.Error())
		}
	}(file)

	reader := bufio.NewReader(file)
	data, err := io.ReadAll(reader)
	if err != nil {
		panic(err)
	}

	var properties Config
	if err := yaml.Unmarshal(data, &properties); err != nil {
		panic(err)
	}

	return &properties, nil
}
