package out

import (
	"encoding/json"
	"fmt"
	. "github.com/xied5531/xsh/base"
	. "github.com/xied5531/xsh/sh"
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
		if ar.Err != nil {
			fmt.Printf("%s\n", ar.Err.Error())
			return
		}
		for address, response := range ar.Result {
			fmt.Printf("[%-18s] ---------------------------------------------------------\n", address)
			for _, r := range response {
				if r.Stdout != "" {
					fmt.Printf("%s\n", r.Stdout)
				}
				if r.Stderr != "" {
					fmt.Printf("%s\n", "Warn =>")
					fmt.Printf("%s\n", r.Stderr)
				}
				if r.Err != nil {
					fmt.Printf("%s ", "Error =>")
					fmt.Printf("%s\n", r.Err.Error())
				}

				if r.Stdout == "" && r.Stderr == "" && r.Err == nil {
					fmt.Println()
				}
			}
			if len(response) > 1 {
				fmt.Printf("%s\n", "------")
			}
		}
	} else {
		ErrLogf("text error for result: %+v\n", ar)
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
