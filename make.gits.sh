#!/bin/bash

echo "Cleaning up"

rm -rf ../screencorder-publishgits
mkdir -p ../screencorder-publishgits

cd ../screencorder-publishgits

mkdir linux.x86
mkdir linux.amd64
mkdir linux.arm
mkdir linux.arm64

echo "Cloning"

cd ./linux.x86
git clone --depth=1 --branch=linux.x86 --single-branch https://github.com/HasinduLanka/screencorder

cd ../linux.amd64
git clone --depth=1 --branch=linux.amd64 --single-branch https://github.com/HasinduLanka/screencorder

cd ../linux.arm
git clone --depth=1 --branch=linux.arm --single-branch https://github.com/HasinduLanka/screencorder

cd ../linux.arm64
git clone --depth=1 --branch=linux.arm64 --single-branch https://github.com/HasinduLanka/screencorder
