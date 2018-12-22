package dstask

const (
	STATUS_PENDING   = "pending"
	STATUS_ACTIVE    = "active"
	STATUS_RESOLVED  = "resolved"
	STATUS_DELEGATED = "delegated"
	STATUS_DEFERRED  = "deferred"
	STATUS_SOMEDAY   = "someday"
	STATUS_RECURRING = "recurring" // tentative

	CMD_ADD          = "add"
	CMD_START        = "start"
	CMD_ANNOTATE     = "annotate"
	CMD_STOP         = "stop"
	CMD_RESOLVE      = "resolve"
	CMD_CONTEXT      = "context"
	CMD_MODIFY       = "modify"
	CMD_EDIT         = "edit"
	CMD_UNDO         = "undo"
	CMD_GIT          = "git"
	CMD_DAY          = "day"
	CMD_WEEK         = "week"
	CMD_PROJECTS     = "projects"
	CMD_IMPORT_TW    = "import-tw"
	CMD_HELP         = "help"


	// filter: P1 P2 etc
	PRIORITY_CRITICAL = "P0"
	PRIORITY_HIGH     = "P1"
	PRIORITY_NORMAL   = "P2"
	PRIORITY_LOW      = "P3"

	MAX_TASKS_OPEN = 10000
)

// for import (etc) it's necessary to have full context
var ALL_STATUSES = []string{
	STATUS_ACTIVE,
	STATUS_PENDING,
	STATUS_DELEGATED,
	STATUS_DEFERRED,
	STATUS_SOMEDAY,
	STATUS_RECURRING,
	STATUS_RESOLVED,
}

// for most operations, it's not necessary or desirable to load the expensive resolved tasks
var NON_RESOLVED_STATUSES = []string{
	STATUS_ACTIVE,
	STATUS_PENDING,
	STATUS_DELEGATED,
	STATUS_DEFERRED,
	STATUS_SOMEDAY,
	STATUS_RECURRING,
}

var STATUS_ORDER = map[string]int{
	STATUS_ACTIVE:    1,
	STATUS_PENDING:   2,
	STATUS_DELEGATED: 3,
	STATUS_DEFERRED:  4,
	STATUS_SOMEDAY:   5,
	STATUS_RECURRING: 6,
	STATUS_RESOLVED:  7,
}

var ALL_CMDS = []string{
	CMD_ADD,
	CMD_START,
	CMD_ANNOTATE,
	CMD_STOP,
	CMD_RESOLVE,
	CMD_CONTEXT,
	CMD_MODIFY,
	CMD_EDIT,
	CMD_UNDO,
	CMD_GIT,
	CMD_DAY,
	CMD_WEEK,
	CMD_PROJECTS,
	CMD_IMPORT_TW,
	CMD_HELP,
}
