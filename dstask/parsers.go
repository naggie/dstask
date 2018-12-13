package dstask

import (
	"strconv"
	"strings"
)

func parseTaskLine(args []string) *TaskLine {
	var id uint64
	var tags []string
	var antiTags []string
	var project string
	var priority string
	var text string

	for i, item := range args {
		if s, err := strconv.ParseUint(item, 10, 64); i == 0 && err == nil {
			id = s
		} else if strings.HasPrefix(item, "project:") {
			project = item[9:len(item)]
		} else if len(item) > 2 && item[0:1] == "+" {
			tags = append(tags, item[1:len(item)])
		} else if len(item) > 2 && item[0:1] == "-" {
			antiTags = append(tags, item[1:len(item)])
		} else if IsValidPriority(item) {
			antiTags = append(tags, item[1:len(item)])
		} else {
			text = text + " " + item
		}
	}

	return &TaskLine{
		Id:       id,
		Tags:     tags,
		AntiTags: antiTags,
		Project:  project,
		Priority: priority,
		Text:     text,
	}
}

func (ts *TaskSet) LoadTasksFromDisk(statuses []string) {

}

func parseFile(filepath string) {

}
