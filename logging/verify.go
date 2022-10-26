package logging

import (
	"archive/tar"
	"encoding/json"
	"fmt"
	"github.com/st2projects/ssh-sentinel-server/crypto"
	"io"
	"os"
	"time"
)

var pred = &Entry{}

func Verify() {
	logName := fmt.Sprintf("./resources/log/ssh-sentinel-%s.tar", time.Now().UTC().Format("2006-01-02"))
	tarFile, _ := os.OpenFile(logName, os.O_RDONLY, os.ModePerm)
	defer tarFile.Close()

	reader := tar.NewReader(tarFile)

	for {
		header, err := reader.Next()

		if err == io.EOF {
			fmt.Println("EOF")
			break
		} else if err != nil {
			panic(err)
		}

		buff := make([]byte, header.Size)
		_, _ = io.ReadFull(reader, buff)

		msg := &Entry{}

		json.Unmarshal(buff, msg)

		if msg.PredecessorHash == "nil" {
			fmt.Printf("First event H = nill\n")
		} else {
			messagePred := msg.PredecessorHash
			b, _ := json.Marshal(pred)
			predecessorEntryHash := crypto.Sha256sum(b)

			fmt.Printf("Hash matches = %v -> %s == %s\n", messagePred == predecessorEntryHash, messagePred, predecessorEntryHash)
		}

		pred = msg

		fmt.Printf("M = %s, H = %s\n", msg.Event, msg.PredecessorHash)
	}
}
