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

RELEASE_FILE=RELEASE.md

LDFLAGS="-s -w \
    -X \"github.com/naggie/dstask.GIT_COMMIT=$GIT_COMMIT\" \
    -X \"github.com/naggie/dstask.VERSION=$VERSION\" \
    -X \"github.com/naggie/dstask.BUILD_DATE=$BUILD_DATE\"\
"

# get release information
if ! test -f $RELEASE_FILE || head -n 1 $RELEASE_FILE | grep -vq $VERSION; then
    # file doesn't exist or is for old version, replace
    printf "$VERSION\n\n\n" > $RELEASE_FILE
fi

vim "+ normal G $" $RELEASE_FILE


# build
mkdir -p dist

# UPX is disabled due to 40ms overhead, plus:
# see https://github.com/upx/upx/issues/222 -- UPX produces broken darwin executables.

GOOS=linux GOARCH=arm GOARM=5 go build -mod=vendor -ldflags="$LDFLAGS" cmd/dstask.go
# upx -q dstask
mv dstask dist/dstask-linux-arm5

GOOS=linux GOARCH=amd64 go build -mod=vendor -ldflags="$LDFLAGS" cmd/dstask.go
# upx -q dstask
mv dstask dist/dstask-linux-amd64

GOOS=darwin GOARCH=amd64 go build -mod=vendor -ldflags="$LDFLAGS" cmd/dstask.go
#upx -q dstask
mv dstask dist/dstask-darwin-amd64

hub release create \
    --draft \
    -a dist/dstask-linux-arm5#"dstask linux-arm5" \
    -a dist/dstask-linux-amd64#"dstask linux-amd64" \
    -a dist/dstask-darwin-amd64#"dstask darwin-amd64" \
    -F $RELEASE_FILE \
    $1
