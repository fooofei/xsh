package sh

import (
	"fmt"
	. "github.com/xied5531/xsh/base"
	"io"
	"os"
	"path/filepath"
	"runtime"
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

func GetLocalDir(direction string, path string) (string, error) {
	path = strings.Trim(path, " ")
	if runtime.GOOS == "windows" {
		if strings.Contains(path, "/") {
			return "", LocalDirFormatIllegal
		}

		if !strings.HasSuffix(filepath.VolumeName(path), ":") {
			return "", LocalDirTypeIllegal
		}
	} else {
		if strings.Contains(path, "\\") {
			return "", LocalDirFormatIllegal
		}
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

func GetRemoteDir(direction string, path string) (string, error) {
	path = strings.Trim(path, " ")

	if strings.Contains(path, "\\") {
		return "", RemoteDirFormatIllegal
	}

	if !strings.HasPrefix(path, "/") {
		return "", RemoteDirTypeIllegal
	}

	if direction == "upload" && !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	return path, nil
}

func IsLocalDirEmpty(path string) bool {
	me := os.MkdirAll(path, 0755)
	f, err := os.Open(path)
	if err != nil || me != nil {
		return false
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	return err == io.EOF
}

func IsRemoteDirEmpty(session *xSshSession, path string) bool {
	stdout, _, _ := session.OutputStdoutStderr(fmt.Sprintf("mkdir -p \"%s\" ;ls -A \"%s\"", path, path))
	return stdout == ""
}

func CleanPath4Upload(path string) string {
	if runtime.GOOS == "windows" {
		return strings.Replace(path, "/", "+", -1)
	} else {
		return path
	}
}

func CleanPath4Download(path string) string {
	path = strings.Replace(path, "\\", "+", -1)
	path = strings.Replace(path, ":", "=", -1)
	return path
}
