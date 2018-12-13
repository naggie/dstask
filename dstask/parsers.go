package dstask

import (
	"strconv"
	"strings"
)

func ParseTaskLine(args []string) *TaskLine {
	var id int
	var tags []string
	var antiTags []string
	var project string
	var priority string
	var words []string

	for i, item := range args {
		if s, err := strconv.ParseInt(item, 10, 64); i == 0 && err == nil {
			id = int(s)
		} else if strings.HasPrefix(item, "project:") {
			project = item[8:len(item)]
		} else if len(item) > 2 && item[0:1] == "+" {
			tags = append(tags, item[1:len(item)])
		} else if len(item) > 2 && item[0:1] == "-" {
			antiTags = append(antiTags, item[1:len(item)])
		} else if IsValidPriority(item) {
			priority = item
		} else {
			words = append(words, item)
		}
	}

	return &TaskLine{
		Id:       id,
		Tags:     tags,
		AntiTags: antiTags,
		Project:  project,
		Priority: priority,
		Text:     strings.Join(words, " "),
	}
}
