package completions

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/naggie/dstask"
)

// Completions ...
func Completions(conf dstask.Config, args []string, ctx dstask.Query) {
	// given the entire user's command line arguments as the arguments for
	// this cmd, suggest possible candidates for the last arg.
	// see the relevant shell completion bindings in this repository for
	// integration. Note there are various idiosyncrasies with bash
	// involving arg separation.
	var completions []string
	var originalArgs []string
	var prefix string

	// drop dstask _completions dstask to allow parsing what is on actual
	// prompt
	if len(args) > 3 {
		originalArgs = args[3:]
	}

	// args are dstask _completions <user command line>
	// parse command line as normal to set rules
	query := dstask.ParseQuery(originalArgs...)

	// No command and OK to specify command (to run or help)
	// Note that techically we should only specify commands as available
	// completions if the last partial argument is a command substring.
	// However, this is unnecessary as a general substring filter is used at
	// the end of the func.
	// This is exhaustive but the clearest way, IMO.
	if len(query.AntiProjects) == 0 &&
		query.Project == "" &&
		len(query.Tags) == 0 &&
		len(query.AntiTags) == 0 &&
		query.Priority == "" &&
		query.Template == 0 &&
		!query.IgnoreContext &&
		(query.Cmd == dstask.CMD_HELP || query.Cmd == "") {
		for _, cmd := range dstask.ALL_CMDS {
			if !strings.HasPrefix(cmd, "_") {
				completions = append(completions, cmd)
			}
		}
	}

	if dstask.StrSliceContains([]string{
		"",
		dstask.CMD_NEXT,
		dstask.CMD_ADD,
		dstask.CMD_REMOVE,
		dstask.CMD_LOG,
		dstask.CMD_START,
		dstask.CMD_STOP,
		dstask.CMD_DONE,
		dstask.CMD_RESOLVE,
		dstask.CMD_CONTEXT,
		dstask.CMD_MODIFY,
		dstask.CMD_SHOW_NEXT,
		dstask.CMD_SHOW_PROJECTS,
		dstask.CMD_SHOW_ACTIVE,
		dstask.CMD_SHOW_PAUSED,
		dstask.CMD_SHOW_OPEN,
		dstask.CMD_SHOW_RESOLVED,
		dstask.CMD_SHOW_TEMPLATES,
	}, query.Cmd) {
		ts, err := dstask.LoadTaskSet(conf.Repo, conf.IDsFile, false)
		if err != nil {
			log.Printf("completions error: %v\n", err)
			return

		}
		// limit completions to available context, but not if the user is
		// trying to change context, context ignore is on, or modify
		// command is being completed
		if !query.IgnoreContext &&
			query.Cmd != dstask.CMD_CONTEXT &&
			query.Cmd != dstask.CMD_MODIFY {
			ts.Filter(ctx)
		}

		// templates
		if query.Cmd == dstask.CMD_ADD {
			for _, task := range ts.Tasks() {
				if task.Status == dstask.STATUS_TEMPLATE {
					completions = append(completions, "template:"+strconv.Itoa(task.ID))
				}
			}
		}

		// priorities
		completions = append(completions, dstask.PRIORITY_CRITICAL)
		completions = append(completions, dstask.PRIORITY_HIGH)
		completions = append(completions, dstask.PRIORITY_NORMAL)
		completions = append(completions, dstask.PRIORITY_LOW)

		// projects
		for _, project := range ts.GetProjects() {
			completions = append(completions, "project:"+project.Name)
			completions = append(completions, "-project:"+project.Name)
		}

		// tags
		for tag := range ts.GetTags() {
			completions = append(completions, "+"+tag)
			completions = append(completions, "-"+tag)
		}
	}

	if len(originalArgs) > 0 {
		prefix = originalArgs[len(originalArgs)-1]
	}

	for _, completion := range completions {
		if strings.HasPrefix(completion, prefix) && !dstask.StrSliceContains(originalArgs, completion) {
			fmt.Println(completion)
		}
	}
}
