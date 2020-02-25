package dstask

// import tasks from taskwarrior

// see https://taskwarrior.org/docs/design/task.html

// usage: task export | dstask import
// Filters can be used in taskwarrior to specify a subset.

import (
	"encoding/json"
	"os"
	"strings"
	"time"
)

type TwTime struct {
	time.Time
}

// TW time is stored in ISO 8601 basic format, which the default parser does
// not understand.
func (tt *TwTime) UnmarshalJSON(b []byte) error {
	s := string(b)

	if s == "null" {
		return nil
	}

	// drop quotes
	if len(s) > 2 {
		s = s[1 : len(s)-1]
	}

	if len(s) == 16 {
		// convert from basic format to normal format which is RFC3339 compatible
		s = s[0:4] + "-" + s[4:6] + "-" + s[6:11] + ":" + s[11:13] + ":" + s[13:]
	}

	t, err := time.Parse(time.RFC3339, s)
	tt.Time = t

	if err != nil {
		return err
	}

	return nil
}

type TwAnnotation struct {
	Description string
	Entry       string
}

type TwTask struct {
	Description string
	End         TwTime
	Entry       TwTime
	Start       TwTime
	Modified    TwTime
	Due         TwTime
	Status      string
	Project     string
	Priority    string
	Depends     string
	Tags        []string
	UUID        string
	Annotations []TwAnnotation
}

var priorityMap = map[string]string{
	"H": PRIORITY_HIGH,
	"M": PRIORITY_NORMAL,
	"L": PRIORITY_LOW,
	"":  PRIORITY_NORMAL,
}

func (t *TwTask) ConvertAnnotations() string {
	var comments []string

	for _, ann := range t.Annotations {
		comments = append(comments, ann.Description)
	}

	return strings.Join(comments, "\n")
}

// convert a tw status into a dstask status
func (t *TwTask) ConvertStatus() string {
	if !t.Start.Time.IsZero() {
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

// resolved time is not tracked. Give best guess.
func (t *TwTask) GetResolvedTime() time.Time {
	if t.Status == "completed" {
		return t.Modified.Time
	} else {
		return time.Time{}
	}
}

func (ts *TaskSet) ImportFromTaskwarrior() error {
	var twtasks []TwTask
	// from stdin
	err := json.NewDecoder(os.Stdin).Decode(&twtasks)

	if err != nil {
		ExitFail("Failed to decode JSON from stdin")
	}

	for _, twTask := range twtasks {
		ts.LoadTask(Task{
			UUID:         twTask.UUID,
			Status:       twTask.ConvertStatus(),
			WritePending: true,
			Summary:      twTask.Description,
			Tags:         twTask.Tags,
			Project:      twTask.Project,
			Priority:     priorityMap[twTask.Priority],
			Notes:        twTask.ConvertAnnotations(),
			// FieldsFunc required instead of split as split returns a slice of len(1) when empty...
			Dependencies: strings.FieldsFunc(twTask.Depends, func(c rune) bool { return c == ',' }),
			Created:      twTask.Entry.Time,
			Resolved:     twTask.GetResolvedTime(),
			Due:          twTask.Due.Time,
		})
	}

	return nil
}
