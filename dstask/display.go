package dstask

import (
	"fmt"
	"golang.org/x/sys/unix"
	"os"
)

const (
	// keep it readable
	TABLE_MAX_WIDTH = 160
)

// should use a better console library after first POC

/// display list of filtered tasks with context and filter
func (ts *TaskSet) Display() {
	for _, t := range ts.Tasks {
		fmt.Printf("%+v\n", t)
	}
}

// display a single task in detail, with numbered subtasks
func (t *Task) Display() {

}

type Table struct {
	Header []string
	Rows [][]string
	MaxColWidths []int
	TermWidth int
	TermHeight int
}

func NewTable(header []string) *Table {
	ws, err := unix.IoctlGetWinsize(int(os.Stdout.Fd()), unix.TIOCGWINSZ)
	if err != nil {
		ExitFail("Not a TTY")
	}

	return &Table{
		Header: header,
		MaxColWidths: make([]int, len(header)),
		TermWidth: int(ws.Col),
		TermHeight: int(ws.Row),
	}
}

func (t *Table) AddRow(row []string) {
	if len(row) != len(t.Header) {
		panic("Row is incorrect length")
	}

	for i, cell := range(row) {
		if t.MaxColWidths[i] < len(cell) {
			t.MaxColWidths[i] = len(cell)
		}
	}

	t.Rows = append(t.Rows,row)
}

// get widths appropriate to the terminal size and TABLE_MAX_WIDTH
// cells may require padding or truncation. Cell padding of 1char between
// fields recommended -- not included.
func (t *Table) calcColWidths() []int {
	target := TABLE_MAX_WIDTH

	if t.TermWidth < target {
		target = t.TermWidth
	}

	colWidths := t.MaxColWidths[:]

	for SumInts(colWidths...) > target {
		// find max col width index
		var max, maxi int

		for i,w := range(colWidths) {
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

func (t *Table) Render() {
	// TODO: ansi colours

}
