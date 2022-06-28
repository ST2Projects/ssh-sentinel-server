package http

type KeySignRequest struct {
	Username   string   `json:"username"`
	APIKey     string   `json:"api_key"`
	Principals []string `json:"principals"`
	Key        string   `json:"key"`
}
