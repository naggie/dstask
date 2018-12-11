package main

import (
	"github.com/naggie/dstask/dstask"
)

func main() {
	// importing requires full context
	ts := dstask.LoadTaskSetFromDisk(dstask.ALL_STATUSES)

	err := ts.ImportFromTaskwarrior()

	if err != nil {
		panic(err)
	}
	ts.SortTaskList()
	//ts.Display()
	ts.SaveToDisk()
}
