package main

import (
	"github.com/naggie/dstask/dstask"
	"os"
	"fmt"
)

func main() {
	dstask.Help()

	// importing requires full context
	ts := dstask.LoadTaskSetFromDisk(dstask.ALL_STATUSES)

	//err := ts.ImportFromTaskwarrior()

	//if err != nil {
	//	panic(err)
	//}
	ts.SortTaskList()
	ts.Display()
	ts.SaveToDisk()

	fmt.Printf("%+v\n", dstask.ParseTaskLine(os.Args[1:len(os.Args)]))
}
