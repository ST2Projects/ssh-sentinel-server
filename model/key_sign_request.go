package model

type KeySignRequest struct {
	APIKey     string   `json:"api_key"`
	Principals []string `json:"principals"`
	Key        string   `json:"key"`
}
