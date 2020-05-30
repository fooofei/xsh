package sh

import (
	"context"
	"fmt"
	. "github.com/xied5531/xsh/base"
	"strings"
	"time"
)

type Step struct {
	Type string `yaml:"type"` // command/copy

	// for command
	Commands []string `yaml:"commands"`
	Su       bool     `yaml:"su"`

	// for copy
	Direction string `yaml:"direction"` // upload/download
	Local     string `yaml:"local"`
	Remote    string `yaml:"remote"`
}

type SshAction struct {
	Name   string     `yaml:"name"`
	Single bool       `yaml:"single,omitempty"`
	Group  string     `yaml:"group,omitempty"`
	Detail HostDetail `yaml:"detail,omitempty"`
	Steps  []Step     `yaml:"steps"`
}

type SshActionResult struct {
	Name       string
	Target     string
	StepResult map[string][]SshResponse
	Err        error
}

func (s *SshAction) checkAction() error {
	if len(s.Steps) == 0 {
		return ActionEmptyErr
	}

	su := false
	for _, action := range s.Steps {
		if action.Type == "command" {
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
				return CopyDirectionErr
			}
			if err := checkFullPath(action.Local, action.Remote); err != nil {
				return err
			}
		}
	}

	if su {
		if !s.Single {
			a := XAuthMap[XHostMap[s.Group].Auth]
			if a.SuType == "" {
				return CommandSuErr
			}
		} else {
			if s.Detail.SuType == "" {
				return CommandSuErr
			}
		}
	}

	var err error
	if s.Detail.Password, err = GetPlainPassword(s.Detail.Password); err != nil {
		Error.Printf("action[%s] detail[%s] password decrypt error: %v", s.Name, s.Detail.Address, err)
	}
	if s.Detail.Passphrase, err = GetPlainPassword(s.Detail.Passphrase); err != nil {
		Error.Printf("action[%s] detail[%s] passphrase decrypt error: %v", s.Name, s.Detail.Address, err)
	}
	if s.Detail.SuPass, err = GetPlainPassword(s.Detail.SuPass); err != nil {
		Error.Printf("action[%s] detail[%s] suPass decrypt error: %v", s.Name, s.Detail.Address, err)
	}

	return nil
}

func (s *SshAction) Do() SshActionResult {
	sshActionResult := SshActionResult{
		Name:       s.Name,
		StepResult: make(map[string][]SshResponse),
	}

	if err := s.checkAction(); err != nil {
		sshActionResult.Err = err
		return sshActionResult
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Duration(XConfig.Timeout.ActionTimeoutS)*time.Second)

	if !s.Single {
		sshActionResult.Target = s.Group

		responseCh := make(chan map[string][]SshResponse, 1)
		defer close(responseCh)
		s.do4group(ctx, responseCh)

		sshActionResult.StepResult = <-responseCh
	} else {
		sshActionResult.Target = s.Detail.Address

		responseCh := make(chan []SshResponse, 1)
		defer close(responseCh)
		s.do4host(ctx, s.Detail, responseCh)

		sshActionResult.StepResult[s.Detail.Address] = <-responseCh
	}

	return sshActionResult
}

func (s *SshAction) newSshCopy(hostDetail HostDetail) sshCopy {
	resut := sshCopy{
		SshClient: sshClient{
			Host:       hostDetail.Address,
			Port:       hostDetail.Port,
			UserName:   hostDetail.Username,
			PassWord:   hostDetail.Password,
			PrivateKey: hostDetail.PrivateKey,
			PassPhrase: hostDetail.Passphrase,
		},
	}

	if hostDetail.Port <= 0 {
		resut.SshClient.Port = XConfig.Command.Port
	}

	resut.Copy = scpCopy{
		SshClient: resut.SshClient,
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

func (s *SshAction) do4host(ctx context.Context, hostDetail HostDetail, resultCh chan []SshResponse) {
	responseCh := make(chan SshResponse, 1)
	defer close(responseCh)

	result := make([]SshResponse, 0)

	for _, action := range s.Steps {
		switch action.Type {
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
			sc.Local = strings.TrimSpace(action.Local)
			sc.Remote = strings.TrimSpace(action.Remote)

			s.doCopy(ctx, responseCh, sc)
			result = append(result, <-responseCh)
		}
	}

	resultCh <- result
}

func (s *SshAction) do4group(ctx context.Context, resultCh chan map[string][]SshResponse) {
	responseCh := make(chan []SshResponse, XConfig.Concurrency)
	defer close(responseCh)

	xHost, _ := XHostMap[s.Group]
	go func() {
		for _, hostDetail := range xHost.AllHost {
			go s.do4host(ctx, hostDetail, responseCh)
		}
	}()

	size := len(xHost.AllHost)
	result := make(map[string][]SshResponse)
	for i := 0; i < size; i++ {
		response := <-responseCh
		result[response[0].Address] = response
		printProgress(response, false)
	}
	printProgress(nil, true)

	resultCh <- result
}

func (s *SshAction) doCommand(ctx context.Context, resultCh chan SshResponse, sc sshCommand) {
	rc := make(chan SshResponse, 1)
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
		resultCh <- SshResponse{Address: sc.Session.Client.Host, Err: ActionTimeoutErr}
	}
}

func (s *SshAction) doCopy(ctx context.Context, resultCh chan SshResponse, sc sshCopy) {
	rc := make(chan SshResponse, 1)
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
		resultCh <- SshResponse{Address: sc.SshClient.Host, Err: ActionTimeoutErr}
	}
}

func printProgress(response []SshResponse, end bool) {
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
