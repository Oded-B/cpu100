@echo off
go env -w GOOS=windows GOARCH=amd64
cd /d %GOPATH%\src\github.com\schwarzlichtbezirk\cpu100
go build -o %GOPATH%/bin/cpu100.x64.exe -v github.com/schwarzlichtbezirk/cpu100
