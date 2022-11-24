#!/bin/bash
echo "Build for Windows..."
env GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC="x86_64-w64-mingw32-gcc" go build -o ./build/nes-emu-win.exe ./cmd/main.go
echo "Done!"

echo "Build for Linux..."
env GOOS=linux GOARCH=amd64 CGO_ENABLED=1 CC="x86_64-linux-musl-gcc" CXX="x86_64-linux-musl-g++" go build -o ./build/nes-emu-linux ./cmd/main.go
echo "Done!"

echo "Build for MacOS..."
env GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -o ./build/nes-emu-mac ./cmd/main.go
echo "Done!"

echo "Build for MacOS M1..."
go build -o ./build/nes-emu-mac-m1 ./cmd/main.go
echo "Done!"
