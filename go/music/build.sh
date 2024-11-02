#!/bin/bash

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o musicplayer_windows_amd64.exe

# Build for macOS
GOOS=darwin GOARCH=amd64 go build -o musicplayer_macos_amd64

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o musicplayer_linux_amd64