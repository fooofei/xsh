package sh

import (
	"fmt"
	. "github.com/xied5531/xsh/base"
	"io"
	"strings"
	"time"
)

type sshCommand struct {
	Session  sshSession
	Commands []string
	SuType   string
	SuPass   string
}

func (s sshCommand) run() sshResponse {
	if s.SuType != "" {
		result := s.sudo()
		result.Address = s.Session.Client.Host
		return result
	}

	result := s.do()
	result.Address = s.Session.Client.Host
	return result
}

func (s sshCommand) do() sshResponse {
	cmd := "set -e;set +o history;" + strings.Join(s.Commands, ";")

	return withTimeout(s.Session.run, cmd, time.Duration(XConfig.Timeout.CommandTimeoutS)*time.Second)
}

func (s sshCommand) sudo() sshResponse {
	cmd := "set -e;set +o history;" + strings.Join(s.Commands, ";")

	if XConfig.Command.Pty.Term == "" {
		cmd := "echo '" + s.SuPass + "'|" + s.SuType + " -c '" + strings.Replace(cmd, "'", "\\'", -1) + "'"
		return withTimeout(s.Session.run, cmd, time.Duration(XConfig.Timeout.CommandTimeoutS)*time.Second)
	}

	return s.sudoWithPty()
}

func (s sshCommand) sudoWithPty() sshResponse {
	randName := fmt.Sprintf("%d", time.Now().UnixNano())

	res := withTimeout(s.Session.shell, func(stdin io.WriteCloser) sshResponse {
		if _, err := stdin.Write([]byte("set +o history;stty -echo;" + s.SuType + " || exit 110 \n")); err != nil {
			return sshResponse{Err: err}
		}
		time.Sleep(time.Millisecond * time.Duration(XConfig.Command.Pty.IntervalMS))
		if _, err := stdin.Write([]byte(s.SuPass + "\n")); err != nil {
			return sshResponse{Err: err}
		}
		time.Sleep(time.Millisecond * time.Duration(XConfig.Command.Pty.IntervalMS))

		commands := s.preProcessCmd4Pty(s.Commands, randName)

		for _, cmd := range commands {
			if _, err := stdin.Write([]byte(cmd + "\n")); err != nil {
				return sshResponse{Err: err}
			}
			time.Sleep(time.Millisecond * time.Duration(XConfig.Command.Pty.IntervalMS))
		}
		return sshResponse{}
	}, time.Duration(XConfig.Timeout.CommandTimeoutS)*time.Second)

	return s.postProcessOutput4Pty(res, randName)
}

func (s sshCommand) preProcessCmd4Pty(commands []string, randName string) []string {
	var ret []string

	shName := ".xsh." + randName + ".sh"
	stdoutName := shName + ".Stdout"
	stderrName := shName + ".Stderr"

	ret = append(ret, "stty -echo")
	ret = append(ret, "echo '#!/usr/bin/env sh'> "+shName)
	ret = append(ret, "echo -n '<XSH-STDOUT-BEGIN>' > "+stdoutName)
	ret = append(ret, "echo -n '<XSH-STDERR-BEGIN>' > "+stderrName)

	for _, cmd := range commands {
		if cmd != "" {
			cmd = strings.Replace(cmd, "'", "'\\''", -1)
			ret = append(ret, "echo '"+cmd+" || exit $?' >> "+shName)
		}
	}

	ret = append(ret, "chmod +x "+shName)
	ret = append(ret, "./"+shName+" 1>>"+stdoutName+" 2>>"+stderrName)
	ret = append(ret, "echo '<XSH-STDOUT-END>' >> "+stdoutName)
	ret = append(ret, "echo '<XSH-STDERR-END>' >> "+stderrName)
	ret = append(ret, "cat "+stdoutName)
	ret = append(ret, "cat "+stderrName)
	ret = append(ret, "rm -f "+".xsh."+randName+".*")
	ret = append(ret, "exit")
	ret = append(ret, "exit")

	return ret
}

func (s sshCommand) postProcessOutput4Pty(res sshResponse, randName string) sshResponse {
	stdoutBegin := "<XSH-STDOUT-BEGIN>"
	stdoutEnd := "<XSH-STDOUT-END>"
	stderrBegin := "<XSH-STDERR-BEGIN>"
	stderrEnd := "<XSH-STDERR-END>"

	var stdoutReturn = ""
	stdoutBeginIndex := strings.LastIndex(res.Stdout, stdoutBegin)
	stdoutEndIndex := strings.LastIndex(res.Stdout, stdoutEnd)
	if stdoutBeginIndex >= 0 && stdoutEndIndex > 0 && stdoutEndIndex >= stdoutBeginIndex+len([]byte(stdoutBegin)) {
		stdoutReturn = res.Stdout[stdoutBeginIndex+len([]byte(stdoutBegin)) : stdoutEndIndex]
	}

	var stderrReturn = ""
	stderrBeginIndex := strings.LastIndex(res.Stdout, stderrBegin)
	stderrEndIndex := strings.LastIndex(res.Stdout, stderrEnd)
	if stderrBeginIndex >= 0 && stderrEndIndex > 0 && stderrEndIndex >= stderrBeginIndex+len([]byte(stderrBegin)) {
		stderrReturn = res.Stdout[stderrBeginIndex+len([]byte(stderrBegin)) : stderrEndIndex]
		stderrReturn = strings.Replace(stderrReturn, ".xsh."+randName+".sh:", "", -1)
	}

	if res.Stderr != "" {
		stderrReturn = stderrReturn + res.Stderr
	}

	if stdoutReturn == "" && stderrReturn == "" {
		Warn.Printf("host: %s, Commands: %v, res: %v", s.Session.Client.Host, s.Commands, res)
	}

	return sshResponse{
		Stdout: stdoutReturn,
		Stderr: stderrReturn,
		Err:    nil,
	}
}
