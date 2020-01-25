package main

import (
	"fmt"
	. "github.com/xied5531/xsh/base"
	. "github.com/xied5531/xsh/sh"
	"strings"
)

func runTask() {
	task := SshTask{}

	CurEnv.Output = *Output
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

	CurEnv.Output = *Output
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

	local, le := GetLocalPath(*Direction, *Local)
	remote, re := GetRemotePath(*Direction, *Remote)
	if le != nil || re != nil {
		Error.Printf("path illegal, local: %v; remote: %v", le, re)
		return
	}

	action := SshAction{
		Name:  "Default",
		Group: *Group,
		Steps: []Step{{
			Type:      "copy",
			Direction: *Direction,
			Local:     local,
			Remote:    remote,
		}},
	}

	CurEnv.Output = *Output
	Out(action.Do())
}

func runCrypt() {
	if XConfig.Crypt.Type == "" || XConfig.Crypt.Key == "" {
		Error.Println("crypt type or key not found")
		return
	}

	if *Plain != "" {
		if c, e := GetMaskPassword(*Plain); e != nil {
			fmt.Printf("%s -> error: %s\n", *Plain, e.Error())
		} else {
			fmt.Printf("%s -> %s\n", *Plain, c)
		}
	} else if *Cipher != "" {
		if p, e := GetPlainPassword(*Cipher); e != nil {
			fmt.Printf("%s -> error: %s\n", *Cipher, e.Error())
		} else {
			fmt.Printf("%s -> %s\n", *Cipher, p)
		}
	} else {
		Error.Println("crypt plain or cipher text not found")
	}
}
