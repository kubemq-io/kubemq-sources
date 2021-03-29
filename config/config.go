package config

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/ghodss/yaml"
	"github.com/kubemq-hub/kubemq-sources/pkg/logger"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const defaultApiPort = 8082

var configFile string
var logr = logger.NewLogger("config")
var lastConf *Config
var defaultConfig = &Config{
	Bindings: []BindingConfig{},
	ApiPort:  defaultApiPort,
	LogLevel: "info",
}

type Config struct {
	Bindings []BindingConfig `json:"bindings"`
	ApiPort  int             `json:"apiPort"`
	LogLevel string          `json:"logLevel"`
}

func SetConfigFile(filename string) {
	configFile = filename
}
func (c *Config) hash() string {
	b, err := json.Marshal(c)
	if err != nil {
		return ""
	}
	h := sha256.New()
	_, _ = h.Write(b)
	hash := hex.EncodeToString(h.Sum(nil))
	return hash
}
func (c *Config) copy() *Config {
	b, _ := json.Marshal(c)
	n := &Config{}
	_ = json.Unmarshal(b, n)
	return n
}
func (c *Config) Validate() error {
	if c.ApiPort == 0 {
		c.ApiPort = defaultApiPort
	}
	exitedBindings := map[string]string{}
	for _, binding := range c.Bindings {
		if err := binding.Validate(); err != nil {
			return err
		}
		if _, ok := exitedBindings[binding.Name]; ok {
			return fmt.Errorf("duplicated binding names found: %s", binding.Name)
		} else {
			exitedBindings[binding.Name] = binding.Name
		}
	}
	return nil
}
func getConfigFormat(in []byte) (string, error) {
	c := &Config{}
	yamlErr := yaml.Unmarshal(in, c)
	if yamlErr == nil {
		return "yaml", nil
	}
	jsonErr := json.Unmarshal(in, c)
	if jsonErr == nil {
		return "json", nil
	}

	return "", fmt.Errorf("yaml parsing error: %s, json parsing error: %s", yamlErr.Error(), jsonErr.Error())
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
	fileExt, err := getConfigFormat(data)
	if fileExt == "" {
		return "", err
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
		fileExt, err := getConfigFormat([]byte(envConfigData))
		if fileExt == "" {
			return "", err
		}
		/* #nosec */
		err = ioutil.WriteFile("./config."+fileExt, []byte(envConfigData), 0644)
		if err != nil {
			return "", fmt.Errorf("cannot save environment config file")
		}
		return "./config." + fileExt, nil
	}
	return "", fmt.Errorf("no config data from environment variable")
}
func getConfigFile() (string, error) {
	if configFile != "" {
		loadedConfigFile, err := getConfigDataFromLocalFile(configFile)
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
func createDefaultConfig() (*Config, error) {
	data, err := yaml.Marshal(defaultConfig)
	if err != nil {
		return nil, err
	}
	/* #nosec */
	err = ioutil.WriteFile(configFile, data, 0644)
	if err != nil {
		return nil, err
	}
	return defaultConfig, nil
}
func load() (*Config, error) {
	loadedConfigFile, err := getConfigFile()
	if err != nil {
		return createDefaultConfig()
	} else {
		viper.SetConfigFile(loadedConfigFile)
	}
	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	err = viper.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}
	logr.Infof("%d bindings loaded", len(cfg.Bindings))
	logger.SetLogLevel(cfg.LogLevel)
	return cfg, err
}

func Load(cfgCh chan *Config) (*Config, error) {
	viper.AddConfigPath("./")
	cfg, err := load()
	if err != nil {
		return nil, err
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		cfg, err := load()
		if err != nil {
			logr.Errorf("error loading new configuration file: %s", err.Error())
			return
		}
		if cfg.hash() != lastConf.hash() {
			logr.Info("config file changed, reloading...")
			lastConf = cfg.copy()
			cfgCh <- cfg
		}
	})
	return cfg, err
}
