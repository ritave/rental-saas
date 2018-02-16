#!/usr/bin/env bash

# $1 let it be URI
# $2 let it be a filename

curl $1 -d @$2 --header "Content-Type: application/json"
