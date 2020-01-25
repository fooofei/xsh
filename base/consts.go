package base

const (
	ConfigRootPath      = ".xsh"
	ConfigsFile         = "configs.yaml"
	AuthenticationsFile = "authentications.yaml"
	HostgroupsFile      = "hostgroups.yaml"
	TempFile            = "xsh.temp"
	HisFile             = "xsh.his"
	EnvFile             = "xsh.env"
	LogPath             = ".xsh/logs"
	LogoStr             = "xsh"
	PromptStr           = "[" + LogoStr + "]# "
)

var (
	Keywords = []string{":help", ":show", ":set", ":reload", ":do", ":sudo", ":encrypt", ":decrypt"}
	setopts  = []string{"group=", "address="}
	showopts = []string{"address", "env"}
)
