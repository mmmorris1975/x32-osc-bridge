#!/usr/bin/env bash

source scripts/helpers.sh

# Only release if the commit is associated with a tag (maybe only matches a semver-like tag?)
git describe --tags --exact-match --first-parent || exit 0

cd ${1:-.}

sha256sum * >MANIFEST

# TODO - actually release something
cat MANIFEST