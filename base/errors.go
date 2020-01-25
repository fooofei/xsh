package base

import "errors"

var (
	TimeoutErr       = errors.New("timeout")
	ActionTimeoutErr = errors.New("do action timeout")
	CommandEmptyErr  = errors.New("command can not empty")
	CommandSuErr     = errors.New("command su info illegal")
	CryptTypeUnknown = errors.New("crypt type unknown")
	CryptKeyIllegal  = errors.New("crypt key illegal")
	PromptAborted    = errors.New("prompt aborted")
	PromptHisErr     = errors.New("prompt history error")
)
