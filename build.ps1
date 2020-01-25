$env:GOOS = "linux"
$env:GOARCH = "386"
go build -o bin\xsh-linux-386

$env:GOARCH = "amd64"
go build -o bin\xsh-linux-amd64

$env:GOARCH = "arm64"
go build -o bin\xsh-linux-arm64

$env:GOOS = "windows"
$env:GOARCH = "386"
go build -o bin\xsh-windows-386.exe

$env:GOARCH = "amd64"
go build -o bin\xsh-windows-amd64.exe
