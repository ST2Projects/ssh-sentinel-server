package config

type Config struct {
	CAPrivateKey string `json:"CAPrivateKey"`
	CAPublicKey  string `json:"CAPublicKey"`
	MaxValidTime string `json:"MaxValidTime"`
}
