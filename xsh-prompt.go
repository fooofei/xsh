package main

import (
	"fmt"
	. "github.com/xied5531/xsh/base"
	. "github.com/xied5531/xsh/out"
	. "github.com/xied5531/xsh/sh"
	"log"
	"sort"
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

	switch CurEnv.ActionType {
	case ":do":
		if action := newCommandAction(line, false); action != nil {
			Out(action.Do())
		}
	case ":sudo":
		if action := newCommandAction(line, true); action != nil {
			Out(action.Do())
		}
	case ":copy":
		if action, err := newCopyAction(line); err != nil {
			Out(action.Do())
		} else {
			Error.Println(err.Error())
		}
	}

	return false
}

func newCommandAction(line string, su bool) *SshAction {
	cmds := strings.Split(line, XConfig.CommandSep)
	if len(cmds) == 0 {
		return nil
	}

	action := &SshAction{
		Name:       "Default",
		TargetType: CurEnv.TargetType,
		SubActions: []SubAction{{
			ActionType: "command",
			Commands:   cmds,
			Su:         su,
		}},
	}
	setupActionTaget(action)

	return action
}

func newCopyAction(line string) (*SshAction, error) {
	var direction string
	var fields []string

	if strings.Contains(line, "->") {
		direction = "upload"
		fields = strings.Split(line, "->")
	} else if strings.Contains(line, "<-") {
		direction = "download"
		fields = strings.Split(line, "<-")
	} else {
		return nil, fmt.Errorf("line[%s] format illegal", line)
	}

	local, le := GetLocalPath(direction, fields[0])
	remote, re := GetRemotePath(direction, fields[1])
	if le != nil || re != nil {
		return nil, fmt.Errorf("line[%s] path illegal, local: %v; remote: %v", line, le, re)
	}

	action := &SshAction{
		Name:       "Default",
		TargetType: CurEnv.TargetType,
		SubActions: []SubAction{{
			ActionType: "copy",
			Direction:  direction,
			Local:      local,
			Remote:     remote,
		}},
	}
	setupActionTaget(action)

	return action, nil
}

func setupActionTaget(action *SshAction) {
	if action.TargetType == "group" {
		action.HostGroup = CurEnv.HostGroup
	} else {
		authentication := XAuthMap[CurEnv.Authentication]
		action.HostDetail = HostDetail{
			Address:    CurEnv.HostAddress,
			Username:   authentication.Username,
			Password:   authentication.Password,
			PrivateKey: authentication.PrivateKey,
			Passphrase: authentication.Passphrase,
			SuType:     authentication.SuType,
			SuPass:     authentication.SuPass,
		}
	}
}

func show(line string) {
	fields := strings.Fields(line)
	if len(fields) > 1 {
		switch fields[1] {
		case "address":
			group, _ := XHostMap[CurEnv.HostGroup]
			addresses := make([]string, len(group.AllHost))
			for i, v := range group.AllHost {
				addresses[i] = v.Address
			}
			sort.Strings(addresses)
			PrintYaml(addresses)
		case "env":
			PrintYaml(CurEnv)
		}
	} else {
		fmt.Println(":show argument not enough, please :help first.")
	}
}

func set(line string) {
	fields := strings.Fields(line)
	if len(fields) > 1 {
		for _, v := range fields[1:] {
			currFields := strings.Split(v, "=")
			if len(currFields) != 2 {
				fmt.Println(":set argument illegal, please :help first.")
				return
			}
			switch currFields[0] {
			case "group":
				if group, ok := XHostMap[currFields[1]]; !ok {
					fmt.Printf("group[%s] not found\n", currFields[1])
					return
				} else {
					CurEnv.TargetType = currFields[0]
					CurEnv.HostGroup = group.Name
					CurEnv.Authentication = group.Authentication
				}
			case "address":
				if !CheckIp(currFields[1]) {
					fmt.Printf("address[%s] illegal\n", currFields[1])
					return
				}
				if CurEnv.Authentication == "" {
					fmt.Println("authentication empty, please :set group= first.")
					return
				}
				if !ContainsAddress(currFields[1], XHostMap[CurEnv.HostGroup].AllHost) {
					fmt.Printf("address[%s] not found in group [%s]\n", currFields[1], CurEnv.HostGroup)
					return
				}
				CurEnv.TargetType = currFields[0]
				CurEnv.HostAddress = currFields[1]
			}
		}
		SaveEnv()
	} else {
		fmt.Println(":set argument not enough, please :help first.")
	}
}

func encrypt(line string) {
	fields := strings.Fields(line)
	if len(fields) == 2 {
		fmt.Println(GetMaskPassword(fields[1]))
	} else {
		Error.Printf("line[%s] illegal", line)
	}
}

func decrypt(line string) {
	fields := strings.Fields(line)
	if len(fields) == 2 {
		fmt.Println(GetPlainPassword(fields[1]))
	} else {
		Error.Printf("line[%s] illegal", line)
	}
}

func reload() {
	InitXConfig()
	InitXAuth()
	InitXHost()
}
