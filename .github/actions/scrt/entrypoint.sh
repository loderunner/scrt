#!/bin/sh -l
set -e

key=$1

echo "Retrieving secret for \"$key\""
secret="$(/scrt get $key)"

echo "::add-mask::$secret"
echo "::set-output name=secret::$secret"