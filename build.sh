#!/bin/bash

rm -rf build

echo "Building for Windows"

env GOOS=windows GOARCH=386 go build -o build/windows.x86/m3udownloader.exe .
env GOOS=windows GOARCH=amd64 go build -o build/windows.amd64/m3udownloader.exe .

echo "Building for Linux"

env GOOS=linux GOARCH=386 go build -o build/linux.x86/m3udownloader .
env GOOS=linux GOARCH=amd64 go build -o build/linux.amd64/m3udownloader .
env GOOS=linux GOARCH=arm go build -o build/linux.arm/m3udownloader .
env GOOS=linux GOARCH=arm64 go build -o build/linux.arm64/m3udownloader .

echo "Building for OSX"

env GOOS=darwin GOARCH=amd64 go build -o build/mac.amd64/m3udownloader .
env GOOS=darwin GOARCH=arm64 go build -o build/mac.arm64/m3udownloader .

echo "Done"
