package main

import (
	"fmt"
	. "github.com/xied5531/xsh/base"
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
	if strings.HasPrefix(line, ":show") {
		show(line)
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

	switch line {
	case ":do":
		if CurEnv.Action != ":do" {
			CurEnv.Action = ":do"
			SaveEnv()
		}
		return false
	case ":sudo":
		if CurEnv.Action != ":sudo" {
			CurEnv.Action = ":sudo"
			SaveEnv()
		}
		return false
	case ":copy":
		if CurEnv.Action != ":copy" {
			CurEnv.Action = ":copy"
			SaveEnv()
		}
		return false
	}

	if strings.HasPrefix(line, ":") {
		Error.Printf("command[%s] not support", line)
		return false
	}

	switch CurEnv.Action {
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
			Error.Println(err.Error())
		} else {
			Out(action.Do())
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
		Name:   "Default",
		Single: CurEnv.Single,
		Steps: []Step{{
			Type:     "command",
			Commands: cmds,
			Su:       su,
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

	local, le := GetLocalDir(direction, fields[0])
	remote, re := GetRemoteDir(direction, fields[1])
	if le != nil || re != nil {
		return nil, fmt.Errorf("line[%s] path illegal, local: %v; remote: %v", line, le, re)
	}

	action := &SshAction{
		Name:   "Default",
		Single: CurEnv.Single,
		Steps: []Step{{
			Type:      "copy",
			Direction: direction,
			Local:     local,
			Remote:    remote,
		}},
	}
	setupActionTaget(action)

	return action, nil
}

func setupActionTaget(action *SshAction) {
	if !action.Single {
		action.Group = CurEnv.Group
	} else {
		auth := XAuthMap[CurEnv.Auth]
		action.Detail = HostDetail{
			Address:    CurEnv.Address,
			Username:   auth.Username,
			Password:   auth.Password,
			PrivateKey: auth.PrivateKey,
			Passphrase: auth.Passphrase,
			SuType:     auth.SuType,
			SuPass:     auth.SuPass,
		}
	}
}

func show(line string) {
	fields := strings.Fields(line)
	if len(fields) > 1 {
		switch fields[1] {
		case "group":
			var groups []string
			for k, _ := range XHostMap {
				groups = append(groups, k)
			}
			sort.Strings(groups)
			PrintYaml(groups)
		case "address":
			group, ok := XHostMap[CurEnv.Group]
			if !ok {
				fmt.Printf("current group[%s] not found\n", CurEnv.Group)
				return
			}
			addresses := make([]string, len(group.AllHost))
			for i, v := range group.AllHost {
				addresses[i] = v.Address
			}
			sort.Strings(addresses)
			PrintYaml(addresses)
		case "env":
			PrintYaml(CurEnv)
		case "config":
			PrintYaml(XConfig)
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
					CurEnv.Single = false
					CurEnv.Group = group.Name
					CurEnv.Auth = group.Auth
				}
			case "address":
				if !CheckIp(currFields[1]) {
					fmt.Printf("address[%s] illegal\n", currFields[1])
					return
				}
				if CurEnv.Group == "" {
					fmt.Println("group empty, please :set group= first.")
					return
				}
				group := XHostMap[CurEnv.Group]
				CurEnv.Auth = group.Auth

				if !ContainsAddress(currFields[1], XHostMap[CurEnv.Group].AllHost) {
					fmt.Printf("address[%s] not found in group [%s]\n", currFields[1], CurEnv.Group)
					return
				}
				CurEnv.Single = true
				CurEnv.Address = currFields[1]
			}
		}
		SaveEnv()
	} else {
		fmt.Println(":set argument not enough, please :help first.")
	}
}

func reload() {
	InitXConfig()
	InitXAuth()
	InitXHost()
}

func help() {
	fmt.Println(`NAME:
   xsh - A command line tool that can execute commands on remote hosts or upload and download files.

SYNOPSIS:
   [KEYWORD] [ARG]... 

VERSION:
` + Version + `

DESCRIPTION:
   Please report a issue at ` + XConfig.IssueUrl + ` if possible.

   :help or :h                         Show help info
   :set [group=xxx|address=x.x.x.x]    Set current target hosts
   :show                               Show current information
   :reload                             Reload config auth and host information

   :do                                 Run ssh command as normal user
   :sudo                               Run ssh command as su user from normal user
   :copy                               Upload or download files between local and remote
     local -> remote                   -> means upload, remote must be directory
     local <- remote                   <- means download, local must be directory

   :exit or :quit or :q                Exit xsh
`)
}
