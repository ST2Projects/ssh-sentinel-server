package api

type KeySignRequest struct {
	Username   string      `json:"username"`
	APIKey     string      `json:"api_key"`
	Principals []string    `json:"principals"`
	Key        string      `json:"key"`
	Extensions []Extension `json:"extensions"`
}

type Extension string

const (
	no_touch_required       Extension = "no-touch-required"
	permit_x11_forwarding   Extension = "permit-x11-forwarding"
	permit_agent_forwarding Extension = "permit-agent_forwarding"
	permit_port_forwarding  Extension = "permit-port-forwarding"
	permit_pty              Extension = "permit-pty"
	permit_user_rc          Extension = "permit-user-rc"
)
