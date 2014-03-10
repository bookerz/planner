package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"runtime"
)

type Config struct {
	DBUser         string
	DBPasswd       string
	DBMaxIdleConns int
	DBMaxOpenConns int
	Concurrency    int
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

func LoadConfig(file string) (*Config, error) {
	b, err := ioutil.ReadFile(file)

	if err != nil {
		log.Printf("Unable to load config file '%v', reason -> %v\n", configFile, err)
		return nil, err
	}

	c := &Config{}

	err = json.Unmarshal(b, c)

	if err != nil {
		log.Printf("Unable to parse config file '%v'\n", err)
		return nil, err
	}

	return c, nil
}

func ExampleConfig() string {
	cfg := Config{
		DBUser:   "dbuser",
		DBPasswd: "secret",
	}

	b, _ := json.MarshalIndent(cfg, "", " ")

	return string(b)
}
