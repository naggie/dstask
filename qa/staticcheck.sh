#!/bin/sh
echo "installing staticcheck"
go get honnef.co/go/tools/cmd/staticcheck || exit 2
echo "running staticcheck ./..."
staticcheck ./...
