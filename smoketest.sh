#!/bin/bash
# Just a quick smoke test to check nothing major is broken as a start of CI

# exit on error and print commands
set -x
set -e

export DSTASK_GIT_REPO=$(mktemp -d)
export DSTASK_CONTEXT_FILE=$(mktemp -u)

cleanup() {
    set +x
    set +e
    rm -rf $DSTASK_GIT_REPO
    rm $DSTASK_CONTEXT_FILE
}

trap cleanup EXIT

go build cmd/dstask.go

git -C $DSTASK_GIT_REPO init
git -C $DSTASK_GIT_REPO config user.email "you@example.com"
git -C $DSTASK_GIT_REPO config user.name "Test user"

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

# there should be no staged changes
git -C $DSTASK_GIT_REPO diff-index --quiet --cached HEAD --

# there should be no un-staged changes
git -C $DSTASK_GIT_REPO diff-files
git -C $DSTASK_GIT_REPO diff-files --quiet

# there should be no untracked files changes
test -z "$(git -C $DSTASK_GIT_REPO ls-files --others)"
