#!/usr/bin/env bash

source scripts/helpers.sh

cd ${1:-.}

for d in linux*
do
  # create zip archives (opposed to tar archives, for consistency with macos and windows packaging)
  pkg_zip $d
  # TODO use fpm to create DEB and RPM packages
done
