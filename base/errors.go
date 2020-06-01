package base

import "errors"

var (
	TimeoutErr = errors.New("timeout")

	ActionTimeoutErr = errors.New("do action timeout")
	ActionEmptyErr   = errors.New("action empty")

	CommandEmptyErr = errors.New("command can not empty")
	CommandSuErr    = errors.New("su info illegal")

	CopyDirectionErr = errors.New("copy direction illegal")

	LocalFileExistErr   = errors.New("local file can not exist")
	LocalDirExistErr    = errors.New("local dir can not exist")
	LocalDirNotEmptyErr = errors.New("local dir not empty")
	LocalPathNotFullErr = errors.New("local path must be full path")
	LocalPathIllegalErr = errors.New("local path name illegal")

	RemoteFileExistErr   = errors.New("remote file can not exist")
	RemoteDirExistErr    = errors.New("remote dir can not exist")
	RemoteDirNotEmptyErr = errors.New("remote dir not empty")
	RemotePathNotFullErr = errors.New("remote path must be full path")
	RemotePathIllegalErr = errors.New("remote path name illegal")

	CryptTypeUnknown = errors.New("crypt type unknown")
	CryptKeyIllegal  = errors.New("crypt key illegal")

	PromptAborted = errors.New("prompt aborted")
	PromptHisErr  = errors.New("prompt history error")
)
