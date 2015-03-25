#!/usr/bin/env bash

if cmp -s "LICENSE.md" "$1/LICENSE.md"
then
   echo "not adding a license"
   exit 1
else
   echo "Adding a license"
   cp "LICENSE.md" "$1/LICENSE.md"
   exit 0
fi
