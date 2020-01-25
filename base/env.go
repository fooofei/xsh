package base

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
)

//Single: group/address
type CurrentEnv struct {
	Single  bool   `yaml:"single,omitempty"`
	Group   string `yaml:"group,omitempty"`
	Address string `yaml:"address,omitempty"`
	Auth    string `yaml:"auth,omitempty"`
	Prompt  string `yaml:"prompt,omitempty"`
	Action  string `yaml:"action,omitempty"`
	Output  string `yaml:"output,omitempty"`
}

var CurEnv CurrentEnv

func CheckEnv() bool {
	if !CurEnv.Single {
		return CurEnv.Group != ""
	}

	return CurEnv.Address != "" && CurEnv.Auth != ""
}

func SaveEnv() {
	if CheckEnv() {
		if !CurEnv.Single {
			CurEnv.Prompt = "[" + CurEnv.Group + CurEnv.Action + "]# "
		} else {
			CurEnv.Prompt = "[" + CurEnv.Auth + "@" + CurEnv.Address + CurEnv.Action + "]# "
		}

		d, err := yaml.Marshal(&CurEnv)
		if err != nil {
			Warn.Printf("marshal env error: %v", err)
		}

		if err := ioutil.WriteFile(path.Join(ConfigRootPath, EnvFile), d, os.ModeAppend|0644); err != nil {
			Warn.Printf("save env error: %v", err)
		}
	} else {
		Warn.Printf("save current env error, env: %v", CurEnv)
		initEnv()
		Warn.Printf("after reset env, env: %v", CurEnv)
	}
}

func initEnv() {
	CurEnv.Single = false
	CurEnv.Prompt = PromptStr
	CurEnv.Action = ":do"
	CurEnv.Output = XConfig.Output.Type

	e, err := ioutil.ReadFile(path.Join(ConfigRootPath, EnvFile))
	if err != nil {
		Warn.Printf("read current env error: %v", err)
		return
	}

	if yaml.Unmarshal(e, &CurEnv) != nil {
		Warn.Printf("load current env error: %v", err)
		return
	}

	if CheckEnv() {
		if !CurEnv.Single {
			CurEnv.Prompt = "[" + CurEnv.Group + CurEnv.Action + "]# "
		} else {
			CurEnv.Prompt = "[" + CurEnv.Auth + "@" + CurEnv.Address + CurEnv.Action + "]# "
		}
	}
}
