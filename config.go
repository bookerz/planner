package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	DBUser   string
	DBPasswd string
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
