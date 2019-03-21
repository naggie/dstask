package dstask

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Table struct {
	Header       []string
	Rows         [][]string
	RowStyles    []RowStyle
	Width        int
}

type RowStyle struct {
	// ansi mode
	Mode int
	// xterm 256-colour palette
	Fg int
	Bg int
}

// header may  havetruncated words
func NewTable(w int, header ...string) *Table {
	if w > TABLE_MAX_WIDTH {
		w = TABLE_MAX_WIDTH
	}

	return &Table{
		Header:       header,
		Width:        w,
		RowStyles: []RowStyle{
			RowStyle{
				Mode: MODE_HEADER,
			},
		},
	}
}

func (t *Table) AddRow(row []string, style RowStyle) {
	if len(row) != len(t.Header) {
		panic("Row is incorrect length")
	}

	t.Rows = append(t.Rows, row)
	t.RowStyles = append(t.RowStyles, style)
}

// render table, returning count of rows rendered
// gap of zero means fit terminal exactly by truncating table -- you will want
// a larger gap to account for prompt or other text. A gap of -1 means the row
// count is not limited -- useful for reports or inspecting tasks.
func (t *Table) Render() {
	originalWidths := make([]int, len(t.Header))

	for _, row := range t.Rows {
		for j, cell := range row {
			if originalWidths[j] < len(cell) {
				originalWidths[j] = len(cell)
			}
		}
	}

	// initialise with original size and reduce interatively
	widths := originalWidths[:]

	// account for gaps of 2 chrs
	widthBudget := t.Width - TABLE_COL_GAP*(len(t.Header) - 1)

	for SumInts(widths...) > widthBudget {
		// find max col width index
		var max, maxi int

		for i, w := range widths {
			if w > max {
				max = w
				maxi = i
			}
		}

		// decrement, if 0 abort
		if widths[maxi] == 0 {
			break
		}
		widths[maxi] = widths[maxi] - 1
	}

	rows := append([][]string{t.Header}, t.Rows...)

	for i, row := range rows {
		cells := row[:]
		for i, w := range widths {
			cells[i] = FixStr(cells[i], w)
		}

		line := strings.Join(cells, strings.Repeat(" ", TABLE_COL_GAP))

		mode := t.RowStyles[i].Mode
		fg := t.RowStyles[i].Fg
		bg := t.RowStyles[i].Bg

		// defaults
		if mode == 0 {
			mode = MODE_DEFAULT
		}

		if fg == 0 {
			fg = FG_DEFAULT
		}

		if bg == 0 {
			/// alternate if not specified
			if i%2 != 0 {
				bg = BG_DEFAULT_1
			} else {
				bg = BG_DEFAULT_2
			}
		}

		// print style, line then reset
		fmt.Printf("\033[%d;38;5;%d;48;5;%dm%s\033[0m\n", mode, fg, bg, line)
	}
}

/// display list of filtered tasks with context and filter
func (ts *TaskSet) DisplayByNext() {
	if ts.numTasksLoaded == 0 {
		fmt.Println("\033[31mNo tasks found. Showing help.\033[0m")
		Help("")
	} else if len(ts.tasks) == 0 {
		ExitFail("No matching tasks in given context or filter.")
	} else if len(ts.tasks) == 1 {
		ts.tasks[0].Display()
		return
	} else {
		var tasks []*Task
		w, h := MustGetTermSize()

		h -=8

		if h < 4 {
			ExitFail("Terminal is too small to display next tasks")
		}

		if h > len(ts.tasks) {
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

		if h == len(ts.tasks) {
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
	table.AddRow([]string{"Notes", task.Notes}, RowStyle{})
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

func (ts TaskSet) DisplayByResolved() {
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

		noteLines := strings.Split(t.Notes,"\n")
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
