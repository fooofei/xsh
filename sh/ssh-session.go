package sh

import (
	"bytes"
	"errors"
	. "github.com/xied5531/xsh/base"
	"golang.org/x/crypto/ssh"
	"io"
)

type sshSession struct {
	Client sshClient
}

func (s *sshSession) newSession() (*xSshSession, error) {
	client, err := s.Client.NewClient()
	if err != nil {
		return nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}

	if XConfig.Ssh.Pty.Term != "" {
		if err := session.RequestPty(XConfig.Ssh.Pty.Term, XConfig.Ssh.Pty.Height, XConfig.Ssh.Pty.Width, ssh.TerminalModes{ssh.TTY_OP_ISPEED: uint32(XConfig.Ssh.Pty.Ispeed), ssh.TTY_OP_OSPEED: uint32(XConfig.Ssh.Pty.Ospeed)}); err != nil {
			return nil, err
		}
	}

	return &xSshSession{
		Session: *session,
	}, nil
}

func (s *sshSession) run(arg interface{}) sshResponse {
	command := arg.(string)

	session, err := s.newSession()
	if err != nil {
		return sshResponse{Err: err}
	}
	defer session.Close()

	stdout, stderr, err := session.OutputStdoutStderr(command)
	if err != nil {
		return sshResponse{Err: err}
	}

	return sshResponse{
		Address: s.Client.Host,
		Stdout:  stdout,
		Stderr:  stderr,
		Err:     nil,
	}
}

type Callback func(stdin io.WriteCloser) error

func (s *sshSession) shell(arg interface{}) sshResponse {
	callback := arg.(Callback)

	session, err := s.newSession()
	if err != nil {
		return sshResponse{Err: err}
	}
	defer session.Close()

	stdin, err := session.StdinPipe()
	if err != nil {
		return sshResponse{Err: err}
	}

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer

	session.Stdout = &stdoutBuf
	session.Stderr = &stderrBuf

	err = session.Shell()
	if err != nil {
		return sshResponse{Err: err}
	}

	if err = callback(stdin); err != nil {
		return sshResponse{Err: err}
	}

	err = session.Wait()
	if err != nil {
		return sshResponse{Err: err}
	}

	return sshResponse{
		Address: s.Client.Host,
		Stdout:  stdoutBuf.String(),
		Stderr:  stderrBuf.String(),
		Err:     nil,
	}
}

type xSshSession struct {
	ssh.Session
}

func (x *xSshSession) OutputStdoutStderr(cmd string) (string, string, error) {
	if x.Stdout != nil {
		return "", "", errors.New("ssh: Stdout already set")
	}
	if x.Stderr != nil {
		return "", "", errors.New("ssh: Stderr already set")
	}
	var stdout, stderr bytes.Buffer
	x.Stdout = &stdout
	x.Stderr = &stderr
	err := x.Run(cmd)
	return stdout.String(), stderr.String(), err
}
