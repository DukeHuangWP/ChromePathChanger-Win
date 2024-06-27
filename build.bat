set PATH=%PATH%;%cd%\GoVersionInfo
go generate
go build -ldflags -H=windowsgui -o=msedge.exe