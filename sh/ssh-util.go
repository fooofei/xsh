package sh

import (
	"fmt"
	"github.com/xied5531/xsh/base"
	"time"
)

func withTimeout(fn func(interface{}) sshResponse, arg interface{}, timeout time.Duration) sshResponse {
	result := make(chan sshResponse, 1)
	go func() {
		defer close(result)
		result <- fn(arg)
	}()

	select {
	case r := <-result:
		return r
	case <-time.After(timeout):
		return sshResponse{Err: base.TimeoutErr}
	}
}

func checkCommands(commands []string) error {
	if len(base.XBlackCommandRegexps) > 0 {
		for _, command := range commands {
			for _, reg := range base.XBlackCommandRegexps {
				if cmd := reg.FindStringIndex(command); cmd != nil {
					return fmt.Errorf("command[%s] in black command list", command)
				}
			}
		}
	}

	return nil
}
