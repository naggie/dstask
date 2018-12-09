package dstask

const (
	STATUS_PENDING   = "pending"
	STATUS_ACTIVE    = "active"
	STATUS_RESOLVED  = "resolved"
	STATUS_DELEGATED = "delegated"
	STATUS_DEFERRED  = "deferred"
	STATUS_SOMEDAY   = "someday"

	GIT_REPO   = "~/.dstask/"
	CACHE_FILE = "~/.cache/dstask/completion_cache.gob"

	// filter: P1 P2 etc
	PRIORITY_CRITICAL = 1
	PRIORITY_HIGH = 2
	PRIORITY_NORMAL = 3
	PRIORITY_LOW = 4
)
