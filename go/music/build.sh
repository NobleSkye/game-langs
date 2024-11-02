#!/bin/bash

# build_windows.sh

echo "Building Skye's Music Player for Windows..."

# Set GOOS and GOARCH for Windows
export GOOS=windows
export GOARCH=amd64

# Build the executable
go build -o SkyesMusicPlayer.exe

echo "Build complete. Output: SkyesMusicPlayer.exe"