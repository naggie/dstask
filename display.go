package dstask

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

/// display list of filtered tasks with context and filter
func (ts *TaskSet) DisplayByNext() {
	if ts.tasksLoaded == 0 {
		fmt.Println("\033[31mNo tasks found. Showing help.\033[0m")
		Help("")
	} else if len(ts.tasks) == 0 {
		ExitFail("No matching tasks in given context or filter.")
	} else if len(ts.tasks) == 1 {
		ts.tasks[0].Display()
		if ts.tasks[0].Notes != "" {
			fmt.Printf("\nNotes on task %d:\n\033[38;5;245m%s\033[0m\n\n", ts.tasks[0].ID, ts.tasks[0].Notes)
		}
		return
	} else {
		var tasks []*Task
		w, h := MustGetTermSize()

		h -= 8 // leave room for context message, header and prompt

		if h > len(ts.tasks) || h < 0 {
			tasks = ts.tasks
		} else {
			tasks = ts.tasks[:h]
		}

		table := NewTable(
			w,
			"ID",
			"Priority",
			"Tags",
			"Project",
			"Summary",
		)

		for _, t := range tasks {
			style := t.Style()
			table.AddRow(
				[]string{
					// id should be at least 2 chars wide to match column header
					// (headers can be truncated)
					fmt.Sprintf("%-2d", t.ID),
					t.Priority,
					strings.Join(t.Tags, " "),
					t.Project,
					t.Summary,
				},
				style,
			)
		}

		table.Render()

		if h >= len(ts.tasks) {
			fmt.Printf("\n%v tasks.\n", len(ts.tasks))
		} else {
			fmt.Printf("\n%v tasks, truncated to %v lines.\n", len(ts.tasks), h)
		}
	}
}

func (task *Task) Display() {
	w, _ := MustGetTermSize()

	table := NewTable(
		w,
		"Name",
		"Value",
	)

	table.AddRow([]string{"ID", strconv.Itoa(task.ID)}, RowStyle{})
	table.AddRow([]string{"Priority", task.Priority}, RowStyle{})
	table.AddRow([]string{"Summary", task.Summary}, RowStyle{})
	table.AddRow([]string{"Status", task.Status}, RowStyle{})
	table.AddRow([]string{"Project", task.Project}, RowStyle{})
	table.AddRow([]string{"Tags", strings.Join(task.Tags, ", ")}, RowStyle{})
	table.AddRow([]string{"UUID", task.UUID}, RowStyle{})
	table.AddRow([]string{"Created", task.Created.String()}, RowStyle{})
	if !task.Resolved.IsZero() {
		table.AddRow([]string{"Resolved", task.Resolved.String()}, RowStyle{})
	}
	if !task.Due.IsZero() {
		table.AddRow([]string{"Due", task.Due.String()}, RowStyle{})
	}
	table.Render()
}

func (t *Task) Style() RowStyle {
	now := time.Now()
	style := RowStyle{}

	if t.Status == STATUS_ACTIVE {
		style.Fg = FG_ACTIVE
		style.Bg = BG_ACTIVE
	} else if !t.Due.IsZero() && t.Due.Before(now) {
		style.Fg = FG_PRIORITY_HIGH
	} else if t.Priority == PRIORITY_CRITICAL {
		style.Fg = FG_PRIORITY_CRITICAL
	} else if t.Priority == PRIORITY_HIGH {
		style.Fg = FG_PRIORITY_HIGH
	} else if t.Priority == PRIORITY_LOW {
		style.Fg = FG_PRIORITY_LOW
	}

	if t.Status == STATUS_PAUSED {
		style.Bg = BG_PAUSED
	}

	return style
}

func (ts TaskSet) DisplayByWeek() {
	w, _ := MustGetTermSize()

	table := NewTable(
		w,
		"Resolved",
		"Priority",
		"Tags",
		"Project",
		"Summary",
		"Closing note",
	)

	var lastWeek int

	for _, t := range ts.tasks {
		_, week := t.Resolved.ISOWeek()

		if lastWeek != 0 && week != lastWeek {
			table.Render()
			// insert gap
			fmt.Printf("\n\n> Week %d, starting %s\n\n", week, t.Resolved.Format("Mon 2 Jan 2006"))
			table = NewTable(
				w,
				"Resolved",
				"Priority",
				"Tags",
				"Project",
				"Summary",
				"Closing note",
			)
		}

		noteLines := strings.Split(t.Notes, "\n")
		table.AddRow(
			[]string{
				t.Resolved.Format("Mon 2"),
				t.Priority,
				strings.Join(t.Tags, " "),
				t.Project,
				t.Summary,
				noteLines[len(noteLines)-1],
			},
			t.Style(),
		)

		_, lastWeek = t.Resolved.ISOWeek()
	}

	table.Render()
	fmt.Printf("\n%v tasks.\n", len(ts.tasks))
}

func (ts TaskSet) DisplayProjects() {
	var style RowStyle
	projects := ts.GetProjects()

	w, _ := MustGetTermSize()
	table := NewTable(
		w,
		"Created",
		"Name",
		"Progress",
	)

	for name := range projects {
		project := projects[name]
		if project.TasksNotResolved < project.TasksResolved {
			if project.Active {
				style = RowStyle{Fg: FG_ACTIVE, Bg: BG_ACTIVE}
			} else {
				style = RowStyle{}
			}

			table.AddRow(
				[]string{
					project.Created.Format("Mon 2 Jan 2006"),
					project.Name,
					fmt.Sprintf("%d/%d", project.TasksNotResolved, project.TasksResolved),
				},
				style,
			)
		}
	}

	table.Render()
}

func (ts TaskSet) DisplayCriticalTaskWarning() {
	var critical int

	for _, t := range ts.tasks {
		if (t.Priority == PRIORITY_CRITICAL) {
			critical += 1
		}
	}

	if (critical < ts.tasksLoadedCritical) {
		fmt.Printf(
			"\033[38;5;%dm%v critical tasks outside this context!\033[0m\n",
			FG_PRIORITY_CRITICAL,
			ts.tasksLoadedCritical - critical,
		)
	}
}
