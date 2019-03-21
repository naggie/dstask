package dstask

import (
	"fmt"
	"golang.org/x/sys/unix"
	"os"
	"strconv"
	"strings"
	"time"
)

type RowStyle struct {
	// ansi mode
	Mode int
	// xterm 256-colour palette
	Fg int
	Bg int
}

// should use a better console library after first POC

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
		table := NewTable(
			"ID",
			"Priority",
			"Tags",
			"Project",
			"Summary",
		)

		for _, t := range ts.tasks {
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

		rowsRendered := table.Render(11)

		if rowsRendered == len(ts.tasks) {
			fmt.Printf("\n%v tasks.\n", len(ts.tasks))
		} else {
			fmt.Printf("\n%v tasks, truncated to %v lines.\n", len(ts.tasks), rowsRendered)
		}
	}
}

type Table struct {
	Header       []string
	Rows         [][]string
	TermWidth    int
	TermHeight   int
	RowStyles    []RowStyle
}

// header may  havetruncated words
func NewTable(header ...string) *Table {
	ws, err := unix.IoctlGetWinsize(int(os.Stdout.Fd()), unix.TIOCGWINSZ)
	if err != nil {
		ExitFail("Not a TTY")
	}

	return &Table{
		Header:       header,
		TermWidth:    int(ws.Col),
		TermHeight:   int(ws.Row),
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

// get widths appropriate to the terminal size and TABLE_MAX_WIDTH
// cells may require padding or truncation. Cell padding of 1char between
// fields recommended -- not included.
// A nice characteristic of this, is that if there are no populated cells the
// column will disappear.
func (t *Table) calcColWidths(colGap, limit int) []int {
	targetWidth := TABLE_MAX_WIDTH

	if t.TermWidth < targetWidth {
		targetWidth = t.TermWidth
	}

	if limit < 0 {
		ExitFail("Limit does not make sense")
	}

	if limit > len(t.Rows) {
		limit = len(t.Rows)
	}

	originalWidths := make([]int, len(t.Header))

	for _, row := range t.Rows[:limit] {
		for j, cell := range row {
			if originalWidths[j] < len(cell) {
				originalWidths[j] = len(cell)
			}
		}
	}

	// initialise with original size and reduce interatively
	newWidths := originalWidths[:]

	// account for gaps of 2 chrs
	targetWidth -= colGap*len(t.Header) - 1

	for SumInts(newWidths...) > targetWidth {
		// find max col width index
		var max, maxi int

		for i, w := range newWidths {
			if w > max {
				max = w
				maxi = i
			}
		}

		// decrement, if 0 abort
		if newWidths[maxi] == 0 {
			break
		}
		newWidths[maxi] = newWidths[maxi] - 1
	}

	return newWidths
}

// theme loosely based on https://github.com/GothenburgBitFactory/taskwarrior/blob/2.6.0/doc/rc/dark-256.theme
// render table, returning count of rows rendered
// gap of zero means fit terminal exactly by truncating table -- you will want
// a larger gap to account for prompt or other text. A gap of -1 means the row
// count is not limited -- useful for reports or inspecting tasks.
func (t *Table) Render(gap int) int {
	maxRows := t.TermHeight - gap

	if maxRows < 1 {
		ExitFail("Not enough space to render anything")
	}

	widths := t.calcColWidths(2, maxRows)
	rows := append([][]string{t.Header}, t.Rows...)

	for i, row := range rows {
		cells := row[:]
		for i, w := range widths {
			cells[i] = FixStr(cells[i], w)
		}

		line := strings.Join(cells, "  ")

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

		if gap != -1 && i > maxRows {
			return i
		}
	}

	return len(t.Rows)
}

func (task *Task) Display() {
	table := NewTable(
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
	table.Render(0)
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
	table := NewTable(
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
			table.Render(-1)
			// insert gap
			fmt.Printf("\n\n> Week %d, starting %s\n\n", week, t.Resolved.Format("Mon 2 Jan 2006"))
			table = NewTable(
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

	table.Render(-1)
	fmt.Printf("\n%v tasks.\n", len(ts.tasks))
}
