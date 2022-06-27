#!/bin/sh -l

key=$0

echo "Retrieving secret for $key"
secret="$(/scrt get $key)"

echo "::add-mask::$secret"
echo "::set-output name=secret::$secret"