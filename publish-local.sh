#!/bin/bash


echo "Copying static assests"

rm -rf ../screencorder-publish-local
mkdir -p ../screencorder-publish-local/www

cp -r js ../screencorder-publish-local/www/
cp -r imgs ../screencorder-publish-local/www/
cp .gitignore ../screencorder-publish-local/www/
cp update ../screencorder-publish-local/www/
cp 404.html ../screencorder-publish-local/www/
cp favicon.ico ../screencorder-publish-local/www/
cp index.html ../screencorder-publish-local/www/
cp screen.png ../screencorder-publish-local/www/
cp serviceworker.js ../screencorder-publish-local/www/
cp README.md ../screencorder-publish-local/www/
mkdir ../screencorder-publish-local/www/workspace


cp -ra ../screencorder-publish-local/www ../screencorder-publish-local/linux.x86
cp -ra ../screencorder-publish-local/www ../screencorder-publish-local/linux.amd64
cp -ra ../screencorder-publish-local/www ../screencorder-publish-local/linux.arm
cp -ra ../screencorder-publish-local/www ../screencorder-publish-local/linux.arm64


echo "Building for Linux"

env CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o ../screencorder-publish-local/linux.x86/screencorder .
env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../screencorder-publish-local/linux.amd64/screencorder .
env CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o ../screencorder-publish-local/linux.arm/screencorder .
env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ../screencorder-publish-local/linux.arm64/screencorder .

echo "screencorder-publish local completed"