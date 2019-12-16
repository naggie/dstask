package dstask

import (
	"fmt"
	"strconv"
	"strings"
)

type Table struct {
	Header    []string
	Rows      [][]string
	RowStyles []RowStyle
	Width     int
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
		Header: header,
		Width:  w,
		RowStyles: []RowStyle{
			RowStyle{
				Mode: MODE_HEADER,
			},
		},
	}
}

func FixStr(text string, width int) string {
	// remove after newline
	text = strings.Split(text, "\n")[0]
	if len(text) <= width {
		return fmt.Sprintf("%-"+strconv.Itoa(width)+"v", text)
	} else {
		return text[:width]
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
	widthBudget := t.Width - TABLE_COL_GAP*(len(t.Header)-1)

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

		cells := row[:]
		for i, w := range widths {
			trimmed := FixStr(cells[i], w)

			// support ' / ' markup -- show notes faded. Insert ANSI escape
			// formatting, ensuring reset to original colour for given row.
			if strings.Contains(trimmed, " "+NOTE_MODE_KEYWORD+" ") {
				trimmed = strings.Replace(
					FixStr(cells[i], w+2),
					" "+NOTE_MODE_KEYWORD+" ",
					fmt.Sprintf("\033[38;5;%dm ", FG_NOTE),
					1,
				) + fmt.Sprintf("\033[38;5;%dm", fg)
			}

			cells[i] = trimmed
		}

		line := strings.Join(cells, strings.Repeat(" ", TABLE_COL_GAP))

		// print style, line then reset
		fmt.Printf("\033[%d;38;5;%d;48;5;%dm%s\033[0m\n", mode, fg, bg, line)
	}
}
