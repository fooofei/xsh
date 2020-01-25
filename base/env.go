package base

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
)

//TargetType: group/address
type CurrentEnv struct {
	TargetType     string `yaml:"target_type,omitempty"`
	HostGroup      string `yaml:"host_group,omitempty"`
	HostAddress    string `yaml:"host_address,omitempty"`
	Authentication string `yaml:"authentication,omitempty"`
	PromptStr      string `yaml:"prompt_str,omitempty"`
	ActionType     string `yaml:"action_type,omitempty"`
	OutputType     string `yaml:"output_type,omitempty"`
}

var CurEnv CurrentEnv

func CheckEnv() bool {
	if CurEnv.TargetType == "" {
		return false
	}

	if CurEnv.TargetType == "group" {
		return CurEnv.HostGroup != ""
	}

	return CurEnv.HostAddress != "" && CurEnv.Authentication != ""
}

func SaveEnv() {
	if CheckEnv() {
		if CurEnv.TargetType == "group" {
			CurEnv.PromptStr = "[" + CurEnv.HostGroup + CurEnv.ActionType + "]# "
		} else {
			CurEnv.PromptStr = "[" + CurEnv.Authentication + "@" + CurEnv.HostAddress + CurEnv.ActionType + "]# "
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
		if CurEnv.TargetType == "group" {
			CurEnv.PromptStr = "[" + CurEnv.HostGroup + CurEnv.ActionType + "]# "
		} else {
			CurEnv.PromptStr = "[" + CurEnv.Authentication + "@" + CurEnv.HostAddress + CurEnv.ActionType + "]# "
		}
	} else {
		CurEnv.TargetType = ""
		CurEnv.PromptStr = PromptStr
		CurEnv.ActionType = ":do"
		CurEnv.OutputType = XConfig.Output.Type
	}
}
