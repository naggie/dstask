#!/bin/bash
# Just a quick smoke test to check nothing major is broken as a start of CI

set -x
set -e

go build cmd/dstask.go

mkdir ~/.dstask
git -C ~/.dstask init

git -C ~/.dstask config user.email "you@example.com"
git -C ~/.dstask config user.name "Test user"

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

# we are in context project:bar, adding with another project should fail
./dstask context project:bar
! ./dstask add project:baz test
# ... however, bypassing the context with -- should work
./dstask add project:cheese test --

./dstask context none
./dstask context

./dstask import-tw < etc/taskwarrior-export.json
./dstask next
