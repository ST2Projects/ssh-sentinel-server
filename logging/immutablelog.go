package logging

import (
	"archive/tar"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/st2projects/ssh-sentinel-server/crypto"
	"io"
	"os"
	"sync"
	"time"
)

// TODO provide via config
const logPattern = "ssh-sentinel-%d.json"

var wg sync.WaitGroup
var predecessorEntryHash = "nil"
var log = Log{}

func Info(msg string) {
	wg.Add(1)
	go func() {
		tarLog, created := newLog()
		defer wg.Done()
		defer tarLog.Close()

		if !created {
			_, err := tarLog.Seek(-1024, io.SeekEnd)
			if err != nil {
				panic("Cannot seek end of tar " + err.Error())
			}
		}

		entry := newEntry(msg)
		entryBytes, _ := marshall(entry)

		tarHdr := newTarHeader(int64(len(entryBytes)))

		fmt.Printf("New entry %v\n", tarHdr)
		tarWriter := tar.NewWriter(tarLog)

		err := tarWriter.WriteHeader(tarHdr)
		if err != nil {
			panic("Cannot write header " + err.Error())
		}

		_, err = tarWriter.Write(entryBytes)
		if err != nil {
			panic("Cannot write entry " + err.Error())
		}
		err = tarWriter.Close()
		if err != nil {
			panic(err)
		}

		// Update predessor with this ID
		predecessorEntryHash = crypto.Sha256sum(entryBytes)

	}()
	wg.Wait()
}

func newTarHeader(len int64) *tar.Header {
	return &tar.Header{
		Name: fmt.Sprintf(logPattern, time.Now().UTC().UnixNano()),
		Size: len,
		Uid:  os.Getuid(),
		Gid:  os.Getgid(),
		Mode: 0700,
	}
}

func newEntry(msg string) *Entry {
	entryID := uuid.New().String()
	entry := &Entry{
		ID:              entryID,
		DateTime:        time.Now().UTC(),
		Level:           INFO,
		Event:           msg,
		PredecessorHash: predecessorEntryHash,
	}

	return entry
}

func Infof(msg string, args ...any) {
	// TODO
}

func marshall(entry *Entry) ([]byte, error) {
	return json.Marshal(entry)
}

//func makeFile() *os.File {
//
//	os.Mkdir("resources/log", 0700)
//
//	logfile, err := os.Create(fmt.Sprintf(logPattern, time.Now().UTC().UnixNano()))
//	if err != nil {
//		panic("Could not create log " + err.Error())
//	}
//	return logfile
//}

func newLog() (*os.File, bool) {
	logName := fmt.Sprintf("./resources/log/ssh-sentinel-%s.tar", time.Now().UTC().Format("2006-01-02"))

	created := false
	if _, err := os.Stat(logName); errors.Is(err, os.ErrNotExist) {
		_, err := os.Create(logName)
		if err != nil {
			return nil, false
		}
		created = true
	}

	tarFile, err := os.OpenFile(logName, os.O_RDWR, os.ModePerm)

	if err != nil {
		panic("Failed to open tar for writing " + err.Error())
	}

	return tarFile, created
}
