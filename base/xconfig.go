package base

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"path"
)

type config struct {
	Timeout struct {
		TaskTimeoutS    int `yaml:"task_timeout_s,omitempty"`
		ActionTimeoutS  int `yaml:"action_timeout_s,omitempty"`
		CommandTimeoutS int `yaml:"command_timeout_s,omitempty"`
		DialTimeoutS    int `yaml:"dial_timeout_s,omitempty"`
	} `yaml:"timeout,omitempty"`

	Ssh struct {
		Port    int      `yaml:"port,omitempty"`
		Ciphers []string `yaml:"ciphers,omitempty"`
		Pty     struct {
			Term       string `yaml:"term,omitempty"`
			Width      int    `yaml:"width,omitempty"`
			Height     int    `yaml:"height,omitempty"`
			Ispeed     int    `yaml:"ispeed,omitempty"`
			Ospeed     int    `yaml:"ospeed,omitempty"`
			IntervalMS int    `yaml:"interval_ms,omitempty"`
		} `yaml:"pty,omitempty"`
	} `yaml:"ssh,omitempty"`

	Cache struct {
		ExpirationS      int `yaml:"expiration_s,omitempty"`
		CleanupIntervalS int `yaml:"cleanup_interval_s,omitempty"`
		TickerIntervalS  int `yaml:"ticker_interval_s,omitempty"`
	} `yaml:"cache,omitempty"`

	Crypt struct {
		Type string `yaml:"type,omitempty"`
		Key  string `yaml:"key,omitempty"`
	} `yaml:"crypt,omitempty"`

	Output struct {
		Type    string `yaml:"type,omitempty"`
		Ordered bool   `yaml:"ordered,omitempty"`
	} `yaml:"output,omitempty"`

	Concurrency    int      `yaml:"concurrency,omitempty"`
	CommonCommands []string `yaml:"common_commands,omitempty"`
	CommandSep     string   `yaml:"command_sep,omitempty"`
}

var XConfig = config{}

func initXConfig() {
	var filePath = path.Join(ConfigRootPath, ConfigsFile)
	c, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Can not find configs file[%s].", filePath)
	}

	err = yaml.Unmarshal(c, &XConfig)
	if err != nil {
		log.Fatalf("Configs[%s] unmarshal error: %v", filePath, err)
	}
}
