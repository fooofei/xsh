package sh

import (
	"fmt"
	"github.com/pkg/sftp"
	. "github.com/xied5531/xsh/base"
	"io"
	"os"
	"path/filepath"
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

func (s sshCopy) upload() sshResponse {
	response := sshResponse{
		Status: make([]string, 0),
	}

	session, err := s.Session.newSftpSession()
	if err != nil {
		response.Err = err
		return response
	}

	remote := s.Remote
	if e := sftpMakeDir(session, "", remote); e != nil {
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
				if local != "." && local != ".." {
					if e := sftpMakeDir(session, local, remote); e != nil {
						s = s + " :ERROR: " + e.Error()
					} else {
						s = s + " :OK"
					}
				}
			} else if info.Mode().IsRegular() {
				if e := sftpMakeFile(session, local, remote); e != nil {
					if XConfig.Copy.Skip && e == RemoteFileExistErr {
						s = s + " :WARN:skip file: " + e.Error()
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

func (s sshCopy) download() sshResponse {
	response := sshResponse{
		Status: make([]string, 0),
	}

	session, err := s.Session.newSftpSession()
	if err != nil {
		response.Err = err
		return response
	}

	local := s.Local
	if e := localMakeDir(local, "", ""); e != nil {
		response.Err = fmt.Errorf("make local directory[%s] error: %v", local, e)
		return response
	}

	local = local + string(os.PathSeparator) + s.Session.Client.Host
	if e := localMakeDir(local, "", ""); e != nil {
		response.Err = fmt.Errorf("make local directory[%s] error: %v", local, e)
		return response
	}

	w := session.Walk(s.Remote)
	if w == nil {
		response.Err = RemoteWalkErr
		return response
	}

	for w.Step() {
		stat := w.Stat()
		remote := w.Path()
		res := remote
		if stat.IsDir() {
			if remote == s.Remote {
				continue
			}
			if e := localMakeDir(local, remote, s.Remote); e != nil {
				res = res + " :ERROR: " + e.Error()
			} else {
				res = res + " :OK"
			}
		} else if stat.Mode().IsRegular() {
			if e := localMakeFile(session, local, remote, s.Remote); e != nil {
				if XConfig.Copy.Skip && e == LocalFileExistErr {
					res = res + " :WARN:skip file: " + e.Error()
				} else {
					res = res + " :ERROR: " + e.Error()
				}
			} else {
				res = res + " :OK"
			}
		} else {
			res = res + " :ERROR: file type not support"
		}

		response.Status = append(response.Status, res)
	}

	return response
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

func localMakeDir(local string, remote string, prefix string) error {
	if prefix != "" {
		remote = strings.Replace(remote, prefix, "", 1)
	}
	if strings.HasPrefix(remote, "/") {
		remote = strings.Replace(remote, "/", "", 1)
	}

	if GOOS == "windows" {
		remote = strings.Replace(remote, "\\", "+", -1)
		remote = strings.Replace(remote, "/", "\\", -1)
	}

	if !strings.HasSuffix(local, string(os.PathSeparator)) {
		local = local + string(os.PathSeparator)
	}

	err := os.Mkdir(local+remote, os.ModeDir|0755)
	if !os.IsExist(err) {
		return err
	}

	return nil
}

func localMakeFile(client *sftp.Client, local string, remote string, prefix string) error {
	oldRemote := remote
	if prefix != "" {
		remote = strings.Replace(remote, prefix, "", 1)
	}
	if strings.HasPrefix(remote, "/") {
		remote = strings.Replace(remote, "/", "", 1)
	}

	if GOOS == "windows" {
		remote = strings.Replace(remote, "\\", "+", -1)
		remote = strings.Replace(remote, "/", "\\", -1)
	}

	if !strings.HasSuffix(local, string(os.PathSeparator)) {
		local = local + string(os.PathSeparator)
	}

	if !XConfig.Copy.Override {
		if _, err := os.Stat(local + remote); err == nil {
			return LocalFileExistErr
		}
	}

	rf, err := client.Open(oldRemote)
	if err != nil {
		return err
	}
	defer rf.Close()

	lf, err := os.Create(local + remote)
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

func sftpMakeDir(client *sftp.Client, local string, remote string) error {
	if GOOS == "windows" {
		local = strings.Replace(local, "/", "+", -1)
		local = strings.Replace(local, "\\", "/", -1)
	}
	if !strings.HasSuffix(remote, "/") {
		remote = remote + "/"
	}

	return client.MkdirAll(remote + local)
}

func sftpMakeFile(client *sftp.Client, local string, remote string) error {
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
