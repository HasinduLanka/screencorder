#!/bin/bash

rm -rf build
mkdir build

echo "Building for Linux"

env GOOS=linux GOARCH=386 go build -o build/linux.x86/m3udownloader .
env GOOS=linux GOARCH=amd64 go build -o build/linux.amd64/m3udownloader .
env GOOS=linux GOARCH=arm go build -o build/linux.arm/m3udownloader .
env GOOS=linux GOARCH=arm64 go build -o build/linux.arm64/m3udownloader .

echo "Done"
