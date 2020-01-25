package out

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	. "github.com/luckywinds/xsh/base"
	. "github.com/luckywinds/xsh/sh"
	"gopkg.in/yaml.v2"
)

func Print(v interface{}) {
	if CurEnv.OutputType == "json" {
		printJson(v)
	} else if CurEnv.OutputType == "yaml" {
		printYaml(v)
	} else {
		printText(v)
	}
}

func printText(v interface{}) {
	if ar, ok := v.(SshActionResult); ok {
		for _, c := range ar.Result {
			color.Green("[%-18s] ---------------------------------------------------------\n", c.Address)
			if c.Stdout != "" {
				fmt.Printf("%s\n", c.Stdout)
			}
			if c.Stderr != "" {
				color.Red("%s\n", "Warn =>")
				fmt.Printf("%s\n", c.Stderr)
			}
			if c.Err != nil {
				color.Red("%s\n", "Error =>")
				fmt.Printf("%s\n", c.Err.Error())
			}
		}
	} else {
		Error.Printf("text error for result: %+v\n", ar)
	}
}

func printJson(v interface{}) {
	d, _ := json.MarshalIndent(&v, "", "  ")
	fmt.Println(string(d))
}

func printYaml(v interface{}) {
	d, _ := yaml.Marshal(&v)
	fmt.Println(string(d))
}
