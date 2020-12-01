#!/bin/bash

out=$(gofmt -d -s $(find . -name '*.go' | grep -v vendor | grep -v _gen.go))
if [ "$out" != "" ]; then
	echo "$out"
	echo
	echo "You might want to run something like 'find . -name '*.go' -not -path './vendor/*' | xargs gofmt -w -s'"
	exit 2
fi
exit 0
