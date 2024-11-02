#!/bin/bash
# build.sh - Cross-platform build script for music player

# Ensure script fails on any error
set -e

# Check for target OS argument
if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <target_os>"
    echo "Where <target_os> is either 'linux' or 'windows'."
    exit 1
fi

TARGET_OS=$1

# Set target OS and architecture
if [ "$TARGET_OS" = "linux" ]; then
    unset GOOS
    unset GOARCH
    EXE_NAME="SkyeMusicPlayer"
elif [ "$TARGET_OS" = "windows" ]; then
    export GOOS=windows
    export GOARCH=amd64
    EXE_NAME="SkyeMusicPlayer.exe"
else
    echo "Invalid target OS: $TARGET_OS. Use 'linux' or 'windows'."
    exit 1
fi

# Create build directory structure
mkdir -p build/{font,mp3}

# Copy resources
cp -r font/* build/font/
cp -r mp3/* build/mp3/ 2>/dev/null || echo "No MP3s found - that's okay if you'll add them later"

# Build the executable
echo "Building $TARGET_OS executable..."
go build -o build/$EXE_NAME main.go

echo "Build complete!"
echo "You can find the executable and resources in the 'build' directory"
echo "Run it with: ./build/$EXE_NAME"
