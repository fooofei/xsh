package base

import (
	"flag"
	"fmt"
	"os"
)

var (
	Mode = flag.String("mode", "interaction", "The run mode. value ranges: task/cmd/copy/interaction")

	Task  = flag.String("task", "", "The task yaml when run task")
	Value = flag.String("value", "", "The value yaml when run task")

	Group = flag.String("group", "", "The group name when run cmd or copy")

	Cmd = flag.String("cmd", "", "The command line when run cmd")
	Su  = flag.Bool("su", false, "Su or not when run cmd")

	Direction = flag.String("direction", "", "The direction upload/download when run copy")
	Local     = flag.String("local", "", "The local path when run copy")
	Remote    = flag.String("remote", "", "The remote path when run copy")
)

func initFlag() {
	os.Mkdir(ConfigRootPath, os.ModeDir|0755)

	cmd := os.Args[0]
	flag.Usage = func() {
		fmt.Println(`SYNOPSIS:`, cmd, `[OPTION]... [ARG]...

VERSION:
`+Version+`

DESCRIPTION:`)
		flag.PrintDefaults()
	}
	flag.Parse()
}
