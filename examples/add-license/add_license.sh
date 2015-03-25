#!/usr/bin/env bash

if cmp -s "LICENSE" "$1/LICENSE"
then
   echo "not adding a license"
   exit 1
else
   echo "Adding a license"
   cp "LICENSE" "$1/LICENSE"
   exit 0
fi
