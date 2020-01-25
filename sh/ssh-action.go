package sh

import (
	"context"
	"fmt"
	. "github.com/xied5531/xsh/base"
	"time"
)

type SubAction struct {
	ActionType string // command/copy

	// for ssh
	Commands []string
	Su       bool

	// for sftp/scp
	Direction string // upload/download
	Local     string
	Remote    string
}

type SshAction struct {
	Name       string
	TargetType string // group/Address
	HostGroup  string
	HostDetail HostDetail
	SubActions []SubAction
}

type SshActionResult struct {
	Name   string
	Target string
	Result map[string][]sshResponse
	Err    error
}

func (s *SshAction) checkAction() error {
	if len(s.SubActions) == 0 {
		return ActionEmptyErr
	}

	su := false
	for _, action := range s.SubActions {
		if action.ActionType == "command" {
			if len(action.Commands) == 0 {
				return CommandEmptyErr
			}
			if err := checkCommands(action.Commands); err != nil {
				return err
			}
			if action.Su {
				su = true
			}
		} else {
			if action.Direction != "upload" && action.Direction != "download" {
				return CopyInfoErr
			}
			if action.Local == "" || action.Remote == "" {
				return CopyInfoErr
			}
		}
	}

	if su {
		if s.TargetType == "group" {
			a := XAuthMap[XHostMap[s.HostGroup].Authentication]
			if a.SuType == "" {
				return CommandSuErr
			}
		} else {
			if s.HostDetail.SuType == "" {
				return CommandSuErr
			}
		}
	}

	return nil
}

func (s *SshAction) Do() SshActionResult {
	sshActionResult := SshActionResult{
		Name: s.Name,
	}

	if err := s.checkAction(); err != nil {
		sshActionResult.Err = err
		return sshActionResult
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Duration(XConfig.Timeout.ActionTimeoutS)*time.Second)

	if s.TargetType == "group" {
		sshActionResult.Target = s.HostGroup

		responseCh := make(chan map[string][]sshResponse, 1)
		defer close(responseCh)
		s.do4group(ctx, responseCh)

		sshActionResult.Result = <-responseCh
	} else {
		sshActionResult.Target = s.HostDetail.Address

		responseCh := make(chan []sshResponse, 1)
		defer close(responseCh)
		s.do4host(ctx, s.HostDetail, responseCh)

		sshActionResult.Result[s.HostDetail.Address] = <-responseCh
	}

	return sshActionResult
}

func (s *SshAction) newSshCopy(hostDetail HostDetail) sshCopy {
	resut := sshCopy{
		Session: copySession{
			Client: sshClient{
				Host:       hostDetail.Address,
				Port:       hostDetail.Port,
				UserName:   hostDetail.Username,
				PassWord:   hostDetail.Password,
				PrivateKey: hostDetail.PrivateKey,
				PassPhrase: hostDetail.Passphrase,
			},
		},
	}

	if hostDetail.Port <= 0 {
		resut.Session.Client.Port = XConfig.Command.Port
	}

	return resut
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
		SuType: hostDetail.SuType,
		SuPass: hostDetail.SuPass,
	}

	if hostDetail.Port <= 0 {
		resut.Session.Client.Port = XConfig.Command.Port
	}

	return resut
}

func (s *SshAction) do4host(ctx context.Context, hostDetail HostDetail, resultCh chan []sshResponse) {
	responseCh := make(chan sshResponse, 1)
	defer close(responseCh)

	result := make([]sshResponse, 0)

	for _, action := range s.SubActions {
		switch action.ActionType {
		case "command":
			sc := s.newSshCommand(hostDetail)
			sc.Commands = action.Commands
			if !action.Su {
				sc.SuType = ""
				sc.SuPass = ""
			}

			s.doCommand(ctx, responseCh, sc)
			result = append(result, <-responseCh)
		case "copy":
			sc := s.newSshCopy(hostDetail)
			sc.Direction = action.Direction
			sc.Local = action.Local
			sc.Remote = action.Remote

			s.doCopy(ctx, responseCh, sc)
			result = append(result, <-responseCh)
		}
	}

	resultCh <- result
}

func (s *SshAction) do4group(ctx context.Context, resultCh chan map[string][]sshResponse) {
	responseCh := make(chan []sshResponse, XConfig.Concurrency)
	defer close(responseCh)

	xHost, _ := XHostMap[s.HostGroup]
	go func() {
		for _, hostDetail := range xHost.AllHost {
			go s.do4host(ctx, hostDetail, responseCh)
		}
	}()

	size := len(xHost.AllHost)
	result := make(map[string][]sshResponse)
	for i := 0; i < size; i++ {
		response := <-responseCh
		result[response[0].Address] = response
		printProgress(response, false)
	}
	printProgress(nil, true)

	resultCh <- result
}

func (s *SshAction) doCommand(ctx context.Context, resultCh chan sshResponse, sc sshCommand) {
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
		resultCh <- sshResponse{Address: sc.Session.Client.Host, Err: ActionTimeoutErr}
	}
}

func (s *SshAction) doCopy(ctx context.Context, resultCh chan sshResponse, sc sshCopy) {
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
		resultCh <- sshResponse{Address: sc.Session.Client.Host, Err: ActionTimeoutErr}
	}
}

func printProgress(response []sshResponse, end bool) {
	if XConfig.Output.Type == "text" && XConfig.Output.Progress {
		if end {
			fmt.Println()
			return
		}

		for _, res := range response {
			if res.Err != nil {
				fmt.Print("x")
				return
			}
		}

		fmt.Print(".")
	}
}
