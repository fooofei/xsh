package out

import (
	"encoding/json"
	"fmt"
	. "github.com/xied5531/xsh/base"
	. "github.com/xied5531/xsh/sh"
	"gopkg.in/yaml.v2"
)

func Out(v interface{}) {
	if CurEnv.Output == "json" {
		PrintJson(v)
	} else if CurEnv.Output == "yaml" {
		PrintYaml(v)
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

		printActionResult(ar)
		return
	}

	if tr, ok := v.(SshTaskResult); ok {
		if tr.Err != nil {
			fmt.Printf("%s\n", tr.Err.Error())
			return
		}

		printTaskResult(tr)
		return
	}
}

func printTaskResult(result SshTaskResult) {
	if result.Err != nil {
		fmt.Printf("task[%s] error: %s\n", result.Name, result.Err.Error())
	}

	for _, result := range result.SshActionResults {
		fmt.Printf("[%-36s] =======================================\n", result.Name)
		printActionResult(result)
		fmt.Println()
	}
}

func printActionResult(result SshActionResult) {
	for address, response := range result.Result {
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
			if r.Status != nil {
				for _, s := range r.Status {
					fmt.Printf("%s\n", s)
				}
			}

			if r.Stdout == "" && r.Stderr == "" && r.Err == nil && r.Status == nil {
				fmt.Println()
			}
		}
		if len(response) > 1 {
			fmt.Printf("%s\n", "------")
		}
	}
}

func PrintJson(v interface{}) {
	d, _ := json.MarshalIndent(&v, "", "  ")
	fmt.Println(string(d))
}

func PrintYaml(v interface{}) {
	d, _ := yaml.Marshal(&v)
	fmt.Println(string(d))
}
