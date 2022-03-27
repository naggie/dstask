package dstask

import "os"

func init() {
	if os.Getenv("DSTASK_FAKE_PTY") != "" {
		FAKE_PTY = true
	}
}

var (
	// for CI testing
	FAKE_PTY = false
	// populated by linker flags, see do-release.sh
	GIT_COMMIT = "Unknown"
	VERSION    = "Unknown"
	BUILD_DATE = "Unknown"
)

const (
	STATUS_PENDING   = "pending"
	STATUS_ACTIVE    = "active"
	STATUS_RESOLVED  = "resolved"
	STATUS_DELEGATED = "delegated"
	STATUS_DEFERRED  = "deferred"
	STATUS_PAUSED    = "paused"
	STATUS_RECURRING = "recurring"
	STATUS_TEMPLATE  = "template"

	CMD_NEXT             = "next"
	CMD_ADD              = "add"
	CMD_RM               = "rm"
	CMD_REMOVE           = "remove"
	CMD_TEMPLATE         = "template"
	CMD_LOG              = "log"
	CMD_START            = "start"
	CMD_NOTE             = "note"
	CMD_NOTES            = "notes"
	CMD_STOP             = "stop"
	CMD_DONE             = "done"
	CMD_RESOLVE          = "resolve"
	CMD_CONTEXT          = "context"
	CMD_MODIFY           = "modify"
	CMD_EDIT             = "edit"
	CMD_UNDO             = "undo"
	CMD_SYNC             = "sync"
	CMD_OPEN             = "open"
	CMD_GIT              = "git"
	CMD_SHOW_NEXT        = "show-next"
	CMD_SHOW_PROJECTS    = "show-projects"
	CMD_SHOW_TAGS        = "show-tags"
	CMD_SHOW_ACTIVE      = "show-active"
	CMD_SHOW_PAUSED      = "show-paused"
	CMD_SHOW_OPEN        = "show-open"
	CMD_SHOW_RESOLVED    = "show-resolved"
	CMD_SHOW_TEMPLATES   = "show-templates"
	CMD_SHOW_UNORGANISED = "show-unorganised"
	CMD_COMPLETIONS      = "_completions"
	CMD_HELP             = "help"
	CMD_VERSION          = "version"

	CMD_PRINT_ZSH_COMPLETION = "zsh-completion"
	CMD_PRINT_BASH_COMPLETION = "bash-completion"

	// filter: P1 P2 etc
	PRIORITY_CRITICAL = "P0"
	PRIORITY_HIGH     = "P1"
	PRIORITY_NORMAL   = "P2"
	PRIORITY_LOW      = "P3"

	MAX_TASKS_OPEN    = 10000
	TASK_FILENAME_LEN = 40

	// if the terminal is too short, show this many tasks anyway
	MIN_TASKS_SHOWN = 8

	// reserve this many lines for status messages/prompt
	TERMINAL_HEIGHT_MARGIN = 9

	IGNORE_CONTEXT_KEYWORD = "--"
	NOTE_MODE_KEYWORD      = "/"

	// theme loosely based on https://github.com/GothenburgBitFactory/taskwarrior/blob/2.6.0/doc/rc/dark-256.theme
	TABLE_MAX_WIDTH      = 160 // keep it readable
	TABLE_COL_GAP        = 2   // differentiate columns
	MODE_HEADER          = 4
	FG_DEFAULT           = 250
	BG_DEFAULT_1         = 233
	BG_DEFAULT_2         = 232
	MODE_DEFAULT         = 0
	FG_ACTIVE            = 233
	BG_ACTIVE            = 250
	BG_PAUSED            = 236 // task that has been started then stopped
	FG_PRIORITY_CRITICAL = 160
	FG_PRIORITY_HIGH     = 166
	FG_PRIORITY_NORMAL   = FG_DEFAULT
	FG_PRIORITY_LOW      = 245
	FG_NOTE              = 240
)

// for import (etc) it's necessary to have full context
var ALL_STATUSES = []string{
	STATUS_ACTIVE,
	STATUS_PENDING,
	STATUS_DELEGATED,
	STATUS_DEFERRED,
	STATUS_PAUSED,
	STATUS_RECURRING,
	STATUS_RESOLVED,
	STATUS_TEMPLATE,
}

// statuses which are hidden by default (direct addressing or show- commands
// needed to see them)
var HIDDEN_STATUSES = []string{
	STATUS_RECURRING,
	STATUS_RESOLVED,
	STATUS_TEMPLATE,
}

// incomplete until all statuses are implemented
var VALID_STATUS_TRANSITIONS = [][]string{
	{STATUS_PENDING, STATUS_ACTIVE},
	{STATUS_ACTIVE, STATUS_PAUSED},
	{STATUS_PAUSED, STATUS_ACTIVE},
	{STATUS_PENDING, STATUS_RESOLVED},
	{STATUS_PAUSED, STATUS_RESOLVED},
	{STATUS_ACTIVE, STATUS_RESOLVED},
	{STATUS_PENDING, STATUS_TEMPLATE},
}

// for most operations, it's not necessary or desirable to load the expensive resolved tasks
var NON_RESOLVED_STATUSES = []string{
	STATUS_ACTIVE,
	STATUS_PENDING,
	STATUS_DELEGATED,
	STATUS_DEFERRED,
	STATUS_PAUSED,
	STATUS_RECURRING,
	STATUS_TEMPLATE,
}

var ALL_CMDS = []string{
	CMD_NEXT,
	CMD_ADD,
	CMD_RM,
	CMD_REMOVE,
	CMD_TEMPLATE,
	CMD_LOG,
	CMD_START,
	CMD_NOTE,
	CMD_NOTES,
	CMD_STOP,
	CMD_DONE,
	CMD_RESOLVE,
	CMD_CONTEXT,
	CMD_MODIFY,
	CMD_EDIT,
	CMD_UNDO,
	CMD_SYNC,
	CMD_OPEN,
	CMD_GIT,
	CMD_SHOW_NEXT,
	CMD_SHOW_PROJECTS,
	CMD_SHOW_TAGS,
	CMD_SHOW_ACTIVE,
	CMD_SHOW_PAUSED,
	CMD_SHOW_OPEN,
	CMD_SHOW_RESOLVED,
	CMD_SHOW_TEMPLATES,
	CMD_SHOW_UNORGANISED,
	CMD_COMPLETIONS,
	CMD_PRINT_BASH_COMPLETION,
	CMD_PRINT_ZSH_COMPLETION,
	CMD_HELP,
	CMD_VERSION,
}
