rd /s /q bin
mkdir bin
set GOARCH=arm64
set GOOS=linux
go build -o bin/gateway
