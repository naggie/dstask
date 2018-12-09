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

type TWAnnotation struct {
	Description string `json:"description"`
	Entry string `json: entry`
}

type TWTask struct {
	Description string `json:"description"`
	End string `json:"end"`
	Entry string `json: entry`
	Modified string `json: modified`
	Status string `json: status`
	Tags []string `json: tags`
	Uuid string `json: uuid`
	Annotations []TWAnnotation `json:annotations`
}

func (ts *TaskSet) ImportFromTaskwarrior() error {
	var tWTasks []TWTask
	// from stdin
	err := json.NewDecoder(os.Stdin).Decode(&tWTasks)

	if (err != nil) {
		return err
	}

	for _, tWTask := range tWTasks {
		fmt.Println(tWTask)
	}

	return nil
}
