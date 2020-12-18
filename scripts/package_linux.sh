#!/usr/bin/env bash

source scripts/helpers.sh

cd ${1:-.}

for d in linux*
do
  # create zip archives (opposed to tar archives, for consistency with macos and windows packaging)
  pkg_zip $d
  pkg_deb $d
  pkg_rpm $d
done
