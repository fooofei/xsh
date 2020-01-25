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

func do(action SshAction) {
	Print(action.Do())
}

func sudo(action SshAction) {
	action.Su = true
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
