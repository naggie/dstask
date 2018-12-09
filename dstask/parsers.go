package dstask

import (
	"encoding/json"
	"os"
	"fmt"
)

func parseTaskLine(args []string) {

}

func parseFilterLine(args []string) {

}

func LoadTasks() *TaskSet {
	return &TaskSet{}
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

func (ts *TaskSet) ImportFromTaskwarrior() error {
	var taskWarriorTasks []TaskWarriorTask
	// from stdin
	err := json.NewDecoder(os.Stdin).Decode(&taskWarriorTasks)

	if (err != nil) {
		return err
	}

	fmt.Println(taskWarriorTasks)

	return nil
}
