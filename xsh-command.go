package main

import (
	"fmt"
	. "github.com/xied5531/xsh/base"
	. "github.com/xied5531/xsh/out"
	. "github.com/xied5531/xsh/sh"
	"strings"
)

func show(line string) {
	fields := strings.Fields(line)
	if len(fields) > 1 {
		switch fields[1] {
		case "address":
			group, _ := XHostMap[CurEnv.HostGroup]
			for i, v := range group.AllHost {
				fmt.Printf("[%d] %s\n", i, v.Address)
			}
		case "env":
			fmt.Printf("%+v\n", CurEnv)
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

func do(action SshAction, line string, su bool) {
	cmds := strings.Split(line, XConfig.CommandSep)
	if len(cmds) == 0 {
		return
	}

	action.SubActions[0].ActionType = "command"
	action.SubActions[0].Commands = cmds
	action.SubActions[0].Su = su

	Print(action.Do())
}

func copy(action SshAction, line string) {
	var direction string
	var fields []string

	if strings.Contains(line, "->") {
		direction = "upload"
		fields = strings.Split(line, "->")
	} else if strings.Contains(line, "<-") {
		direction = "download"
		fields = strings.Split(line, "<-")
	} else {
		Error.Printf("line[%s] format illegal", line)
		return
	}

	local, le := GetLocalPath(direction, fields[0])
	remote, re := GetRemotePath(direction, fields[1])
	if le != nil || re != nil {
		Error.Printf("line[%s] path illegal, local: %v; remote: %v", line, le, re)
		return
	}

	action.SubActions[0].Direction = direction
	action.SubActions[0].Local = local
	action.SubActions[0].Remote = remote
	action.SubActions[0].ActionType = "copy"
	Print(action.Do())
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
