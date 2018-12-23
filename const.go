package dstask

const (
	STATUS_PENDING   = "pending"
	STATUS_ACTIVE    = "active"
	STATUS_RESOLVED  = "resolved"
	STATUS_DELEGATED = "delegated"
	STATUS_DEFERRED  = "deferred"
	STATUS_SOMEDAY   = "someday"
	STATUS_RECURRING = "recurring" // tentative

	CMD_NEXT      = "next"
	CMD_ADD       = "add"
	CMD_START     = "start"
	CMD_ANNOTATE  = "annotate"
	CMD_STOP      = "stop"
	CMD_RESOLVE   = "resolve"
	CMD_CONTEXT   = "context"
	CMD_MODIFY    = "modify"
	CMD_EDIT      = "edit"
	CMD_UNDO      = "undo"
	CMD_PUSH      = "push"
	CMD_PULL      = "pull"
	CMD_GIT       = "git"
	CMD_DAY       = "day"
	CMD_WEEK      = "week"
	CMD_PROJECTS  = "projects"
	CMD_IMPORT_TW = "import-tw"
	CMD_HELP      = "help"

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

// incomplete until all statuses are implemented
var VALID_STATUS_TRANSITIONS = [][]string{
	[]string{STATUS_PENDING, STATUS_ACTIVE},
	[]string{STATUS_ACTIVE, STATUS_PENDING},
	[]string{STATUS_ACTIVE, STATUS_RESOLVED},
	[]string{STATUS_PENDING, STATUS_RESOLVED},
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

var ALL_CMDS = []string{
	CMD_NEXT,
	CMD_ADD,
	CMD_START,
	CMD_ANNOTATE,
	CMD_STOP,
	CMD_RESOLVE,
	CMD_CONTEXT,
	CMD_MODIFY,
	CMD_EDIT,
	CMD_UNDO,
	CMD_PUSH,
	CMD_PULL,
	CMD_GIT,
	CMD_DAY,
	CMD_WEEK,
	CMD_PROJECTS,
	CMD_IMPORT_TW,
	CMD_HELP,
}
