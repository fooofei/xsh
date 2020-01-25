package main

import (
	. "github.com/xied5531/xsh/out"
	. "github.com/xied5531/xsh/sh"
)

func runTask() {
	task := SshTask{}
	Out(task.Do())
}
