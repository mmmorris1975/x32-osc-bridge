#!/usr/bin/env bash

source scripts/helpers.sh

cd ${1:-.}

for d in darwin* windows*
do
  pkg_zip $d
done
