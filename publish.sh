#!/bin/bash

echo "Cloning github state"

./make.gits.sh


echo "Copying static assests"

rm -rf ../publish
mkdir -p ../publish/www

cp -r js ../publish/www/
cp -r mirror ../publish/www/
cp .gitignore ../publish/www/
cp update ../publish/www/
cp 404.html ../publish/www/
cp favicon.ico ../publish/www/
cp index.html ../publish/www/
cp manifest.json ../publish/www/
cp screen.png ../publish/www/
cp serviceworker.js ../publish/www/
cp README.md ../publish/www/
mkdir ../publish/www/workspace

cp -ra ../publish/www ../publish/linux.x86/
cp -ra ../publish/www ../publish/linux.amd64/
cp -ra ../publish/www ../publish/linux.arm/
cp -ra ../publish/www ../publish/linux.arm64/


echo "Making .gits"

./copy.gits.sh


echo "Building for Linux"

env GOOS=linux GOARCH=386 go build -o ../publish/linux.x86/screencorder .
env GOOS=linux GOARCH=amd64 go build -o ../publish/linux.amd64/screencorder .
env GOOS=linux GOARCH=arm go build -o ../publish/linux.arm/screencorder .
env GOOS=linux GOARCH=arm64 go build -o ../publish/linux.arm64/screencorder .

./pullgits.sh $1


echo "Publish completed"