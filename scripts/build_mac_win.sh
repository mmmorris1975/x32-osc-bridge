#!/usr/bin/env bash

source scripts/helpers.sh

export GOARCH="amd64"

# MacOS
export GOOS="darwin"
build $1

# Windows
export GOOS="windows"
build $1