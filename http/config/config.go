package config

import (
	"bufio"
	"errors"
	"io"
	"os"
	"perftest/http/logger"
	"strings"
)

var (
	log = logger.LOGGER
	ErrConfigPropertyNotExist = errors.New("config property not found")

	HTTP_REQUEST_CONFIG_PREFIX = "http.request"
)

type (
	GoTestConfig struct {
		configPath        string
		data map[string]interface{}
		httpRequestConfig map[string]interface{}
	}

	OnNotFound func (conf *GoTestConfig) interface{}
)

func NewConfig(configPath string) (conf *GoTestConfig) {
	conf = &GoTestConfig{
		configPath: configPath,
	}
	configFile, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("open local config file error, %s", err.Error())
		return
	}
	bufReader := bufio.NewReader(configFile)
	for line, _, err := bufReader.ReadLine(); err != io.EOF; {
		kv := strings.Split(string(line), "=")
		k := strings.TrimSpace(kv[0])
		v := strings.TrimSpace(kv[1])
		conf.data[k] = v

		if strings.HasPrefix(k, HTTP_REQUEST_CONFIG_PREFIX) {
			conf.httpRequestConfig[k] = v
		}
	}
	return
}

// Get get value from config
func (conf *GoTestConfig) Get(key string) interface{} {
	if conf.data == nil {
		return nil
	}
	return conf.data[key]
}

// Set set value to config
func (conf *GoTestConfig) Set(key string, value interface{}) {
	if conf.data == nil {
		conf.data = make(map[string]interface{})
	}
	conf.data[key] = value
}

// GetString get string value
func (conf *GoTestConfig) GetString(key string) string {
	v := conf.Get(key)
	if v == nil {
		return ""
	}
	str, _ := v.(string)
	return str
}

// GetBool get bool value
func (conf *GoTestConfig) GetBool(key string) bool {
	v := conf.Get(key)
	if v == nil {
		return false
	}
	b, _ := v.(bool)
	return b
}

// GetInt get int value
func (conf *GoTestConfig) GetInt(key string) int {
	v := conf.Get(key)
	if v == nil {
		return 0
	}
	i, _ := v.(int)
	return i
}
