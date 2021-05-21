#!/bin/bash


echo "Configuring GITs"

cd ../build/linux.x86
git reset --hard 
mkdir -p workspace
git branch -m linux.x86
git switch linux.x86
git reset --hard 
git pull

cd ../linux.amd64
git reset --hard 
mkdir -p workspace
git branch -m linux.amd64
git switch linux.amd64
git reset --hard 
git pull

cd ../linux.arm
git reset --hard 
mkdir -p workspace
git branch -m linux.arm
git switch linux.arm
git reset --hard 
git pull

cd ../linux.arm64
git reset --hard 
mkdir -p workspace
git branch -m linux.arm64
git switch linux.arm64
git reset --hard 
git pull
