package sh

import (
	"context"
	. "github.com/xied5531/xsh/base"
	"time"
)

// TargetType: group/Address
type SshAction struct {
	Name       string
	TargetType string
	HostGroup  string
	Commands   []string
	HostDetail HostDetail
	Su         bool
}

type SshActionResult struct {
	Name   string
	Target string
	Result []sshResponse
	Err    error
}

func (s *SshAction) Do() SshActionResult {
	result := SshActionResult{
		Name: s.Name,
	}

	if err := checkCommands(s.Commands); err != nil {
		result.Err = err
		return result
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Duration(XConfig.Timeout.ActionTimeoutS)*time.Second)

	if s.TargetType == "group" {
		result.Target = s.HostGroup
		result.Result = s.do4group(ctx)
	} else {
		result.Target = s.HostDetail.Address
		result.Result = s.do4host(ctx)
	}

	return result
}

func (s *SshAction) newSshCommand(hostDetail HostDetail) sshCommand {
	resut := sshCommand{
		Session: sshSession{
			Client: sshClient{
				Host:       hostDetail.Address,
				Port:       hostDetail.Port,
				UserName:   hostDetail.Username,
				PassWord:   hostDetail.Password,
				PrivateKey: hostDetail.PrivateKey,
				PassPhrase: hostDetail.Passphrase,
			},
		},
		Commands: s.Commands,
		SuType:   "",
		SuPass:   "",
	}

	if hostDetail.Port <= 0 {
		resut.Session.Client.Port = XConfig.Ssh.Port
	}

	if s.Su {
		if hostDetail.SuType != "" {
			resut.SuType = hostDetail.SuType
			resut.SuPass = hostDetail.SuPass
		} else {
			resut.SuType = "XXX"
			resut.SuPass = "XXX"
		}
	}

	return resut
}

func (s *SshAction) do4host(ctx context.Context) []sshResponse {
	resultCh := make(chan sshResponse, 1)
	defer close(resultCh)

	go s.doCommands(ctx, resultCh, s.newSshCommand(s.HostDetail))

	results := make([]sshResponse, 1)
	results = append(results, <-resultCh)

	return results
}

func (s *SshAction) do4group(ctx context.Context) []sshResponse {
	resultCh := make(chan sshResponse, XConfig.Concurrency)
	defer close(resultCh)

	xHost, _ := XHostMap[s.HostGroup]
	go func() {
		for _, hostDetail := range xHost.AllHost {
			go s.doCommands(ctx, resultCh, s.newSshCommand(hostDetail))
		}
	}()

	size := len(xHost.AllHost)
	results := make([]sshResponse, 0)
	for i := 0; i < size; i++ {
		results = append(results, <-resultCh)
	}

	return results
}

func (s *SshAction) doCommands(ctx context.Context, resultCh chan sshResponse, sc sshCommand) {
	rc := make(chan sshResponse, 1)
	go func() {
		defer close(rc)
		select {
		case <-ctx.Done():
			return
		default:
			rc <- sc.run()
		}
	}()

	select {
	case r := <-rc:
		resultCh <- r
	case <-ctx.Done():
		resultCh <- sshResponse{Err: ActionTimeoutErr}
	}
}
