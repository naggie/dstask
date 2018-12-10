package main

import (
	"github.com/naggie/dstask/dstask"
)

func main() {
	// check git repository
	// do action
	ts := dstask.LoadTaskSetFromDisk()
	err := ts.ImportFromTaskwarrior()

	if (err != nil) {
		panic(err)
	}
	ts.SortTaskList()
	ts.Display()
	// commit to git repository
	// write bash completion cache
}
