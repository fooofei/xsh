package base

import "errors"

var (
	TimeoutErr = errors.New("timeout")

	ActionTimeoutErr = errors.New("do action timeout")
	ActionEmptyErr   = errors.New("action empty")

	CommandEmptyErr = errors.New("command can not empty")
	CommandSuErr    = errors.New("command su info illegal")

	CopyInfoErr        = errors.New("copy info illegal")
	RemoteFileExistErr = errors.New("remote file existed")
	LocalFileExistErr  = errors.New("local file existed")
	RemoteWalkErr      = errors.New("remote walk existed")

	CryptTypeUnknown = errors.New("crypt type unknown")
	CryptKeyIllegal  = errors.New("crypt key illegal")

	PromptAborted = errors.New("prompt aborted")
	PromptHisErr  = errors.New("prompt history error")
)
