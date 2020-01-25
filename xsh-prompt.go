package main

import (
	"fmt"
	. "github.com/luckywinds/xsh/base"
	. "github.com/luckywinds/xsh/sh"
	"log"
	"strings"
)

func runPrompt() {
	l, err := NewLiner()
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		line, err := Prompt(l)
		if err == PromptAborted {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		}

		Debug.Printf("line: %s", line)
		if runCommand(line) {
			break
		}
	}
}

func runCommand(line string) bool {
	line = strings.Trim(line, " ")
	if line == "" {
		return false
	}
	if line == ":exit" {
		return true
	}
	if line == ":help" {
		help()
		return false
	}
	if strings.HasPrefix(line, ":set") {
		set(line)
	}
	if !CheckEnv() {
		fmt.Println("please :set group= first")
		return false
	}

	if strings.HasPrefix(line, ":show") {
		show(line)
		return false
	}

	switch line {
	case ":do":
		if CurEnv.ActionType != ":do" {
			CurEnv.ActionType = ":do"
			SaveEnv()
			return false
		}
	case ":sudo":
		if CurEnv.ActionType != ":sudo" {
			CurEnv.ActionType = ":sudo"
			SaveEnv()
			return false
		}
	}

	if strings.HasPrefix(line, ":") {
		Error.Printf("command not support, line= %s", line)
		return false
	}

	cmds := strings.Split(line, XConfig.CommandSep)
	if len(cmds) == 0 {
		return false
	}

	sshAction := newSshAction()
	sshAction.Commands = cmds

	switch CurEnv.ActionType {
	case ":do":
		do(sshAction)
	case ":sudo":
		sudo(sshAction)
	}

	return false
}

func newSshAction() SshAction {
	result := SshAction{
		Name:       "Default",
		TargetType: CurEnv.TargetType,
	}

	if result.TargetType == "group" {
		result.HostGroup = CurEnv.HostGroup

	} else {
		authentication := XAuthMap[CurEnv.Authentication]
		result.HostDetail = HostDetail{
			Address:    CurEnv.HostAddress,
			Username:   authentication.Username,
			Password:   authentication.Password,
			PrivateKey: authentication.PrivateKey,
			Passphrase: authentication.Passphrase,
			SuType:     authentication.SuType,
			SuPass:     authentication.SuPass,
		}
	}
	return result
}