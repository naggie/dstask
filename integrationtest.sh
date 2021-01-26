#!/bin/bash
# Just a quick smoke test to check nothing major is broken as a start of CI

# This script should only test commands work without crashing. Behavioural
# tests should use the go test system.

# pipe into cat to disable tty

# exit on error and print commands
set -x
set -e

# isolated db locations (repo2 is used for a sync target)
export DSTASK_GIT_REPO=$(mktemp --directory)
export UPSTREAM_BARE_REPO=$(mktemp --directory)

if [[ -d "dstask" ]]; then
    rm -r dstask
fi

cleanup() {
    set +x
    set +e
    rm -rf $DSTASK_GIT_REPO
    rm -rf $UPSTREAM_BARE_REPO
    rm dstask
}

trap cleanup EXIT

go build -o dstask -mod=vendor cmd/dstask/main.go

# initialise git repo
BRANCH="branch_${RANDOM}"
git -C $DSTASK_GIT_REPO init
git -C $DSTASK_GIT_REPO checkout -b $BRANCH
git -C $DSTASK_GIT_REPO config user.email "you@example.com"
git -C $DSTASK_GIT_REPO config user.name "Test user"
git -C $UPSTREAM_BARE_REPO init --bare


# general task state management and commands
./dstask help
./dstask add test task +foo project:bar
./dstask start 1
./dstask stop 1
./dstask remove 1
./dstask add re-add add test task +foo project:bar
./dstask add another task
./dstask +foo
./dstask -foo
./dstask -project:foo
./dstask 2 done
./dstask undo
./dstask log something
./dstask +foo
./dstask note 1 this is a note
./dstask context +foo
./dstask next
./dstask 1 done
./dstask show-resolved

# TODO we set FAKE_PTY because we do not have a non-tty
# rendering scheme for show-projects
export DSTASK_FAKE_PTY=1
./dstask show-projects
unset DSTASK_FAKE_PTY

# -- to remove current context which is +foo
./dstask add -- unorganised task
./dstask show-unorganised

# we are in context project:bar, adding with another project should fail
./dstask context project:bar
! ./dstask add project:baz test
# ... however, bypassing the context with -- should work
./dstask add project:cheese test --

# test context listing with a context and no context
./dstask context
./dstask context none
./dstask context

# try to resolve a task with an incomplete tasklist: should fail
./dstask add a tasklist task -- "- [ ] incomplete task"
! ./dstask 2 done

# test template functions
./dstask add task to copy
./dstask add template:2 +copiedTask
./dstask template 5
# Task 5 should now be a template
./dstask show-templates +copiedTask
# copy Template with some modifications
./dstask add template:5 -copiedTask +copiedTemplate
./dstask show-open +copiedTemplate
# Create new template from CMD line.
./dstask template give me some things to do P1 +uniqueTag
./dstask show-templates +uniqueTag

# test import
./dstask import-tw < etc/taskwarrior-export.json
./dstask next

# test git command pass through
./dstask git status

# set the bare repository as upstream origin to test sync against, and then push to it
git -C $DSTASK_GIT_REPO remote add origin $UPSTREAM_BARE_REPO
git -C $DSTASK_GIT_REPO push origin $BRANCH
git -C $DSTASK_GIT_REPO branch --set-upstream-to=origin/$BRANCH $BRANCH
./dstask sync

# cause independent changes (could be from separate downstream repositories,
# but its easier and equivalent to simulate this with a hard reset)
./dstask add eggs test task +foo project:bar
./dstask sync
./dstask git reset --hard HEAD~1
./dstask add bacon test task +foo project:bar
./dstask sync

# there should be no staged changes
git -C $DSTASK_GIT_REPO diff-index --quiet --cached HEAD --

# there should be no un-staged changes
git -C $DSTASK_GIT_REPO diff-files --quiet

# there should be no untracked files changes
test -z "$(git -C $DSTASK_GIT_REPO ls-files --others)"

# regression test: nil pointer dereference when showing by-week tables of zero length
./dstask show-resolved project:doesnotexist

# dstask should be able to create a git repo if it does not exist before
# executing command. Also help should work before initialisation (and should
# not init)
export DSTASK_GIT_REPO=$(mktemp --directory --dry-run)
./dstask help
./dstask add test
