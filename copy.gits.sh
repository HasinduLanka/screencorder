#!/bin/bash

echo "Copying gits"

cd ../screencorder-publish

cd ..

cp -ra ./screencorder-publishgits/linux.x86/screencorder/.git ./screencorder-publish/linux.x86/
cp -ra ./screencorder-publishgits/linux.amd64/screencorder/.git ./screencorder-publish/linux.amd64/
cp -ra ./screencorder-publishgits/linux.arm/screencorder/.git ./screencorder-publish/linux.arm/
cp -ra ./screencorder-publishgits/linux.arm64/screencorder/.git ./screencorder-publish/linux.arm64/

