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
	remote = strings.TrimSpace(remote)
	if strings.Contains(remote, "\\") {
		return RemotePathIllegalErr
	}
	if !strings.HasPrefix(remote, "/") {
		return RemotePathNotFullErr
	}

	local = strings.TrimSpace(local)
	if runtime.GOOS == "windows" {
		if strings.Contains(local, "/") {
			return LocalPathIllegalErr
		}

		if !strings.HasSuffix(filepath.VolumeName(local), ":") {
			return LocalPathNotFullErr
		}
	} else {
		if strings.Contains(local, "\\") {
			return LocalPathIllegalErr
		}

		if !strings.HasPrefix(local, "/") {
			return LocalPathNotFullErr
		}
	}

	return nil
}

func isLocalFileExist(file string) bool {
	target := file
	if strings.HasSuffix(file, string(os.PathSeparator)) {
		target = filepath.Dir(target)
	}
	if f, err := os.Stat(target); err != nil {
		return false
	} else {
		return !f.IsDir()
	}
}

func isLocalDirExist(path string) bool {
	target := path
	if strings.HasSuffix(path, string(os.PathSeparator)) {
		target = filepath.Dir(target)
	}
	if f, err := os.Stat(target); err != nil {
		return false
	} else {
		return f.IsDir()
	}
}

func makeLocalDir(path string) {
	if err := os.MkdirAll(path, 0755); err != nil {
		Warn.Printf("make local dir error: %s", err.Error())
	}
}

func isLocalDirEmpty(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		Warn.Printf("open local dir error: %s", err.Error())
		return false
	}
	defer f.Close()

	if _, err = f.Readdirnames(1); err != io.EOF {
		return false
	}

	return true
}

func isRemoteFileExist(session sshSession, file string) bool {
	res := session.run(fmt.Sprintf("test -f '%s'", file))
	return res.Err == nil
}

func isRemoteDirExist(session sshSession, path string) bool {
	res := session.run(fmt.Sprintf("test -d '%s'", path))
	return res.Err == nil
}

func makeRemoteDir(session sshSession, path string) {
	_ = session.run(fmt.Sprintf("mkdir -p '%s'", path))
}

func isRemoteDirEmpty(session sshSession, path string) bool {
	res := session.run(fmt.Sprintf("ls -A '%s'", path))
	return res.Stdout == ""
}
