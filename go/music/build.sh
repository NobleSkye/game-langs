#!/bin/bash

# build_windows.sh

# Ensure GOPATH is set
if [ -z "$GOPATH" ]; then
    echo "GOPATH is not set. Setting it to $HOME/go"
    export GOPATH=$HOME/go
fi

# Set GOOS and GOARCH for Windows
export GOOS=windows
export GOARCH=amd64

# Create a build directory
mkdir -p build

# Compile the program
echo "Compiling Skye's Music Player for Windows..."
go build -o build/SkyesMusicPlayer.exe

# Copy necessary files
echo "Copying font and music files..."
mkdir -p build/font
cp font/opendyslexic.ttf build/font/
mkdir -p build/music
cp music/*.mp3 build/music/

echo "Build complete. Check the 'build' directory for the executable and resources."