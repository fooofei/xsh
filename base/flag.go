package base

import (
	"flag"
	"fmt"
	"os"
)

var (
	taskFile  = flag.String("t", "", "The tasks yaml.")
	valueFile = flag.String("v", "", "The values yaml.")
	authName  = flag.String("A", "", "The auth name.")
	hostName  = flag.String("H", "", "The host name.")
	sshport   = flag.Int("P", 22, "The ssh port.")
	cmdline   = flag.String("L", "", "The command line.")
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
