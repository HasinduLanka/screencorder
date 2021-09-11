#!/bin/bash

echo "Syncing with github"

cd ../screencorder-publish/


cd ./linux.x86
git add .
git update-index --add --chmod=+x screencorder
git update-index --add --chmod=+x update
git commit -m "build-$1"
git push origin HEAD

cd ../linux.amd64
git add .
git update-index --add --chmod=+x screencorder
git update-index --add --chmod=+x update
git commit -m "build-$1"
git push origin HEAD

cd ../linux.arm
git add .
git update-index --add --chmod=+x screencorder
git update-index --add --chmod=+x update
git commit -m "build-$1"
git push origin HEAD

cd ../linux.arm64
git add .
git update-index --add --chmod=+x screencorder
git update-index --add --chmod=+x update
git commit -m "build-$1"
git push origin HEAD
