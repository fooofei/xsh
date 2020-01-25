package base

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	Mode = flag.String("mode", "prompt", "The running mode. value ranges: task/cmd/copy/crypt/prompt")

	Task  = flag.String("task", "", "The task yaml file in task mode")
	Value = flag.String("value", "", "The value yaml file in task mode")

	Group = flag.String("group", "", "The group name in cmd or copy mode")

	Cmd = flag.String("cmd", "", "The command line in cmd mode")
	Su  = flag.Bool("su", false, "Su or not in cmd mode")

	Direction = flag.String("direction", "", "The direction upload/download in copy mode")
	Local     = flag.String("local", "", "The local path in copy mode")
	Remote    = flag.String("remote", "", "The remote path in copy mode")

	Plain  = flag.String("plain", "", "The plain text to encrypt in crypt mode")
	Cipher = flag.String("cipher", "", "The cipher text to decrypt in crypt mode")
)

func initFlag() {
	if err := os.Mkdir(RootPath, os.ModeDir|0755); err != nil && !os.IsExist(err) {
		log.Fatalf("mkdir %s error: %v\n", RootPath, err)
	}
	if err := os.Mkdir(LogPath, os.ModeDir|0755); err != nil && !os.IsExist(err) {
		log.Fatalf("mkdir %s error: %v\n", LogPath, err)
	}
	if err := os.Mkdir(TempPath, os.ModeDir|0755); err != nil && !os.IsExist(err) {
		log.Fatalf("mkdir %s error: %v\n", TempPath, err)
	}

	cmd := os.Args[0]
	flag.Usage = func() {
		fmt.Println(`NAME:
   xsh - A command line tool that can execute commands on remote hosts or upload and download files.

SYNOPSIS:`, cmd, `[OPTION]... [ARG]...

VERSION:
`+Version+`

DESCRIPTION:`)
		flag.PrintDefaults()
	}
	flag.Parse()
}
