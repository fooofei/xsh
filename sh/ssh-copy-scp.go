package sh

import (
	"github.com/hnakamur/go-scp"
	"os"
	"path/filepath"
	"strings"
)

type scpCopy struct {
	SshClient sshClient
}

func (s scpCopy) newSession() (*scp.SCP, error) {
	client, err := s.SshClient.NewClient()
	if err != nil {
		return nil, err
	}

	return scp.NewSCP(client), nil
}

func (s scpCopy) Download(local, remote string) ([]string, error) {
	local = local + s.SshClient.Host + string(os.PathSeparator)
	if err := os.Mkdir(local, os.ModeDir|0755); err != nil && !os.IsExist(err) {
		return nil, err
	}

	if strings.HasSuffix(remote, "/") {
		return s.downloadDir(local, remote)
	} else {
		return s.downloadFile(local, remote)
	}
}

func (s scpCopy) downloadFile(local, remote string) ([]string, error) {
	local = local + filepath.Base(remote)

	session, err := s.newSession()
	if err != nil {
		return nil, err
	}

	if err := session.ReceiveFile(remote, local); err != nil {
		return nil, err
	} else {
		return []string{local + " <- " + remote + " :FILE:OK"}, nil
	}
}

func (s scpCopy) downloadDir(local, remote string) ([]string, error) {
	local = local + filepath.Base(filepath.Dir(remote))

	session, err := s.newSession()
	if err != nil {
		return nil, err
	}

	if err := session.ReceiveDir(remote, local, nil); err != nil {
		return nil, err
	} else {
		return []string{local + " <- " + remote + " :DIR:OK"}, nil
	}
}

func (s scpCopy) Upload(local, remote string) ([]string, error) {
	if strings.HasSuffix(local, string(os.PathSeparator)) {
		return s.uploadDir(local, remote)
	} else {
		return s.uploadFile(local, remote)
	}
}

func (s scpCopy) uploadFile(local, remote string) ([]string, error) {
	remote = remote + filepath.Base(local)

	session, err := s.newSession()
	if err != nil {
		return nil, err
	}

	if err := session.SendFile(local, remote); err != nil {
		return nil, err
	} else {
		return []string{local + " -> " + remote + " :FILE:OK"}, nil
	}
}

func (s scpCopy) uploadDir(local, remote string) ([]string, error) {
	remote = remote + filepath.Base(filepath.Dir(local))

	session, err := s.newSession()
	if err != nil {
		return nil, err
	}

	if err := session.SendDir(local, remote, nil); err != nil {
		return nil, err
	} else {
		return []string{local + " -> " + remote + " :DIR:OK"}, nil
	}
}
