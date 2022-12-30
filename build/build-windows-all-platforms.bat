ECHO OFF
CLS
:MENU
ECHO.
ECHO ...............................................
ECHO Terminator-Samurai-VPN
ECHO ...............................................
ECHO.
ECHO 1 - Linux amd64
ECHO 2 - Linux arm64
ECHO 3 - OSx amd64
ECHO 4 - Windows amd64
ECHO.
SET /P M=Select OS Type then press ENTER:
IF %M%==1 GOTO LINAMD
IF %M%==2 GOTO LINARM
IF %M%==3 GOTO OSX
IF %M%==4 GOTO WIN
:LINAMD
go env -w GOARCH=amd64
go env -w GOOS=linux
go build -o .\bin\TSVPN-linux-amd64 .\main.go
GOTO MENU
:LINARM
go env -w GOARCH=arm64
go env -w GOOS=linux
go build -o .\bin\TSVPN-linux-arm64 .\main.go
GOTO MENU
:OSX
go env -w GOARCH=amd64
go env -w GOOS=darwin
go build -o .\bin\TSVPN-darwin-amd64 .\main.go
GOTO MENU
:WIN
go env -w GOARCH=amd64
go env -w GOOS=windows
go build -o .\bin\TSVPN-win-amd64.exe .\main.go
GOTO MENU