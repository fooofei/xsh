package base

import (
	"flag"
	"fmt"
	"os"
)

var (
	RunMode = flag.String("m", "interaction", "The run mode. value ranges: task/command/interaction(default)")

	TaskFile  = flag.String("t", "", "The task yaml when run task")
	ValueFile = flag.String("v", "", "The value yaml when run task")

	GroupName = flag.String("g", "", "The group name when run command")
	SshPort   = flag.Int("p", 22, "The ssh port when run command, default: 22")
	Cmdline   = flag.String("c", "", "The command line when run command")
)

func initFlag() {
	os.Mkdir(ConfigRootPath, os.ModeDir|0755)

	cmd := os.Args[0]
	flag.Usage = func() {
		fmt.Println(`Usage:`, cmd, `[<options>]

Options:`)
		flag.PrintDefaults()
	}
	flag.Parse()
}
