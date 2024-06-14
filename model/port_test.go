package model

import "testing"

func TestHTTPConfig_Default(t *testing.T) {
	config := HTTPConfig{}.Default()

	if config.Port != 8080 {
		t.Errorf("default port should be 8080 but got %d", config.Port)
	}

	if config.ListenOn != "0.0.0.0" {
		t.Error("default listenOn should be 0.0.0.0")
	}
}
