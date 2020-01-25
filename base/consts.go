package base

const (
	ConfigRootPath      string = ".xsh"
	ConfigsFile         string = "configs.yaml"
	AuthenticationsFile string = "authentications.yaml"
	HostgroupsFile      string = "hostgroups.yaml"
	TempFile            string = "xsh.temp"
	HisFile             string = "xsh.his"
	EnvFile             string = "xsh.env"
	LogPath             string = ".xsh/logs"
	LogoStr             string = "xsh"
	PromptStr           string = "[" + LogoStr + "]# "
)

var (
	Keywords = []string{":show", ":set", ":exit", ":do", ":sudo"}
	setopts  = []string{"group=", "address="}
	showopts = []string{"address"}
)
