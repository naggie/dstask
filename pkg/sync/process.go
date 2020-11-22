package sync

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/naggie/dstask"
	"gopkg.in/yaml.v2"
)

// ProcessTask handles a task that came in from a sync.Source (for now, that's only Github),
// merging it with the local copy as needed as defined in doc/dstask-sync.md
func ProcessTask(repo string, task dstask.Task) error {

	// note that locally, we may have the task as any state.
	// try to find it from any of the states, if found, load it and delete it, merge with Github, then save it again
	// this is quite naive but can be optimized later

	var found bool
	var localTask dstask.Task

	for _, status := range dstask.ALL_STATUSES {
		filepath := dstask.MustGetRepoPath(repo, status, task.UUID+".yml")

		// TODO differentiate between "does not exist" and "file exist but got an error while loading"
		// for now, we assume errors mean "do not exist"

		data, err := ioutil.ReadFile(filepath)
		if err != nil {
			continue
		}
		err = yaml.Unmarshal(data, &localTask)
		if err != nil {
			return fmt.Errorf("Failed to unmarshal %q: %s", filepath, err.Error())
		}
		found = true
		err = os.Remove(filepath)
		if err != nil {
			return err
		}
		break
	}
	if found {
		if localTask.Notes != "" {
			task.Notes = localTask.Notes
		}
		if task.Status == "pending" && (localTask.Status == "active" || localTask.Status == "paused") {
			task.Status = localTask.Status
		}
	}
	task.SaveToDisk(repo)
	return nil
}
