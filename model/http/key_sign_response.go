package http

type KeySignResponse struct {
	Success   bool   `json:"success"`
	Message   any    `json:"message"`
	SignedKey string `json:"signedKey"`
	NotBefore uint64 `json:"notBefore"`
	NotAfter  uint64 `json:"notAfter"`
}

func NewKeySignResponse(success bool, message any) *KeySignResponse {
	return &KeySignResponse{
		success,
		message,
		"",
		0,
		0}
}
