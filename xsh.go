package main

import . "github.com/xied5531/xsh/base"

func main() {
	switch *RunMode {
	case "task":
		runTask()
	default:
		runPrompt()
	}
}
