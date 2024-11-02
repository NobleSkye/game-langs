@echo off

:: Build for Windows
set GOOS=windows
set GOARCH=amd64
go build -o SkyeMusicPlayer_windows_amd64.exe

:: Build for macOS
set GOOS=darwin
set GOARCH=amd64
GOOS=darwin GOARCH=amd64 go build -o SkyesMusicPlayer main.go

:: Build for Linux
set GOOS=linux
set GOARCH=amd64
go build -o SkyesMusicPlayer_linux_amd64