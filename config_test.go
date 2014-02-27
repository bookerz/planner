package main

import (
	"testing"
)

var testResourcesDir string = "test-resources/"

func TestConfig(t *testing.T) {
	file := testResourcesDir + "config.json"

	c, err := LoadConfig(file)

	if err != nil {
		t.Fail()
	}

	if c == nil {
		t.Error("c can not be nil hear ")
	}

	user := "user"

	if c.DBUser != user {
		t.Errorf("expected %v but got %v", user, c.DBUser)
	}

	passwd := "password"

	if c.DBPasswd != passwd {
		t.Errorf("expected %v but got %v", passwd, c.DBPasswd)
	}

}
