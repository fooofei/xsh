package main

import (
	. "github.com/xied5531/xsh/base"
	. "github.com/xied5531/xsh/out"
	. "github.com/xied5531/xsh/sh"
	"strings"
)

func runTask() {
	task := SshTask{}
	Out(task.Do())
}

func runCmd() {
	if _, ok := XHostMap[*Group]; !ok {
		Error.Printf("host group[%s] not found\n", *Group)
		return
	}

	if len(strings.Trim(*Cmd, " ")) == 0 {
		Error.Printf("command line[%s] empty\n", *Cmd)
		return
	}

	action := SshAction{
		Name:  "Default",
		Group: *Group,
		Steps: []Step{{
			Type:     "command",
			Commands: strings.Split(*Cmd, XConfig.CommandSep),
			Su:       *Su,
		}},
	}

	Out(action.Do())
}

func runCopy() {
	if _, ok := XHostMap[*Group]; !ok {
		Error.Printf("host group[%s] not found\n", *Group)
		return
	}

	if *Direction == "" || *Local == "" || *Remote == "" {
		Error.Println("direction or local or remote not found")
		return
	}

	action := SshAction{
		Name:  "Default",
		Group: *Group,
		Steps: []Step{{
			Type:      "copy",
			Direction: *Direction,
			Local:     *Local,
			Remote:    *Remote,
		}},
	}

	Out(action.Do())
}
