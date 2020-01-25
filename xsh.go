package main

import . "github.com/xied5531/xsh/base"

func main() {
	switch *Mode {
	case "task":
		runTask()
	case "cmd":
		runCmd()
	case "copy":
		runCopy()
	case "prompt":
		runPrompt()
	default:
		Error.Printf("run mode[%s] not support", *Mode)
	}
}
