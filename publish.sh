#!/bin/bash

echo "Cloning github state"

./make.gits.sh


echo "Copying static assests"

rm -rf ../screencorder-publish
mkdir -p ../screencorder-publish/www

cp -r js ../screencorder-publish/www/
cp -r mirror ../screencorder-publish/www/
cp -r imgs ../screencorder-publish/www/
cp .gitignore ../screencorder-publish/www/
cp update ../screencorder-publish/www/
cp 404.html ../screencorder-publish/www/
cp favicon.ico ../screencorder-publish/www/
cp index.html ../screencorder-publish/www/
cp screen.png ../screencorder-publish/www/
cp serviceworker.js ../screencorder-publish/www/
cp README.md ../screencorder-publish/www/
mkdir ../screencorder-publish/www/workspace

cp -ra ../screencorder-publish/www ../screencorder-publish/linux.x86/
cp -ra ../screencorder-publish/www ../screencorder-publish/linux.amd64/
cp -ra ../screencorder-publish/www ../screencorder-publish/linux.arm/
cp -ra ../screencorder-publish/www ../screencorder-publish/linux.arm64/


echo "Making .gits"

./copy.gits.sh


echo "Building for Linux"

env GOOS=linux GOARCH=386 go build -o ../screencorder-publish/linux.x86/screencorder .
env GOOS=linux GOARCH=amd64 go build -o ../screencorder-publish/linux.amd64/screencorder .
env GOOS=linux GOARCH=arm go build -o ../screencorder-publish/linux.arm/screencorder .
env GOOS=linux GOARCH=arm64 go build -o ../screencorder-publish/linux.arm64/screencorder .

./pullgits.sh $1


echo "screencorder-publish completed"