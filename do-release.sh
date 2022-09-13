#!/bin/sh
set -e

if [ $# -eq 0 ]; then
    echo "Usage: $0 <tag>"
    echo "Release version required as argument"
    exit 1
fi

VERSION="$1"
GIT_COMMIT=$(git rev-list -1 HEAD)
BUILD_DATE=$(date)
export CGO_ENABLED=0

RELEASE_FILE=RELEASE.md

LDFLAGS="-s -w \
    -X \"github.com/naggie/dstask.GIT_COMMIT=$GIT_COMMIT\" \
    -X \"github.com/naggie/dstask.VERSION=$VERSION\" \
    -X \"github.com/naggie/dstask.BUILD_DATE=$BUILD_DATE\"\
"

# build
mkdir -p dist

# UPX is disabled due to 40ms overhead, plus:
# see https://github.com/upx/upx/issues/222 -- UPX produces broken darwin executables.

GOOS=linux GOARCH=arm GOARM=5 go build -o dstask -mod=vendor -ldflags="$LDFLAGS" cmd/dstask/main.go
GOOS=linux GOARCH=arm GOARM=5 go build -o dstask-import -mod=vendor -ldflags="$LDFLAGS" cmd/dstask-import/main.go
# upx -q dstask
mv dstask dist/dstask-linux-arm5
mv dstask-import dist/dstask-import-linux-arm5

GOOS=linux GOARCH=amd64 go build -o dstask -mod=vendor -ldflags="$LDFLAGS" cmd/dstask/main.go
GOOS=linux GOARCH=amd64 go build -o dstask-import -mod=vendor -ldflags="$LDFLAGS" cmd/dstask-import/main.go
# upx -q dstask
mv dstask dist/dstask-linux-amd64
mv dstask-import dist/dstask-import-linux-amd64

GOOS=darwin GOARCH=amd64 go build -o dstask -mod=vendor -ldflags="$LDFLAGS" cmd/dstask/main.go
GOOS=darwin GOARCH=amd64 go build -o dstask-import -mod=vendor -ldflags="$LDFLAGS" cmd/dstask-import/main.go
#upx -q dstask
mv dstask dist/dstask-darwin-amd64
mv dstask-import dist/dstask-import-darwin-amd64

# github.com/cli/cli
# https://github.com/cli/cli/releases/download/v2.15.0/gh_2.15.0_linux_amd64.deb
# do: gh auth login
gh release create \
    --title $VERSION \
    --notes-file $RELEASE_FILE \
    --draft \
    $VERSION \
    dist/dstask-linux-arm5#"dstask linux-arm5" \
    dist/dstask-linux-amd64#"dstask linux-amd64" \
    dist/dstask-darwin-amd64#"dstask darwin-amd64" \
    dist/dstask-import-linux-arm5#"dstask-import linux-arm5" \
    dist/dstask-import-linux-amd64#"dstask-import linux-amd64" \
    dist/dstask-import-darwin-amd64#"dstask-import darwin-amd64" \
