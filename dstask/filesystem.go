package dstask

// an interface to the filesystem/git based database -- loading, saving, committing

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
)

func MustGetRepoDirectory(directory ...string) string {
	root := MustExpandHome(GIT_REPO)
	return path.Join(append([]string{root}, directory...)...)
}

func LoadTaskSetFromDisk(statuses []string) *TaskSet {
	ts := &TaskSet{
		knownUuids: make(map[string]bool),
		IDRoster:   LoadIDRoster(),
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
				uuid: uuid,
				status: status,
				originalFilepath: filepath,
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

			ts.AddTask(&t)
		}
	}

	return ts
}

func (t *Task) SaveToDisk() {
	filepath := MustGetRepoDirectory(t.status, t.uuid+".yml")
	//fmt.Println(filepath)
	d, err := yaml.Marshal(&t)
	//fmt.Println(string(d), err)

	err = ioutil.WriteFile(filepath, d, 0600)
	if err != nil {
		ExitFail("Failed to write task")
	}
}

// may be removed
func (ts *TaskSet) SaveToDisk() {
	for _, task := range ts.tasks {
		task.SaveToDisk()
	}

	ts.IDRoster.SaveToDisk()
}
