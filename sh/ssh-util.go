package sh

import (
	"fmt"
	. "github.com/xied5531/xsh/base"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func withTimeout(fn func(interface{}) SshResponse, arg interface{}, timeout time.Duration) SshResponse {
	result := make(chan SshResponse, 1)
	go func() {
		defer close(result)
		result <- fn(arg)
	}()

	select {
	case r := <-result:
		return r
	case <-time.After(timeout):
		return SshResponse{Err: TimeoutErr}
	}
}

func checkCommands(commands []string) error {
	if len(XBlackCommandRegexps) > 0 {
		for _, command := range commands {
			for _, reg := range XBlackCommandRegexps {
				if cmd := reg.FindStringIndex(command); cmd != nil {
					return fmt.Errorf("command[%s] in black command list", command)
				}
			}
		}
	}

	return nil
}

func GetLocalPath(direction string, path string) (string, error) {
	path = strings.Trim(path, " ")
	if filepath.VolumeName(path) == "" {
		return "", fmt.Errorf("local path must be full path")
	}

	p, e := filepath.Abs(path)
	if e != nil {
		return "", e
	}

	if direction == "download" && !strings.HasSuffix(path, string(os.PathSeparator)) {
		p = p + string(os.PathSeparator)
	}

	return p, nil
}

func GetRemotePath(direction string, path string) (string, error) {
	path = strings.Trim(path, " ")
	if !strings.HasPrefix(path, "/") {
		return "", fmt.Errorf("remote path must be full path")
	}

	if direction == "upload" && !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	return path, nil
}
