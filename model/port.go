package model

type HTTPConfig struct {
	Port     int
	ListenOn string
}

func (h HTTPConfig) Default() *HTTPConfig {
	return &HTTPConfig{
		Port:     80,
		ListenOn: "0.0.0.0",
	}
}
