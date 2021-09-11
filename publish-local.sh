#!/bin/bash


echo "Copying static assests"

rm -rf ../publish-local
mkdir -p ../publish-local/www

cp -r js ../publish-local/www/
cp -r mirror ../publish-local/www/
cp -r imgs ../publish-local/www/
cp .gitignore ../publish-local/www/
cp update ../publish-local/www/
cp 404.html ../publish-local/www/
cp favicon.ico ../publish-local/www/
cp index.html ../publish-local/www/
cp screen.png ../publish-local/www/
cp serviceworker.js ../publish-local/www/
cp README.md ../publish-local/www/
mkdir ../publish-local/www/workspace


cp -ra ../publish-local/www ../publish-local/linux.x86
cp -ra ../publish-local/www ../publish-local/linux.amd64
cp -ra ../publish-local/www ../publish-local/linux.arm
cp -ra ../publish-local/www ../publish-local/linux.arm64


echo "Building for Linux"

env GOOS=linux GOARCH=386 go build -o ../publish-local/linux.x86/screencorder .
env GOOS=linux GOARCH=amd64 go build -o ../publish-local/linux.amd64/screencorder .
env GOOS=linux GOARCH=arm go build -o ../publish-local/linux.arm/screencorder .
env GOOS=linux GOARCH=arm64 go build -o ../publish-local/linux.arm64/screencorder .

echo "Publish local completed"