package dstask

import (
	"encoding/json"
	"os"
)

func parseTaskLine(args []string) {

}

func parseFilterLine(args []string) {

}

func LoadTasks() *TaskSet {

}

func parseFile(filepath string) {

}

type TaskWarriorTask struct {
	Description string `json:"description"`
	End string `json:"end"`
	Entry string `json: entry`
	Id int `json: id`
	Modified string `json: modified`
	Status string `json: status`
	Tags []string `json: tags`
	Uuid string `json: uuid`
}

func (ts *TaskSet) ImportFromTaskwarrior() {
	var taskWarriorTasks []TaskWarriorTask
	// from stdin
	err := json.NewDecoder(os.Stdin).Decode(&taskWarriorTasks)

	if (err != nil) {
		return err
	}
}
