@echo off

:: Generate syso file
goversioninfo -platform-specific=true
set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=1
set CC=x86_64-w64-mingw32-gcc
go build -ldflags="-H windowsgui -s -w" -o booster_pump.exe

:: How to run build.bat
:: build.bat windows
:: build.bat linux
:: build.bat darwin