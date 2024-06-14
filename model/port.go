package model

type HTTPConfig struct {
	Port     int
	ListenOn string
}

func (h HTTPConfig) Default() *HTTPConfig {
	return &HTTPConfig{
		Port:     8080,
		ListenOn: "0.0.0.0",
	}
}
