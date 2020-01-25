package main

import (
	"fmt"
	. "github.com/xied5531/xsh/base"
	. "github.com/xied5531/xsh/sh"
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
	if line == ":exit" || line == ":quit" || line == ":q" {
		return true
	}
	if line == ":help" || line == ":h" {
		help()
		return false
	}
	if line == ":reload" {
		reload()
		return false
	}
	if strings.HasPrefix(line, ":encrypt") {
		encrypt(line)
		return false
	}
	if strings.HasPrefix(line, ":decrypt") {
		decrypt(line)
		return false
	}
	if strings.HasPrefix(line, ":set") {
		set(line)
		return false
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
		}
		return false
	case ":sudo":
		if CurEnv.ActionType != ":sudo" {
			CurEnv.ActionType = ":sudo"
			SaveEnv()
		}
		return false
	case ":copy":
		if CurEnv.ActionType != ":copy" {
			CurEnv.ActionType = ":copy"
			SaveEnv()
		}
		return false
	}

	if strings.HasPrefix(line, ":") {
		Error.Printf("command[%s] not support", line)
		return false
	}

	sshAction := newSshAction()
	switch CurEnv.ActionType {
	case ":do":
		do(sshAction, line, false)
	case ":sudo":
		do(sshAction, line, true)
	case ":copy":
		copy(sshAction, line)
	}

	return false
}

func newSshAction() SshAction {
	result := SshAction{
		Name:       "Default",
		TargetType: CurEnv.TargetType,
		SubActions: make([]SubAction, 1),
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
