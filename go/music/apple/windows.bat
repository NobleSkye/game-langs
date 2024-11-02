@echo off
setlocal

:: Set GOPATH if it's not already set
if "%GOPATH%"=="" (
    echo GOPATH is not set. Setting it to %USERPROFILE%\go
    set GOPATH=%USERPROFILE%\go
)

:: Add GOPATH\bin to PATH
set PATH=%GOPATH%\bin;%PATH%

:: Install dependencies
echo Installing dependencies...
go get -u github.com/hajimehoshi/ebiten/v2
go get -u github.com/sqweek/dialog
go get -u golang.org/x/image/font/opentype

:: Build the application
echo Building the application...
go build -o SkyesMusicPlayer.exe

:: Check if build was successful
if %ERRORLEVEL% neq 0 (
    echo Build failed. Please check the error messages above.
    pause
    exit /b 1
)

echo Build successful. SkyesMusicPlayer.exe has been created.
pause