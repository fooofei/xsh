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

func checkFullPath(local string, remote string) error {
	remote = strings.Trim(remote, " ")
	if strings.Contains(remote, "\\") {
		return RemoteDirFormatIllegal
	}
	if !strings.HasPrefix(remote, "/") {
		return RemoteDirTypeIllegal
	}

	local = strings.Trim(local, " ")
	if runtime.GOOS == "windows" {
		if strings.Contains(local, "/") {
			return LocalDirFormatIllegal
		}

		if !strings.HasSuffix(filepath.VolumeName(local), ":") {
			return LocalDirTypeIllegal
		}
	} else {
		if strings.Contains(local, "\\") {
			return LocalDirFormatIllegal
		}

		if !strings.HasSuffix(local, "/") {
			return LocalDirTypeIllegal
		}
	}

	return nil
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
