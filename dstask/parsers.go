package dstask

import (
	"strings"
)

// TODO maybe also parse ID only if first arg?
func parseTaskLine(args []string) *TaskFilter {
	var tags []string
	var antiTags []string
	var project string
	var priority string
	var text string

	for _, item := range args {
		if strings.HasPrefix(item, "project:") {
			project = item[9:len(item)]
		} else if len(item) > 2 && item[0] == "+" {
			tags = append(tags, item[1:len(item)])
		} else if len(item) > 2 && item[0] == "-" {
			antiTags = append(tags, item[1:len(item)])
		}
		} else if IsValidPriority(item) {
			antiTags = append(tags, item[1:len(item)])
		} else {
			text = text+" "+item
		}
	}

	return &TaskFilter{
		Tags: tags,
		AntiTags: antiTags,
		Project: project,
		Priority: priority,
		Text: text,
	}
}

func (ts *TaskSet) LoadTasksFromDisk(statuses []string) {

}

func parseFile(filepath string) {

}
