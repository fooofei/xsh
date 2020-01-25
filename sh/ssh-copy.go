package sh

import (
	"fmt"
	"github.com/pkg/sftp"
	. "github.com/xied5531/xsh/base"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type sshCopy struct {
	Session   copySession
	Direction string
	Local     string
	Remote    string
}

func (s sshCopy) run() sshResponse {
	var result sshResponse
	result = withTimeout(s.copy, s.Direction, time.Duration(XConfig.Timeout.CopyTimeoutS)*time.Second)
	result.Address = s.Session.Client.Host
	return result
}

func (s sshCopy) copy(arg interface{}) sshResponse {
	d := arg.(string)
	if d == "upload" {
		return s.upload()
	} else {
		return s.download()
	}
}

//remote must be directory
func (s sshCopy) upload() sshResponse {
	response := sshResponse{
		Status: make([]string, 0),
	}

	session, err := s.Session.newSftpSession()
	if err != nil {
		response.Err = err
		return response
	}

	localRoot := s.Local
	remote := s.Remote
	if e := session.MkdirAll(remote); e != nil {
		response.Err = fmt.Errorf("make remote directory[%s] error: %v", remote, e)
		return response
	}

	if err := filepath.Walk(s.Local, func(local string, info os.FileInfo, err error) error {
		s := local
		if err != nil {
			s = s + " :ERROR: " + err.Error()
		} else if info == nil {
			s = s + " :ERROR: file info nil"
		} else {
			if info.IsDir() {
				if local == localRoot {
					return nil
				}
				if e := session.MkdirAll(remote + strings.Replace(local, localRoot, "", 1)); e != nil {
					s = s + " :ERROR: " + e.Error()
				} else {
					s = s + " :OK"
				}
			} else if info.Mode().IsRegular() {
				remoteName := remote + strings.Replace(strings.Replace(local, localRoot, "", 1), "\\", "/", -1)
				if e := sftpUploadFile(session, local, remoteName); e != nil {
					if XConfig.Copy.Skip && e == RemoteFileExistErr {
						s = s + " :WARN: skip because exist"
					} else {
						s = s + " :ERROR: " + e.Error()
					}
				} else {
					s = s + " :OK"
				}
			} else {
				s = s + " :ERROR: file type not support"
			}
		}

		response.Status = append(response.Status, s)
		return nil
	}); err != nil {
		response.Err = err
	}

	return response
}

//local must be directory
func (s sshCopy) download() sshResponse {
	response := sshResponse{}

	if err := os.Mkdir(s.Local, os.ModeDir|0755); err != nil && !os.IsExist(err) {
		response.Err = err
		return response
	}

	local := s.Local + s.Session.Client.Host + string(os.PathSeparator)
	if err := os.Mkdir(local, os.ModeDir|0755); err != nil && !os.IsExist(err) {
		response.Err = err
		return response
	}

	session, err := s.Session.newSftpSession()
	if err != nil {
		response.Err = err
		return response
	}

	if strings.HasSuffix(s.Remote, "/") {
		status, e := s.downloadDir(session, local, s.Remote)
		if e != nil {
			response.Err = e
			response.Status = status
		} else {
			response.Status = status
		}
	} else {
		localName := local + filepath.Base(s.Remote)
		e := s.downloadFile(session, localName, s.Remote)
		if e != nil {
			response.Status = []string{localName + " <- " + s.Remote + " :ERROR: " + e.Error()}
		} else {
			response.Status = []string{localName + " <- " + s.Remote + " :OK"}
		}
	}

	return response
}

func (s sshCopy) downloadFile(session *sftp.Client, local, remote string) error {
	if _, err := os.Stat(local); err == nil {
		return LocalFileExistErr
	}

	stat, err := session.Stat(remote)
	if err != nil || !stat.Mode().IsRegular() {
		return fmt.Errorf("can only download regular file: %v", err)
	}

	rf, err := session.Open(remote)
	if err != nil {
		return err
	}
	defer rf.Close()

	lf, err := os.Create(local)
	if err != nil {
		return err
	}
	defer lf.Close()

	_, err = io.Copy(lf, rf)
	if err != nil {
		return err
	}

	return nil
}

func (s sshCopy) downloadDir(session *sftp.Client, local, remote string) ([]string, error) {
	status := make([]string, 0)

	w := session.Walk(s.Remote)
	if w == nil {
		return status, RemoteWalkErr
	}

	for w.Step() {
		stat := w.Stat()
		remoteName := w.Path()

		if remoteName == remote {
			continue
		}

		if stat.Mode().IsRegular() {
			remoteTmp := strings.Replace(remoteName, remote, "", 1)
			remoteTmp = strings.Replace(remoteTmp, "\\", "+", -1)
			if runtime.GOOS == "windows" {
				remoteTmp = strings.Replace(remoteTmp, "/", "\\", -1)
			}

			localName := local + remoteTmp
			if e := s.downloadFile(session, localName, remoteName); e != nil {
				if !XConfig.Copy.Override && e == LocalFileExistErr {
					status = append(status, localName+" <- "+remoteName+" :WARN: skip because exist")
				} else {
					status = append(status, localName+" <- "+remoteName+" :ERROR: "+e.Error())
				}
			} else {
				status = append(status, localName+" <- "+remoteName+" :OK")
			}
		} else if stat.IsDir() {
			localName := local + strings.Replace(remoteName, remote, "", 1)
			e := os.Mkdir(localName, os.ModeDir|0755)
			if e != nil && !os.IsExist(e) {
				status = append(status, localName+" <- "+remoteName+" :ERROR: "+e.Error())
			} else {
				status = append(status, localName+" <- "+remoteName+" :OK")
			}
		} else {
			status = append(status, remoteName+" :ERROR: file type not support")
		}
	}

	return status, nil
}

type copySession struct {
	Client sshClient
}

func (s copySession) newSftpSession() (*sftp.Client, error) {
	client, err := s.Client.NewClient()
	if err != nil {
		return nil, err
	}

	session, err := sftp.NewClient(client, sftp.MaxPacket(XConfig.Copy.SftpMaxPacketSize))
	if err != nil {
		return nil, err
	}

	return session, nil
}

func sftpUploadFile(client *sftp.Client, local string, remote string) error {
	if GOOS == "windows" {
		local = strings.Replace(local, "/", "+", -1)
		local = strings.Replace(local, "\\", "/", -1)
	}
	if !strings.HasSuffix(remote, "/") {
		remote = remote + "/"
	}

	if !XConfig.Copy.Override {
		if _, err := client.Stat(remote + local); err == nil {
			return RemoteFileExistErr
		}
	}

	stat, err := os.Stat(local)
	if err != nil || !stat.Mode().IsRegular() {
		return fmt.Errorf("can only upload regular file: %v", err)
	}

	rf, err := client.Create(remote + local)
	if err != nil {
		return err
	}
	defer rf.Close()

	lf, err := os.Open(local)
	if err != nil {
		return err
	}
	defer lf.Close()

	_, err = io.Copy(rf, lf)
	if err != nil {
		return err
	}

	return nil
}
