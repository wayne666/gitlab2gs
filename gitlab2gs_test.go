package main

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	config := loadConfig("./config.json.sample")
	if config.GitlabUser != "username" {
		t.Errorf("load config file error")
	}

	if config.GitlabPassword != "password" {
		t.Error("load config file error")
	}
}
