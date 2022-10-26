package logging

import (
	"github.com/google/uuid"
	"time"
)

type Entry struct {
	ID              string    `json:"id"`
	DateTime        time.Time `json:"dateTime"`
	Level           Level     `json:"level"`
	Event           string    `json:"event"`
	PredecessorHash string    `json:"predecessorHash"`
}

type Level string

const (
	INFO  Level = "INFO"
	WARN  Level = "WARN"
	ERROR Level = "ERROR"
)

type Log struct {
	Entries []*Entry `json:"entries"`
}

func (e *Entry) New() (*Entry, string) {
	id := uuid.New().String()
	e.ID = id

	return e, id
}
