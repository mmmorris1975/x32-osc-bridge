#!/usr/bin/env bash

source scripts/helpers.sh

# Only release if the commit is associated with a tag (maybe only matches a semver-like tag?)
git describe --tags --exact-match --first-parent || exit 0

wget -O - https://github.com/cli/cli/releases/download/v1.4.0/gh_1.4.0_linux_amd64.tar.gz | tar xzf - -C /var/tmp
GH=$(find /var/tmp/gh* -name gh)

cd ${1:-.}

sha256sum * >${NAME}_${VER}.sha256sum
cat *sha256sum

$GH release create $VER *.deb *.rpm *.zip *sha256sum -R mmmorris1975/$NAME -n "release $VER"