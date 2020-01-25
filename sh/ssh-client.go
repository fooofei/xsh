package sh

import (
	"fmt"
	. "github.com/xied5531/xsh/base"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
	"time"
)

type sshClient struct {
	Host string
	Port int

	UserName string
	PassWord string

	PrivateKey string
	PassPhrase string
}

func (s *sshClient) getCacheKey() string {
	return s.UserName + "+" + s.PrivateKey + "@" + s.Host + ":" + fmt.Sprintf("%d", s.Port)
}

func (s *sshClient) NewClient() (*ssh.Client, error) {
	key := s.getCacheKey()
	client, err := _sshClientCache.get(key)

	if err != nil {
		client, err = s.dial()
		if err != nil {
			return nil, err
		}
		_sshClientCache.add(key, client)
		return client, nil
	}

	return client, nil
}

func (s *sshClient) dial() (*ssh.Client, error) {
	var err error
	authMethods := make([]ssh.AuthMethod, 0)
	if s.PassWord != "" {
		authMethods = append(authMethods, ssh.Password(s.PassWord))

		keyboardInteractiveChallenge := func(
			user,
			instruction string,
			questions []string,
			echos []bool,
		) (answers []string, err error) {
			if len(questions) == 0 {
				return []string{}, nil
			}
			return []string{s.PassWord}, nil
		}
		authMethods = append(authMethods, ssh.KeyboardInteractive(keyboardInteractiveChallenge))
	}
	if s.PrivateKey != "" {
		var (
			pemBytes []byte
			signer   ssh.Signer
		)
		pemBytes, err = ioutil.ReadFile(s.PrivateKey)
		if err != nil {
			return nil, err
		}
		if s.PassPhrase != "" {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(pemBytes, []byte(s.PassPhrase))
		} else {
			signer, err = ssh.ParsePrivateKey(pemBytes)
		}
		if err != nil {
			return nil, err
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	clientConfig := &ssh.ClientConfig{
		User:    s.UserName,
		Auth:    authMethods,
		Timeout: time.Duration(XConfig.Timeout.DialTimeoutS) * time.Second,
		Config: ssh.Config{
			Ciphers: XConfig.Command.Ciphers,
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	return ssh.Dial("tcp", fmt.Sprintf("%s:%d", s.Host, s.Port), clientConfig)
}
