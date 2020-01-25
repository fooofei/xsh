package base

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

var (
	Debug *log.Logger
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger

	logfileprefix = ConfigRootPath + "/logs/xsh-"
	logfilesuffix = ".log"
)

func initLog() {
	rotate()

	now := fmt.Sprintf("%s", time.Now().UTC())
	now = strings.Replace(now, " ", "_", -1)
	now = strings.Replace(now, ":", "-", -1)

	logfile := logfileprefix + now + logfilesuffix
	lf, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("open the log file [%s] error: %v\n", logfile, err)
	}

	Debug = log.New(lf, "D: ", log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(lf, "I: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warn = log.New(lf, "W: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(io.MultiWriter(os.Stderr, lf), "E: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func rotate() {
	if err := os.Mkdir(LogPath, os.ModeDir|0755); err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("mkdir %s error: %v\n", LogPath, err)
		}
	}

	lfs, err := filepath.Glob(logfileprefix + "*")
	if err != nil {
		log.Fatalf("glob the log file [%s] error: %v\n", logfileprefix, err)
	}
	if len(lfs) > 10 {
		sort.Sort(sort.Reverse(sort.StringSlice(lfs)))
		for _, value := range lfs[10:] {
			os.Remove(value)
		}
	}
}
