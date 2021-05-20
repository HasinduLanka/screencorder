#!/bin/bash

echo Commit $1

echo "Cleaning up"

rm -rf ../build
mkdir -p ../build/www

echo "Scanning static assests"

cp -r js ../build/www/
cp -ra .git ../build/www/
cp .gitignore ../build/www/
cp update ../build/www/
cp 404.html ../build/www/
cp favicon.ico ../build/www/
cp index.html ../build/www/
cp manifest.json ../build/www/
cp screen.png ../build/www/
cp serviceworker.js ../build/www/

echo "Copying static assests"

cp -ra ../build/www ../build/linux.x86
cp -ra ../build/www ../build/linux.amd64
cp -ra ../build/www ../build/linux.arm
cp -ra ../build/www ../build/linux.arm64

echo "Building for Linux"

env GOOS=linux GOARCH=386 go build -o ../build/linux.x86/m3udownloader .
env GOOS=linux GOARCH=amd64 go build -o ../build/linux.amd64/m3udownloader .
env GOOS=linux GOARCH=arm go build -o ../build/linux.arm/m3udownloader .
env GOOS=linux GOARCH=arm64 go build -o ../build/linux.arm64/m3udownloader .

echo "Configuring GITs"

cd ../build/linux.x86
git branch -m linux.x86
git switch linux.x86
git add .
git update-index --add --chmod=+x m3udownloader
git update-index --add --chmod=+x update
git commit -m "build-$1"
mkdir workspace
git push origin HEAD

cd ../linux.amd64
git branch -m linux.amd64
git switch linux.amd64
git add .
git update-index --add --chmod=+x m3udownloader
git update-index --add --chmod=+x update
git commit -m "build-$1"
mkdir workspace
git push origin HEAD

cd ../linux.arm
git branch -m linux.arm
git switch linux.arm
git add .
git update-index --add --chmod=+x m3udownloader
git update-index --add --chmod=+x update
git commit -m "build-$1"
mkdir workspace
git push origin HEAD

cd ../linux.arm64
git branch -m linux.arm64
git switch linux.arm64
git add .
git update-index --add --chmod=+x m3udownloader
git update-index --add --chmod=+x update
git commit -m "build-$1"
mkdir workspace
git push origin HEAD

echo "Done"
