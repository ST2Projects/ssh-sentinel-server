package clientcli

import (
	"errors"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func ExpandPath(path string) string {
	usr, _ := user.Current()
	homeDir := usr.HomeDir

	if path == "~" {
		path = homeDir
	} else if strings.HasPrefix(path, "~/") {
		path = filepath.Join(homeDir, path[2:])
	}

	return path
}

func PathExists(path string) bool {
	_, err := os.Stat(path)

	return !(err != nil && errors.Is(err, os.ErrNotExist))
}
