package base

import "runtime"

const (
	Version        = "1.0.0"
	ConfigRootPath = ".xsh"
	CfgFile        = "cfg.yaml"
	AuthFile       = "auth.yaml"
	HostFile       = "host.yaml"
	HisFile        = "xsh.his"
	EnvFile        = "xsh.env"
	TempPath       = ".xsh/temp"
	LogPath        = ".xsh/logs"
	PromptStr      = "[xsh]# "
)

var (
	Keywords = []string{":help", ":show", ":set", ":reload", ":do", ":sudo", ":copy", ":encrypt", ":decrypt"}
	setopts  = []string{"group=", "address="}
	showopts = []string{"group", "address", "env"}
)

var (
	GOOS = runtime.GOOS
)
