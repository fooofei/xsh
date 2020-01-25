package main

import "fmt"

const version = "1.0.0"

func help() {
	fmt.Println(`NAME:
   xsh - A simple command line tool to run commands on remote host addresses.

USAGE:
   xsh [OPTIONS|COMMANDS] [arguments...] 

VERSION:
` + version + `

OPTIONS:
   :help                                  show help info
   :set [group=xxx|address=x.x.x.x]       load environment
   :show                                  show address list of current group
   :exit                                  quit

COMMANDS:
   :do                    run ssh command as normal user, default.
   :sudo                  run ssh command as su user from normal user
`)
}
