package main

import (
	"github.com/naggie/dstask/dstask"
)

func main() {
	// importing requires full context
	ts := dstask.LoadTaskSetFromDisk(FULL_CONTEXT)

	err := ts.ImportFromTaskwarrior()

	if (err != nil) {
		panic(err)
	}
	ts.SortTaskList()
	ts.Display()
}
