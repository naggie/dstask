package dstask

const (
	STATUS_PENDING   = "pending"
	STATUS_ACTIVE    = "active"
	STATUS_RESOLVED  = "resolved"
	STATUS_DELEGATED = "delegated"
	STATUS_DEFERRED  = "deferred"
	STATUS_SOMEDAY   = "someday"
	STATUS_RECURRING = "recurring"  // tentative

	GIT_REPO   = "~/.dstask/"
	CACHE_FILE = "~/.cache/dstask/completion_cache.gob"

	// filter: P1 P2 etc
	PRIORITY_CRITICAL = "P1"
	PRIORITY_HIGH     = "P2"
	PRIORITY_NORMAL   = "P3"
	PRIORITY_LOW      = "P4"
)
