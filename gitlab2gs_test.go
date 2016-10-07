package main

import (
	"testing"
)

func TestloadConfig(t *testing.T) {
	config := loadConfig("./config.json.sample")
	if config.GitlabUser != "username" {
		t.Errorf("load config file error")
	}

	if config.GitlabPassword != "password" {
		t.Errorf("load config file error")
	}
}
