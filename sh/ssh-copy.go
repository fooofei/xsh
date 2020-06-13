package sh

import (
	"fmt"
	. "github.com/xied5531/xsh/base"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type copyItf interface {
	Download(local, remote string) ([]string, error)
	Upload(local, remote string) ([]string, error)
}

type sshCopy struct {
	SshClient sshClient
	Direction string
	Local     string
	Remote    string
	Copy      copyItf
}

func (s sshCopy) run() SshResponse {
	var result SshResponse
	result = withTimeout(s.copy, s.Direction, time.Duration(XConfig.Timeout.CopyTimeoutS)*time.Second)
	result.Address = s.SshClient.Host
	return result
}

func (s sshCopy) copy(arg interface{}) SshResponse {
	s.Local = strings.TrimSpace(s.Local)
	s.Remote = strings.TrimSpace(s.Remote)

	d := arg.(string)
	if d == "upload" {
		return s.upload()
	} else {
		return s.download()
	}
}

//remote must be directory
func (s sshCopy) upload() SshResponse {
	response := SshResponse{}

	if !strings.HasSuffix(s.Remote, "/") {
		s.Remote = s.Remote + "/"
	}

	xs := sshSession{
		Client: s.SshClient,
	}
	if err := checkPath4Upload(xs, s.Local, s.Remote); err != nil {
		response.Err = err
		return response
	}

	res, err := s.Copy.Upload(s.Local, s.Remote)
	response.Err = err
	response.Status = res

	return response
}

//local must be directory
func (s sshCopy) download() SshResponse {
	response := SshResponse{}

	if !strings.HasSuffix(s.Local, string(os.PathSeparator)) {
		s.Local = s.Local + string(os.PathSeparator)
	}

	xs := sshSession{
		Client: s.SshClient,
	}
	if err := checkPath4Download(xs, s.Local, s.Remote); err != nil {
		response.Err = err
		return response
	}

	res, err := s.Copy.Download(s.Local, s.Remote)
	response.Err = err
	response.Status = res

	return response
}

func checkPath4Download(session sshSession, local, remote string) error {
	if strings.HasSuffix(remote, "/") {
		if ok, err := isRemoteDirExist(session, remote); !ok {
			return fmt.Errorf("remote dir not exist: %s, error: %v", remote, err)
		}
	} else {
		if ok, err := isRemoteFileExist(session, remote); !ok {
			return fmt.Errorf("remote file not exist: %s, error: %v", remote, err)
		}
	}

	if ok, err := isLocalFileExist(local); ok {
		return fmt.Errorf("local file exist: %s, error: %v", local, err)
	}

	if ok, err := makeLocalDir(local); !ok {
		return fmt.Errorf("make local dir error: %s, error: %v", local, err)
	}

	if XConfig.Copy.OverrideCheck {
		if ok, err := isLocalDirEmpty(local); !ok {
			return fmt.Errorf("local dir not empty: %s, error: %v", local, err)
		}
	}

	return nil
}

func checkPath4Upload(session sshSession, local, remote string) error {
	target := filepath.Base(local)

	if strings.HasSuffix(local, string(os.PathSeparator)) {
		if ok, err := isLocalDirExist(local); !ok {
			return fmt.Errorf("local dir not exist: %s, error: %v", local, err)
		}

		if ok, err := isRemoteFileExist(session, remote+target); ok {
			return fmt.Errorf("remote file exist: %s, error: %v", remote+target, err)
		}

		if ok, err := makeRemoteDir(session, remote); !ok {
			return fmt.Errorf("make remote dir error: %s, error: %v", remote, err)
		}

		if XConfig.Copy.OverrideCheck {
			if ok, err := isRemoteDirEmpty(session, remote+target); !ok {
				return fmt.Errorf("remote dir not empty: %s, error: %v", remote+target, err)
			}
		}
	} else {
		if ok, err := isLocalFileExist(local); !ok {
			return fmt.Errorf("local file not exist: %s, error: %v", local, err)
		}

		if ok, err := isRemoteDirExist(session, remote+target); ok {
			return fmt.Errorf("remote dir exist: %s, error: %v", remote+target, err)
		}

		if ok, err := makeRemoteDir(session, remote); !ok {
			return fmt.Errorf("make remote dir error: %s, error: %v", remote, err)
		}

		if XConfig.Copy.OverrideCheck {
			if ok, err := isRemoteFileExist(session, remote+target); ok {
				return fmt.Errorf("remote file exist: %s, error: %v", remote+target, err)
			}
		}
	}

	return nil
}
