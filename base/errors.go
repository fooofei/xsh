package base

import "errors"

var (
	TimeoutErr = errors.New("timeout")

	ActionTimeoutErr = errors.New("do action timeout")
	ActionEmptyErr   = errors.New("action empty")

	CommandEmptyErr = errors.New("command can not empty")
	CommandSuErr    = errors.New("su info illegal")

	CopyDirectionErr = errors.New("copy direction illegal")

	LocalPathNotFullErr = errors.New("local path must be full path")
	LocalPathIllegalErr = errors.New("local path name illegal")

	RemotePathNotFullErr = errors.New("remote path must be full path")
	RemotePathIllegalErr = errors.New("remote path name illegal")

	CryptTypeUnknown = errors.New("crypt type unknown")
	CryptKeyIllegal  = errors.New("crypt key illegal")

	PromptAborted = errors.New("prompt aborted")
	PromptHisErr  = errors.New("prompt history error")
)
