package dstask

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type TwAnnotation struct {
	Description string `json:"description"`
	Entry       string `json: entry`
}

// see https://taskwarrior.org/docs/design/task.html
type TwTask struct {
	Description string         `json:"description"`
	End         string         `json:"end"`
	Entry       string         `json: entry`
	Start       string         `json: start`
	Modified    string         `json: modified`
	Status      string         `json: status`
	Project     string         `json: project`
	Priority    string         `json: priority`
	Depends     string         `json: depends`
	Tags        []string       `json: tags`
	Uuid        string         `json: uuid`
	Annotations []TwAnnotation `json:annotations`
}

var priorityMap = map[string]string{
	"H": PRIORITY_HIGH,
	"M": PRIORITY_NORMAL,
	"L": PRIORITY_LOW,
}

// convert a tw status into a dstask status
func (t *TwTask) ConvertStatus() string {
	if t.Start != "" {
		return STATUS_ACTIVE
	}

	switch t.Status {
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
		return t.Status
	}
}

func (t *TwTask) ConvertAnnotations() []string {
	var comments []string

	for _, ann := range t.Annotations {
		comments = append(comments, ann.Description)
	}

	return comments
}

func (ts *TaskSet) ImportFromTaskwarrior() error {
	var twTasks []TwTask
	// from stdin
	err := json.NewDecoder(os.Stdin).Decode(&twTasks)

	if err != nil {
		return err
	}

	for _, twTask := range twTasks {
		fmt.Println(twTask)
		ts.tasks = append(ts.tasks, Task{
			uuid:         twTask.Uuid,
			status:       twTask.ConvertStatus(),
			Summary:      twTask.Description,
			Tags:         twTask.Tags,
			Project:      twTask.Project,
			Priority:     priorityMap[twTask.Priority],
			Comments:     twTask.ConvertAnnotations(),
			Dependencies: strings.Split(twTask.Depends, ","),
			//Created
			//Modified
			//Resolved
			//Due
		})
	}

	return nil
}
