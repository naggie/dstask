package main

import (
	"github.com/naggie/dstask/dstask"
)

func main() {
	//dstask.Help()

	// importing requires full context
	ts := dstask.LoadTaskSetFromDisk(dstask.NORMAL_STATUSES)

	//err := ts.ImportFromTaskwarrior()

	//if err != nil {
	//	panic(err)
	//}
	ts.SortTaskList()
	ts.AssignIDs()
	ts.Display()
	ts.SaveToDisk()

	//fmt.Printf("%+v\n", dstask.ParseTaskLine(os.Args[1:len(os.Args)]))
}
