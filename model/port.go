package model

type HTTPConfig struct {
	HttpPort  int
	HttpsPort int
}

func (h HTTPConfig) Default() *HTTPConfig {
	return &HTTPConfig{
		HttpPort:  80,
		HttpsPort: 443,
	}
}
