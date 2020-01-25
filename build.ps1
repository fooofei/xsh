$env:GOOS = "linux"
$env:GOARCH = "386"
go build -o release\xsh-linux-386

$env:GOARCH = "amd64"
go build -o release\xsh-linux-amd64

$env:GOARCH = "arm64"
go build -o release\xsh-linux-arm64

$env:GOOS = "windows"
$env:GOARCH = "386"
go build -o release\xsh-windows-386.exe

$env:GOARCH = "amd64"
go build -o release\xsh-windows-amd64.exe
