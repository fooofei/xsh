package sh

import (
	"github.com/pkg/sftp"
	. "github.com/xied5531/xsh/base"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type sftpCopy struct {
	SshClient sshClient
}

func (s sftpCopy) newSession() (*sftp.Client, error) {
	client, err := s.SshClient.NewClient()
	if err != nil {
		return nil, err
	}

	session, err := sftp.NewClient(client, sftp.MaxPacket(XConfig.Copy.SftpMaxPacketSize))
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s sftpCopy) Download(local, remote string) ([]string, error) {
	local = local + s.SshClient.Host + string(os.PathSeparator)
	if err := os.Mkdir(local, os.ModeDir|0755); err != nil && !os.IsExist(err) {
		return nil, err
	}

	session, err := s.newSession()
	if err != nil {
		return nil, err
	}

	result := make([]string, 0)

	w := session.Walk(remote)
	if w == nil {
		return result, RemoteWalkErr
	}

	remoteDir := filepath.Dir(remote)
	if runtime.GOOS == "windows" {
		remoteDir = strings.Replace(remoteDir, "\\", "/", -1)
	}

	for w.Step() {
		stat := w.Stat()
		remoteName := w.Path()

		baseName := CleanPath4Download(strings.Replace(remoteName, remoteDir, "", 1))
		if strings.HasPrefix(baseName, "/") {
			baseName = strings.Replace(baseName, "/", "", 1)
		}
		if runtime.GOOS == "windows" {
			baseName = strings.Replace(baseName, "/", "\\", -1)
		}
		target := local + baseName

		if stat.Mode().IsRegular() {
			if e := s.downloadFile(session, target, remoteName); e != nil {
				result = append(result, target+" <- "+remoteName+" :ERROR: "+e.Error())
			} else {
				result = append(result, target+" <- "+remoteName+" :OK")
			}
		} else if stat.IsDir() {
			if e := os.Mkdir(target, os.ModeDir|0755); e != nil && !os.IsExist(e) {
				result = append(result, target+" <- "+remoteName+" :ERROR: "+e.Error())
			} else {
				result = append(result, target+" <- "+remoteName+" :OK")
			}
		} else {
			result = append(result, remoteName+" :ERROR: type not support")
		}
	}

	return result, nil

}

func (s sftpCopy) Upload(local, remote string) ([]string, error) {
	session, err := s.newSession()
	if err != nil {
		return nil, err
	}

	localDir := filepath.Dir(local)

	result := make([]string, 0)
	err = filepath.Walk(local, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			result = append(result, path+" :ERROR: "+err.Error())
			return nil
		}

		if info == nil {
			result = append(result, path+" :ERROR: not found")
			return nil
		}

		baseName := CleanPath4Upload(strings.Replace(path, localDir, "", 1))
		if strings.HasPrefix(baseName, string(os.PathSeparator)) {
			baseName = strings.Replace(baseName, string(os.PathSeparator), "", 1)
		}
		if runtime.GOOS == "windows" {
			baseName = strings.Replace(baseName, "\\", "/", -1)
		}
		target := remote + baseName

		if info.IsDir() {
			if e := session.MkdirAll(target); e != nil && !os.IsExist(err) {
				result = append(result, path+" -> "+target+" :ERROR: "+err.Error())
			} else {
				result = append(result, path+" -> "+target+" :OK")
			}
			return nil
		}

		if info.Mode().IsRegular() {
			if e := uploadFile(session, path, target); e != nil {
				result = append(result, path+" -> "+target+" :ERROR: "+e.Error())
			} else {
				result = append(result, path+" -> "+target+" :OK")
			}
			return nil
		}

		result = append(result, path+" :ERROR: type not support")
		return nil
	})

	if err != nil {
		return result, LocalWalkErr
	}

	return result, nil
}

func (s sftpCopy) downloadFile(session *sftp.Client, local, remote string) error {
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

func uploadFile(client *sftp.Client, local string, remote string) error {
	rf, err := client.Create(remote)
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
