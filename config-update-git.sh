#!/bin/bash

cd ../git/

git switch linux.x86
git reset --hard 
git config pull.ff only
git pull

git switch linux.amd64
git reset --hard 
git config pull.ff only
git pull

git switch linux.arm
git reset --hard 
git config pull.ff only
git pull

git switch linux.arm64
git reset --hard 
git config pull.ff only
git pull

