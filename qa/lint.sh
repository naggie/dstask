#!/bin/bash

# find the dir we exist within...
DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
# and cd into root project dir
cd ${DIR}/..

if ! which golangci-lint &>/dev/null; then
    # run the install from a temp dir. we don't want 'go get' updating our go.mod/go.sum files
    dir=$(mktemp -d)
    cd $dir
    GO111MODULE=on go get 'github.com/golangci/golangci-lint/cmd/golangci-lint@v1.35.2'
    cd -
fi

exec golangci-lint run
