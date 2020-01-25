package sh

import (
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
