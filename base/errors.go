package base

import "errors"

var (
	TimeoutErr = errors.New("timeout")

	ActionTimeoutErr = errors.New("do action timeout")
	ActionEmptyErr   = errors.New("action empty")

	CommandEmptyErr = errors.New("command can not empty")
	CommandSuErr    = errors.New("su info illegal")

	CopyInfoErr = errors.New("copy info illegal")

	LocalDirFormatIllegal = errors.New("local dir format illegal")
	LocalDirTypeIllegal   = errors.New("local dir type illegal")
	LocalDirNotEmptyErr   = errors.New("local dir not empty")
	LocalWalkErr          = errors.New("local walk error")

	RemoteDirFormatIllegal = errors.New("remote dir format illegal")
	RemoteDirTypeIllegal   = errors.New("remote dir type illegal")
	RemoteDirNotEmptyErr   = errors.New("remote dir not empty")
	RemoteWalkErr          = errors.New("remote walk error")

	CryptTypeUnknown = errors.New("crypt type unknown")
	CryptKeyIllegal  = errors.New("crypt key illegal")

	PromptAborted = errors.New("prompt aborted")
	PromptHisErr  = errors.New("prompt history error")
)
