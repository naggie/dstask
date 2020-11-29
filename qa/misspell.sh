#!/bin/sh
echo "installing/updating misspell"
go get -u github.com/client9/misspell/cmd/misspell || exit 2
echo "running misspell"
misspell -error $(find . -type f | grep -v '^\./\.git' | grep -v '\./vendor')
