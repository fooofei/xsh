export GOOS=linux
export GOARCH=386
go build -o release/xsh-linux-386

export GOARCH=amd64
go build -o release/xsh-linux-amd64

export GOARCH=arm64
go build -o release/xsh-linux-arm64

export GOOS=windows
export GOARCH=386
go build -o release/xsh-windows-386.exe

export GOARCH=amd64
go build -o release/xsh-windows-amd64.exe
