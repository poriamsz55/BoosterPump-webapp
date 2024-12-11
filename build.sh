#!/bin/bash

# Generate syso file
goversioninfo -platform-specific=true

# Get build os argument
build_os=$1

# Build booster_pump
# if build_os is windows
if [ "$build_os" == "windows" ]; then
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -ldflags="-H windowsgui -s -w" -o booster_pump.exe

# if build_os is linux
elif [ "$build_os" == "linux" ]; then
GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o booster_pump

# if build_os is darwin
elif [ "$build_os" == "darwin" ]; then
GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -o booster_pump
fi

# How to run build.sh
# ./build.sh windows
# ./build.sh linux
# ./build.sh darwin