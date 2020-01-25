package base

import "runtime"

const (
	Version    = "1.0.0"
	RootPath   = ".xsh"
	ConfigFile = "config.yaml"
	AuthFile   = "auth.yaml"
	HostFile   = "host.yaml"
	HisFile    = "xsh.his"
	EnvFile    = "xsh.env"
	TempPath   = ".xsh/temp"
	LogPath    = ".xsh/logs"
	PromptStr  = "[xsh]# "
)

var (
	Keywords = []string{":help", ":show", ":set", ":reload", ":do", ":sudo", ":copy"}
	setopts  = []string{"group=", "address="}
	showopts = []string{"group", "address", "env"}
)

var (
	GOOS = runtime.GOOS
)
