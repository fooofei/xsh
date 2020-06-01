package sh

import (
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

	if err := checkPath4Download(s.Local); err != nil {
		response.Err = err
		return response
	}

	res, err := s.Copy.Download(s.Local, s.Remote)
	response.Err = err
	response.Status = res

	return response
}

func checkPath4Download(local string) error {
	if isLocalFileExist(local) {
		return LocalFileExistErr
	}

	makeLocalDir(local)
	if XConfig.Copy.DirEmptyCheck && !isLocalDirEmpty(local) {
		return LocalDirNotEmptyErr
	}

	return nil
}

func checkPath4Upload(session sshSession, local, remote string) error {
	target := filepath.Base(local)

	if strings.HasSuffix(local, string(os.PathSeparator)) {
		if isRemoteFileExist(session, remote+target) {
			return RemoteFileExistErr
		}
	} else {
		if isRemoteDirExist(session, remote+target) {
			return RemoteDirExistErr
		}
	}

	makeRemoteDir(session, remote+target)
	if XConfig.Copy.DirEmptyCheck && !isRemoteDirEmpty(session, remote+target) {
		return RemoteDirNotEmptyErr
	}

	return nil
}
