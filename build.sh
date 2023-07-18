#!/bin/bash

rm -f build/ffwn.zip
7z a build/ffwn.zip form/*
env GOOS=linux GOARCH=amd64 go build -o build/ffwn-checkout_linux -ldflags '-s -w'
env GOOS=windows GOARCH=amd64 go build -o build/ffwn-checkout.exe -ldflags '-s -w'
env GOOS=darwin GOARCH=amd64 go build -o build/ffwn-checkout_osx -ldflags '-s -w'
env GOOS=linux GOARCH=arm64 go build -o build/ffwn-checkout_rpi64 -ldflags '-s -w'
