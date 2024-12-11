#!/bin/bash

# Generate syso file
goversioninfo -platform-specific=true

GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -ldflags="-H windowsgui -s -w" -o booster_pump.exe