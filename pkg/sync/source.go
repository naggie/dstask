package sync

import "github.com/naggie/dstask"

type Source interface {
	// Next gets the next batch of tasks, or an error
	// caller should keep calling this until an empty slice is returned
	Next() ([]dstask.Task, error)
}
