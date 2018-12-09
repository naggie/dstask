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

// see https://taskwarrior.org/docs/design/task.html
type TwTask struct {
	Description string `json:"description"`
	End string `json:"end"`
	Entry string `json: entry`
	Start string `json: start`
	Modified string `json: modified`
	Status string `json: status`
	Tags []string `json: tags`
	Uuid string `json: uuid`
	Annotations []TWAnnotation `json:annotations`
}

// convert a tw status into a dstask status
func convertStatus(twStatus string, start string) string {
	if start != "" {
		return STATUS_ACTIVE
	}

	switch twStatus {
		case "completed":
			return STATUS_RESOLVED
		case "deleted":
			return STATUS_RESOLVED
		case "waiting":
			return STATUS_PENDING
		case "recurring":
			// TODO -- implement reccurence
			//return STATUS_RECURRING
			return STATUS_RESOLVED
		default:
			return twStatus
	}
}

func (ts *TaskSet) ImportFromTaskwarrior() error {
	var twTasks []TwTask
	// from stdin
	err := json.NewDecoder(os.Stdin).Decode(&twTasks)

	if (err != nil) {
		return err
	}

	for _, twTask := range twTasks {
		fmt.Println(twTask)
		ts.tasks = append(ts.tasks, Task{
			uuid: twTask.Uuid,
			status: convertStatus(twTask.Status, twTask.Start),
		})
	}

	return nil
}
