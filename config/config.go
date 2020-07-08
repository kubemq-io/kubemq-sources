package config

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var configFile = pflag.String("config", "", "set config file name")

type Config struct {
	Sources  []Metadata
	Targets  []Metadata
	Bindings []Binding
}

func (c *Config) Validate() error {

	if len(c.Targets) == 0 {
		return fmt.Errorf("at least one target must be defined")
	}
	for i, target := range c.Targets {
		if err := target.Validate(); err != nil {
			return fmt.Errorf("target entry %d configuration error, %w", i, err)
		}
	}
	if len(c.Sources) == 0 {
		return fmt.Errorf("at least one source must be defined")
	}
	for i, source := range c.Sources {
		if err := source.Validate(); err != nil {
			return fmt.Errorf("source entry %d configuration error, %w", i, err)
		}
	}

	if len(c.Bindings) == 0 {
		return fmt.Errorf("at least one binding must be defined")
	}
	for i, binding := range c.Bindings {
		if err := binding.Validate(); err != nil {
			return fmt.Errorf("binding entry %d configuration error, %w", i, err)
		}
	}
	return nil
}
func getConfigFormat(in []byte) string {
	c := &Config{}
	err := yaml.Unmarshal(in, c)
	if err == nil {
		return "yaml"
	}
	err = json.Unmarshal(in, c)
	if err == nil {
		return "json"
	}
	return ""
}
func decodeBase64(in string) string {
	// base64 string cannot contain space so this is indication of base64 string
	if !strings.Contains(in, " ") {
		sDec, err := base64.StdEncoding.DecodeString(in)
		if err != nil {
			log.Println(fmt.Sprintf("error decoding config file base64 string: %s ", err.Error()))
			return in
		}
		return string(sDec)
	}
	return in
}
func getConfigDataFromLocalFile(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	fileExt := getConfigFormat(data)
	if fileExt == "" {
		return "", fmt.Errorf("invalid file format")
	}
	if strings.HasSuffix(filename, "."+fileExt) {
		return filename, nil
	}
	return filename + "." + fileExt, nil
}

func getConfigDataFromEnv() (string, error) {
	envConfigData, ok := os.LookupEnv("CONFIG")
	envConfigData = decodeBase64(envConfigData)
	if ok {
		fileExt := getConfigFormat([]byte(envConfigData))
		if fileExt == "" {
			return "", fmt.Errorf("invalid environment config format")
		}
		err := ioutil.WriteFile("./config."+fileExt, []byte(envConfigData), 0644)
		if err != nil {
			return "", fmt.Errorf("cannot save environment config file")
		}
		return "./config." + fileExt, nil
	}
	return "", fmt.Errorf("no config data from environment variable")
}
func getConfigFile() (string, error) {
	if *configFile != "" {
		loadedConfigFile, err := getConfigDataFromLocalFile(*configFile)
		if err != nil {
			return "", err
		}
		return loadedConfigFile, nil
	} else {
		loadedConfigFile, err := getConfigDataFromEnv()
		if err != nil {
			return "", err
		}
		return loadedConfigFile, nil
	}
}

func Load() (*Config, error) {
	pflag.Parse()
	viper.AddConfigPath("./")
	loadedConfigFile, err := getConfigFile()
	if err != nil {
		return nil, fmt.Errorf("error loading configuration, %w", err)
	} else {
		viper.SetConfigFile(loadedConfigFile)
	}
	fmt.Println(loadedConfigFile)
	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	err = viper.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, err
}
