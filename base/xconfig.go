package base

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"path"
	"regexp"
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

	Concurrency         int      `yaml:"concurrency,omitempty"`
	CommonCommands      []string `yaml:"common_commands,omitempty"`
	BlackCommandRegexps []string `yaml:"black_command_regexps,omitempty"`
	CommandSep          string   `yaml:"command_sep,omitempty"`
}

var XConfig = config{}

func initXConfig() {
	setupXConfigDefault()

	var filePath = path.Join(ConfigRootPath, ConfigsFile)
	c, err := ioutil.ReadFile(filePath)
	if err != nil {
		Warn.Printf("Can not find configs file[%s].", filePath)
	}

	err = yaml.Unmarshal(c, &XConfig)
	if err != nil {
		log.Fatalf("Configs[%s] unmarshal error: %v", filePath, err)
	}

	Debug.Printf("XConfig: %+v", XConfig)
}

var XBlackCommandRegexps = make([]*regexp.Regexp, 0)

func setupXConfigDefault() {
	if XConfig.Timeout.TaskTimeoutS <= 0 {
		XConfig.Timeout.TaskTimeoutS = 21600
	}
	if XConfig.Timeout.ActionTimeoutS <= 0 {
		XConfig.Timeout.ActionTimeoutS = 3600
	}
	if XConfig.Timeout.CommandTimeoutS <= 0 {
		XConfig.Timeout.CommandTimeoutS = 300
	}
	if XConfig.Timeout.DialTimeoutS <= 0 {
		XConfig.Timeout.DialTimeoutS = 5
	}

	if XConfig.Ssh.Port <= 0 {
		XConfig.Ssh.Port = 22
	}
	if len(XConfig.Ssh.Ciphers) == 0 {
		XConfig.Ssh.Ciphers = []string{"aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-gcm@openssh.com", "arcfour256", "arcfour128", "aes128-cbc", "3des-cbc", "aes192-cbc", "aes256-cbc"}
	}
	if XConfig.Ssh.Pty.Width <= 0 {
		XConfig.Ssh.Pty.Width = 80
	}
	if XConfig.Ssh.Pty.Height <= 0 {
		XConfig.Ssh.Pty.Height = 60
	}
	if XConfig.Ssh.Pty.Ispeed <= 0 {
		XConfig.Ssh.Pty.Ispeed = 14400
	}
	if XConfig.Ssh.Pty.Ospeed <= 0 {
		XConfig.Ssh.Pty.Ospeed = 14400
	}
	if XConfig.Ssh.Pty.IntervalMS <= 0 {
		XConfig.Ssh.Pty.IntervalMS = 100
	}

	if XConfig.Cache.ExpirationS <= 0 {
		XConfig.Cache.ExpirationS = 900
	}
	if XConfig.Cache.CleanupIntervalS <= 0 {
		XConfig.Cache.CleanupIntervalS = 10
	}
	if XConfig.Cache.TickerIntervalS <= 0 {
		XConfig.Cache.TickerIntervalS = 10
	}

	if XConfig.Output.Type == "" {
		XConfig.Output.Type = "text"
	}

	if XConfig.Concurrency <= 0 {
		XConfig.Concurrency = 20
	}

	if len(XConfig.CommonCommands) == 0 {
		XConfig.CommonCommands = []string{"ls", "mkdir", "pwd", "cd", "cp", "mv", "cat", "grep", "find", "ping", "df", "ps"}
	}
	if len(XConfig.BlackCommandRegexps) == 0 {
		XConfig.BlackCommandRegexps = []string{"^\\s*(vi|vim)\\s+", "^\\s*top\\s*$"}
	}

	if XConfig.CommandSep == "" {
		XConfig.CommandSep = ";"
	}

	for _, blackCommandRegexp := range XConfig.BlackCommandRegexps {
		if blackCommand, err := regexp.Compile(blackCommandRegexp); err != nil {
			Warn.Printf("BlackCommandRegexp[%s] illegal, err: %v", blackCommandRegexp, err)
		} else {
			XBlackCommandRegexps = append(XBlackCommandRegexps, blackCommand)
		}
	}
}
