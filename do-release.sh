#!/bin/sh

if [ $# -eq 0 ]; then
    echo "Usage: $0 <tag>"
    echo "Release version required as argument"
    exit 1
fi

mkdir -p dist

GOOS=linux GOARCH=arm GOARM=5 go build -mod=vendor -ldflags="-s -w" cmd/dstask.go
upx -q dstask
mv dstask dist/dstask-linux-arm5

GOOS=linux GOARCH=amd64 go build -mod=vendor -ldflags="-s -w" cmd/dstask.go
upx -q dstask
mv dstask dist/dstask-linux-amd64

GOOS=darwin GOARCH=amd64 go build -mod=vendor -ldflags="-s -w" cmd/dstask.go
upx -q dstask
mv dstask dist/dstask-darwin-amd64

hub release create \
    -a dist/dstask-linux-arm5#"dstask linux-arm5" \
    -a dist/dstask-linux-amd64#"dstask linux-amd64" \
    -a dist/dstask-darwin-amd64#"dstask darwin-amd64" \
    $1

rm -rf dist/tmp
