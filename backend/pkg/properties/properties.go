package properties

import (
	"bufio"
	"fmt"
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
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"ProgramSettings"`
}

func GetConfig(path string) *Config {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Panicln("Something wrong with config: " + err.Error())
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

	return &properties
}
