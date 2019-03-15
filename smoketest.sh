#!/bin/bash
# Just a quick smoke test to check nothing major is broken as a start of CI

set -x
set -e

go build cmd/dstask.go

mkdir ~/.dstask
git -C ~/.dstask init

./dstask add test task +foo project:bar
./dstask start 1
./dstask stop 1
./dstask add another task
./dstask +foo
./dstask -foo
./dstask -project:foo
./dstask 2 done
./dstask 1 done
./dstask undo
./dstask log something
./dstask +foo
./dstask note 1 this is a note
./dstask context +foo
./dstask next
./dstask 1 done
./dstask resolved
./dstask show-projects
