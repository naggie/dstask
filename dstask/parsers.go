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

type TaskWarriorAnnotation struct {
	Description string `json:"description"`
	Entry string `json: entry`
}

type TaskWarriorTask struct {
	Description string `json:"description"`
	End string `json:"end"`
	Entry string `json: entry`
	Modified string `json: modified`
	Status string `json: status`
	Tags []string `json: tags`
	Uuid string `json: uuid`
	Annotations []TaskWarriorAnnotation `json:annotations`
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
