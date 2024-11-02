#!/bin/bash
# build.sh - Linux build script for music player (Windows target)

# Ensure script fails on any error
set -e

# Set Windows as target
export GOOS=windows
export GOARCH=amd64

# Create build directory structure
mkdir -p build/{font,mp3}

# Copy resources
cp -r font/* build/font/
cp -r mp3/* build/mp3/ 2>/dev/null || echo "No MP3s found - that's okay if you'll add them later"

# Build the executable
echo "Building Windows executable..."
go build -o build/SkyeMusicPlayer.exe main.go

echo "Build complete!"
echo "You can find the executable and resources in the 'build' directory"