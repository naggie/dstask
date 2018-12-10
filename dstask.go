package main

import (
	"github.com/naggie/dstask/dstask"
)

func main() {
	// check git repository
	// do action
	ts := dstask.NewTaskSet()
	err := ts.ImportFromTaskwarrior()

	if (err != nil) {
		panic(err)
	}
	ts.Display()
	// commit to git repository
	// write bash completion cache
}
