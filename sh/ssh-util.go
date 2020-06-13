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

func isLocalFileExist(file string) (bool, error) {
	target := file
	if strings.HasSuffix(file, string(os.PathSeparator)) {
		target = filepath.Dir(target)
	}
	if f, err := os.Stat(target); err != nil {
		return false, err
	} else {
		return !f.IsDir(), fmt.Errorf("file can not be directory")
	}
}

func isLocalDirExist(path string) (bool, error) {
	target := path
	if strings.HasSuffix(path, string(os.PathSeparator)) {
		target = filepath.Dir(target)
	}
	if f, err := os.Stat(target); err != nil {
		return false, err
	} else {
		return f.IsDir(), fmt.Errorf("path must be directory")
	}
}

func makeLocalDir(path string) (bool, error) {
	if err := os.MkdirAll(path, 0755); err != nil {
		Warn.Printf("make local dir error: %s", err.Error())
		return false, err
	}
	return true, nil
}

func isLocalDirEmpty(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		Warn.Printf("open local dir error: %s", err.Error())
		return false, err
	}
	defer f.Close()

	if _, err = f.Readdirnames(1); err != io.EOF {
		return false, err
	}

	return true, nil
}

func isRemoteFileExist(session sshSession, file string) (bool, error) {
	res := session.run(fmt.Sprintf("test -f '%s'", file))
	if res.Err != nil {
		Warn.Printf("isRemoteFileExist res: %+v", res)
		return false, res.Err
	}
	return true, nil
}

func isRemoteDirExist(session sshSession, path string) (bool, error) {
	res := session.run(fmt.Sprintf("test -d '%s'", path))
	if res.Err != nil {
		Warn.Printf("isRemoteDirExist res: %+v", res)
		return false, res.Err
	}
	return true, nil
}

func makeRemoteDir(session sshSession, path string) (bool, error) {
	res := session.run(fmt.Sprintf("mkdir -p '%s'", path))
	if res.Err != nil {
		Warn.Printf("makeRemoteDir res: %+v", res)
		return false, res.Err
	}
	return true, nil
}

func isRemoteDirEmpty(session sshSession, path string) (bool, error) {
	res := session.run(fmt.Sprintf("ls -A '%s'", path))
	if res.Err != nil {
		Warn.Printf("isRemoteDirEmpty res: %+v", res)
		return false, res.Err
	}
	return res.Stdout == "", nil
}
