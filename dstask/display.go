package dstask

import (
	"fmt"
	"golang.org/x/sys/unix"
	"os"
	"strings"
)

const (
	// keep it readable
	TABLE_MAX_WIDTH = 160

	// styles for rows
	STYLE_HEADER = iota
	STYLE_ACTIVE
	STYLE_PRIORITY_CRITICAL
	STYLE_PRIORITY_HIGH
	STYLE_PRIORITY_NORMAL
	STYLE_PRIORITY_LOW
)

// should use a better console library after first POC

/// display list of filtered tasks with context and filter
func (ts *TaskSet) Display() {
	table := NewTable(
		"ID",
		"Priority",
		"Tags",
		"Project",
		"Summary",
	)

	for _, t := range ts.Tasks {
		style := STYLE_PRIORITY_NORMAL

		// TODO important if overdue
		if t.status == STATUS_ACTIVE {
			style = STYLE_ACTIVE
		} else if t.Priority == PRIORITY_CRITICAL {
			style = STYLE_PRIORITY_CRITICAL
		} else if t.Priority == PRIORITY_HIGH {
			style = STYLE_PRIORITY_HIGH
		} else if t.Priority == PRIORITY_LOW {
			style = STYLE_PRIORITY_LOW
		}

		table.AddRow(
			[]string{
				// id should be at least 2 chars wide to match column header
				// (headers can be truncated)
				fmt.Sprintf("%-2d", t.id),
				t.Priority,
				strings.Join(t.Tags, " "),
				t.Project,
				t.Summary,
			},
			style,
		)
	}

	// push off prompt
	fmt.Printf("\n\n")

	// TODO print current context here

	rowsRendered := table.Render(10)

	if rowsRendered == len(ts.Tasks) {
		fmt.Printf("\n%v tasks.\n", len(ts.Tasks))
	} else {
		fmt.Printf("\n%v tasks, truncated to %v lines.\n", len(ts.Tasks), rowsRendered)
	}
}

// display a single task in detail, with numbered subtasks
func (t *Task) Display() {

}

type Table struct {
	Header       []string
	Rows         [][]string
	MaxColWidths []int
	TermWidth    int
	TermHeight   int
	RowStyles    []int
}

// header may  havetruncated words
func NewTable(header ...string) *Table {
	ws, err := unix.IoctlGetWinsize(int(os.Stdout.Fd()), unix.TIOCGWINSZ)
	if err != nil {
		ExitFail("Not a TTY")
	}

	return &Table{
		Header:       header,
		MaxColWidths: make([]int, len(header)),
		TermWidth:    int(ws.Col),
		TermHeight:   int(ws.Row),
		RowStyles:    []int{STYLE_HEADER},
	}
}

func (t *Table) AddRow(row []string, style int) {
	if len(row) != len(t.Header) {
		panic("Row is incorrect length")
	}

	for i, cell := range row {
		if t.MaxColWidths[i] < len(cell) {
			t.MaxColWidths[i] = len(cell)
		}
	}

	t.Rows = append(t.Rows, row)
	t.RowStyles = append(t.RowStyles, style)
}

// get widths appropriate to the terminal size and TABLE_MAX_WIDTH
// cells may require padding or truncation. Cell padding of 1char between
// fields recommended -- not included.
func (t *Table) calcColWidths(gap int) []int {
	target := TABLE_MAX_WIDTH

	if t.TermWidth < target {
		target = t.TermWidth
	}

	colWidths := t.MaxColWidths[:]

	// account for gaps
	target -= gap*len(colWidths) - 1

	for SumInts(colWidths...) > target {
		// find max col width index
		var max, maxi int

		for i, w := range colWidths {
			if w > max {
				max = w
				maxi = i
			}
		}

		// decrement, if 0 abort
		if colWidths[maxi] == 0 {
			break
		}
		colWidths[maxi] = colWidths[maxi] - 1
	}

	return colWidths
}

// theme loosely based on https://github.com/GothenburgBitFactory/taskwarrior/blob/2.6.0/doc/rc/dark-256.theme
// render table, returning count of rows rendered
func (t *Table) Render(gap int) int {
	// TODO highlight overdue, high priority, low priority in progress
	// TODO alternate row colours (tw)
	// TODO see screenshot for reference https://taskwarrior.org/docs/themes.html#default

	var fg, bg, mode int

	widths := t.calcColWidths(2)
	maxRows := t.TermHeight - gap
	rows := append([][]string{t.Header}, t.Rows...)

	for i, row := range rows {
		cells := row[:]
		for i, w := range widths {
			cells[i] = FixStr(cells[i], w)
		}

		line := strings.Join(cells, "  ")

		// default
		fg = 250
		bg = 232
		mode = 0

		// default row style alternates. FG and BG can be overridden
		// independently.
		if i%2 != 0 {
			bg = 233
		} else {
			bg = 232
		}

		switch t.RowStyles[i] {
		case STYLE_HEADER:
			// header -- underline
			mode = 4
		case STYLE_ACTIVE:
			fg = 255
			bg = 166
		case STYLE_PRIORITY_CRITICAL:
			fg = 160
		case STYLE_PRIORITY_HIGH:
			fg = 166
		case STYLE_PRIORITY_LOW:
			fg = 245
		}

		// print style, line then reset
		//fmt.Printf("\033[%dm\033[38;5;%dm\033[48;5;%dm%s\033[0m\n", mode, fg, bg, line)
		fmt.Printf("\033[%d;38;5;%d;48;5;%dm%s\033[0m\n", mode, fg, bg, line)

		if i > maxRows {
			return i
		}
	}

	return len(t.Rows)
}
