#!/bin/bash

cd ../git/

git branch -m linux.x86
git switch linux.x86
git reset --hard 
git config pull.ff only
git pull

git branch -m linux.amd64
git switch linux.amd64
git reset --hard 
git config pull.ff only
git pull

git branch -m linux.arm
git switch linux.arm
git reset --hard 
git config pull.ff only
git pull

git branch -m linux.arm64
git switch linux.arm64
git reset --hard 
git config pull.ff only
git pull

