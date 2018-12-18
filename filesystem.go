package dstask

// an interface to the filesystem/git based database -- loading, saving, committing

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"time"
)

func MustGetRepoDirectory(directory ...string) string {
	root := MustExpandHome(GIT_REPO)
	return path.Join(append([]string{root}, directory...)...)
}

func LoadTaskSetFromDisk(statuses []string) *TaskSet {
	ts := &TaskSet{
		tasksByID:   make(map[int]*Task),
		tasksByUuid: make(map[string]*Task),
	}

	gitDotGitLocation := MustGetRepoDirectory(".git")

	if _, err := os.Stat(gitDotGitLocation); os.IsNotExist(err) {
		ExitFail("Could not find git repository at " + GIT_REPO + ", please clone or create")
	}

	for _, status := range statuses {
		dir := MustGetRepoDirectory(status)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err = os.Mkdir(dir, 0700)
			if err != nil {
				ExitFail("Failed to create directory in git repository")
			}
		}

		files, err := ioutil.ReadDir(dir)
		if err != nil {
			ExitFail("Failed to read " + dir)
		}

		for _, file := range files {
			filepath := path.Join(dir, file.Name())

			if len(file.Name()) != 40 {
				// not <uuid4>.yml
				continue
			}

			uuid := file.Name()[0:36]

			if !IsValidUuid4String(uuid) {
				continue
			}

			t := Task{
				Uuid:   uuid,
				Status: status,
			}

			data, err := ioutil.ReadFile(filepath)
			if err != nil {
				ExitFail("Failed to read " + filepath)
			}
			err = yaml.Unmarshal(data, &t)
			if err != nil {
				// TODO present error to user, specific error message is important
				ExitFail("Failed to parse " + filepath)
			}

			ts.AddTask(t)
		}
	}

	return ts
}

func (t *Task) SaveToDisk() {
	if !t.WritePending {
		return
	}

	t.Modified = time.Now()

	filepath := MustGetRepoDirectory(t.Status, t.Uuid+".yml")
	d, err := yaml.Marshal(&t)

	err = ioutil.WriteFile(filepath, d, 0600)
	if err != nil {
		ExitFail("Failed to write task")
	}

	// delete from all other locations to make sure there is only one copy
	// that exists
	for _, st := range ALL_STATUSES {
		if st == t.Status {
			continue
		}

		filepath := MustGetRepoDirectory(st, t.Uuid+".yml")

		if _, err := os.Stat(filepath); !os.IsNotExist(err) {
			err := os.Remove(filepath)
			if err != nil {
				ExitFail("Failed to delete " + filepath)
			}
		}
	}
}

// may be removed
func (ts *TaskSet) SaveToDisk(commitMsg string) {
	for _, task := range ts.tasks {
		task.SaveToDisk()
	}

	// git add all changed/created files
	// could optimise this to be given an explicit list of
	// added/modified/deleted files -- only if slow.
	MustRunGitCmd("add", ".")
	MustRunGitCmd("commit", "--no-gpg-sign", "-m", commitMsg)
}
