package dstask

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

/// display list of filtered tasks with context and filter
func (ts *TaskSet) DisplayByNext(truncate bool) {
	tasks := ts.Tasks()
	total := len(tasks)
	if ts.NumTotal() == 0 {
		fmt.Println("\033[31mNo tasks found. Showing help.\033[0m")
		Help("")
	} else if len(tasks) == 0 {
		ExitFail("No matching tasks in given context or filter.")
	} else if len(tasks) == 1 {
		task := tasks[0]
		task.Display()
		if task.Notes != "" {
			fmt.Printf("\nNotes on task %d:\n\033[38;5;245m%s\033[0m\n\n", task.ID, task.Notes)
		}
		return
	} else {
		w, h := MustGetTermSize()

		maxTasks := h - TERMINAL_HEIGHT_MARGIN // leave room for context message, header and prompt

		if maxTasks < MIN_TASKS_SHOWN {
			maxTasks = MIN_TASKS_SHOWN
		}

		if truncate && maxTasks < len(tasks) {
			tasks = tasks[:maxTasks]
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
					t.LongSummary(),
				},
				style,
			)
		}

		table.Render()

		if truncate && maxTasks < total {
			fmt.Printf("\n%v tasks, truncated to %v lines.\n", total, h)
		} else {
			fmt.Printf("\n%v tasks.\n", total)
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

// TODO combine with previous with interface, plus computed Project status
func (p *Project) Style() RowStyle {
	style := RowStyle{}

	if p.Active {
		style.Fg = FG_ACTIVE
		style.Bg = BG_ACTIVE
	} else if p.Priority == PRIORITY_CRITICAL {
		style.Fg = FG_PRIORITY_CRITICAL
	} else if p.Priority == PRIORITY_HIGH {
		style.Fg = FG_PRIORITY_HIGH
	} else if p.Priority == PRIORITY_LOW {
		style.Fg = FG_PRIORITY_LOW
	}

	return style
}

func (ts TaskSet) DisplayByWeek() {
	w, _ := MustGetTermSize()
	var table *Table
	var lastWeek int
	tasks := ts.Tasks()

	for _, t := range tasks {
		_, week := t.Resolved.ISOWeek()

		// guaranteed true for first iteration, ISOweek starts with 1.
		if week != lastWeek {
			if table != nil && len(table.Rows) > 0 {
				table.Render()
			}
			// insert gap
			fmt.Printf("\n\n> Week %d, starting %s\n\n", week, t.Resolved.Format("Mon 2 Jan 2006"))
			table = NewTable(
				w,
				"Resolved",
				"Priority",
				"Tags",
				"Project",
				"Summary",
			)
		}

		table.AddRow(
			[]string{
				t.Resolved.Format("Mon 2"),
				t.Priority,
				strings.Join(t.Tags, " "),
				t.Project,
				t.LongSummary(),
			},
			t.Style(),
		)

		_, lastWeek = t.Resolved.ISOWeek()
	}

	if table != nil {
		table.Render()
	}
	fmt.Printf("%v tasks.\n", len(tasks))
}

func (ts TaskSet) DisplayProjects() {
	projects := ts.GetProjects()
	w, _ := MustGetTermSize()
	table := NewTable(
		w,
		"Name",
		"Progress",
		"Created",
	)

	for name := range projects {
		project := projects[name]
		if project.TasksResolved < project.Tasks {

			table.AddRow(
				[]string{
					project.Name,
					fmt.Sprintf("%d/%d", project.TasksResolved, project.Tasks),
					project.Created.Format("Mon 2 Jan 2006"),
				},
				project.Style(),
			)
		}
	}

	table.Render()
}

func (ts TaskSet) DisplayCriticalTaskWarning() {
	var critical int
	var totalCritical int

	for _, t := range ts.Tasks() {
		if t.Priority == PRIORITY_CRITICAL {
			critical += 1
		}
	}

	for _, t := range ts.AllTasks() {
		if t.Priority == PRIORITY_CRITICAL {
			totalCritical += 1
		}
	}

	if critical < totalCritical {
		fmt.Printf(
			"\033[38;5;%dm%v critical tasks outside this context! Use `dstask -- P0` to see them.\033[0m\n",
			FG_PRIORITY_CRITICAL,
			totalCritical-critical,
		)
	}
}
