package base

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	Mode = flag.String("mode", "prompt", "running mode. value ranges: task/cmd/copy/crypt/prompt")

	Output = flag.String("o", "text", "output mode in task cmd and copy mode. value ranges: text/json/yaml")

	Task  = flag.String("task", "", "task yaml file in task mode")
	Value = flag.String("value", "", "value yaml file in task mode")

	Group = flag.String("group", "", "group name in cmd and copy mode")

	Cmd = flag.String("cmd", "", "command line in cmd mode")
	Su  = flag.Bool("su", false, "su or not in cmd mode")

	Direction = flag.String("direction", "", "upload or download in copy mode")
	Local     = flag.String("local", "", "local path in copy mode")
	Remote    = flag.String("remote", "", "remote path in copy mode")

	Plain     = flag.String("plain", "", "plain text to encrypt in crypt mode")
	Cipher    = flag.String("cipher", "", "cipher text to decrypt in crypt mode")
	CryptType = flag.String("ctype", "aes", "crypt or decrypt type in crypt mode")
	CryptKey  = flag.String("ckey", "", "crypt or decrypt key in crypt mode, length must be 32")
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

	checkFlag()
}

func checkFlag() {
	if *Mode != "task" && *Mode != "cmd" && *Mode != "copy" && *Mode != "crypt" && *Mode != "prompt" {
		log.Fatalf("mode[%s] illegal", *Mode)
	}

	if *Output != "text" && *Output != "json" && *Output != "yaml" {
		log.Fatalf("o[%s] illegal", *Output)
	}

	if *Direction != "" && *Direction != "upload" && *Direction != "download" {
		log.Fatalf("direction[%s] illegal", *Direction)
	}

	if *Mode == "crypt" {
		if (*Plain == "" && *Cipher == "") || *CryptKey == "" || len(*CryptKey) != 32 {
			log.Fatal("plain or cipher or ckey illegal")
		}
	}
}
