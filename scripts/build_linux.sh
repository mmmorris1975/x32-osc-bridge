#!/usr/bin/env bash

set -e
set -o pipefail

source scripts/helpers.sh

export GOOS="linux"

# amd64
export GOARCH="amd64"
build $1

# ARMv7
export GOARCH="arm"
export GOARM="7"
build $1
unset GOARM

# ARM64
export GOARCH="arm64"
build $1
