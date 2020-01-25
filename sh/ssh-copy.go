package sh

import (
	. "github.com/xied5531/xsh/base"
	"os"
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

	xs, err := sshSession{
		Client: s.SshClient,
	}.newSession()
	if err != nil {
		response.Err = err
		return response
	}
	defer xs.Close()

	if !IsRemoteDirEmpty(xs, s.Remote) && XConfig.Copy.PathEmptyCheck {
		response.Err = RemoteDirNotEmptyErr
		return response
	}

	if !strings.HasSuffix(s.Remote, "/") {
		s.Remote = s.Remote + "/"
	}

	res, err := s.Copy.Upload(s.Local, s.Remote)
	response.Err = err
	response.Status = res

	return response
}

//local must be directory
func (s sshCopy) download() SshResponse {
	response := SshResponse{}

	if !IsLocalDirEmpty(s.Local) && XConfig.Copy.PathEmptyCheck {
		response.Err = LocalDirNotEmptyErr
		return response
	}

	if !strings.HasSuffix(s.Local, string(os.PathSeparator)) {
		s.Local = s.Local + string(os.PathSeparator)
	}

	res, err := s.Copy.Download(s.Local, s.Remote)
	response.Err = err
	response.Status = res

	return response
}
