package main

import (
	"encoding/json"
	log "github.com/golang/glog"
	"io/ioutil"
	"runtime"
)

type Config struct {
	DBUser         string
	DBPasswd       string
	DBMaxIdleConns int
	DBMaxOpenConns int
	Concurrency    int
	WebRoot        string
}

func (c *Config) getMaxIdleConns() int {
	if c.DBMaxIdleConns == 0 {
		return 10
	}
	return c.DBMaxIdleConns
}

func (c *Config) getMaxOpenConns() int {
	if c.DBMaxOpenConns == 0 {
		return 10
	}
	return c.DBMaxOpenConns
}

func (c *Config) getConcurrency() int {
	if c.Concurrency == 0 {
		return runtime.NumCPU()
	}
	return c.Concurrency
}

func (c *Config) getWebRoot() string {
	return c.WebRoot
}

func LoadConfig(file string) (*Config, error) {
	b, err := ioutil.ReadFile(file)

	if err != nil {
		log.Warningf("Unable to load config file '%v', reason -> %v\n", configFile, err)
		return nil, err
	}

	c := &Config{}

	err = json.Unmarshal(b, c)

	if err != nil {
		log.Warningf("Unable to parse config file '%v'\n", err)
		return nil, err
	}

	return c, nil
}

func ExampleConfig() string {
	cfg := Config{
		DBUser:         "dbuser",
		DBPasswd:       "secret",
		DBMaxIdleConns: 10,
		DBMaxOpenConns: 10,
		Concurrency:    10,
		WebRoot:        "./web/app/",
	}

	b, _ := json.MarshalIndent(cfg, "", " ")

	return string(b)
}
