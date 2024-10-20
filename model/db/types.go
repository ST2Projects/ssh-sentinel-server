package db

type Event string

const (
	Login Event = "login"
	Sign  Event = "sign"
	Fetch Event = "fetch_ca"
)
