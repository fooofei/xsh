package main

import (
	"fmt"
	. "github.com/xied5531/xsh/base"
)

func help() {
	fmt.Println(`NAME:
   xsh - A simple command line tool to run commands on remote host addresses.

SYNOPSIS:
   xsh [COMMANDS] [arguments...] 

VERSION:
` + Version + `

DESCRIPTION:
   Please report a issue at ` + XConfig.IssueUrl + ` if possible.

   :help or :h                         Show help info
   :set [group=xxx|address=x.x.x.x]    Load environment
   :show                               Show address list of current group

   :do                                 Run ssh command as normal user, default.
   :sudo                               Run ssh command as su user from normal user
   :copy                               Upload or download file or directory between local and remote
       local -> remote                 -> means upload, remote must be directory
       local <- remote                 <- means download, local must be directory

   :encrypt passwd                     Encrypt passwd
   :decrypt passwd                     Decrypt passwd

   :exit or :quit or :q                Stop xsh
`)
}
