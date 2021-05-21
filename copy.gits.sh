#!/bin/bash

echo "Copying gits"

cd ../publish

cd ..

cp -ra ./publishgits/linux.x86/screencorder/.git ./publish/linux.x86/
cp -ra ./publishgits/linux.amd64/screencorder/.git ./publish/linux.amd64/
cp -ra ./publishgits/linux.arm/screencorder/.git ./publish/linux.arm/
cp -ra ./publishgits/linux.arm64/screencorder/.git ./publish/linux.arm64/

