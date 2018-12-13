package dstask

// maintains an association of ID to UUID to allow the user to address tasks
// quickly. IDs are positive integers, and recycled when tasks are resolved.
// The roster stores its state on disk, in the cache. IDs start from 1.

// Getting an ID from a UUID is cheap, but the reverse is more expensive. This
// fits the use pattern of listing IDs given a UUID, and operating on one task
// at a time.

// NOT THREAD SAFE

import (
	"os"
)

type IDRoster struct {
	IDs map[string]uint64
	// IDs that can be re-used after task resolution
	recycledIDs []uint64
	// incremented for new ID when required
	lastId uint64
}


// TODO make directories for ID_ROSTER_FILE
func LoadIDRoster() *IDRoster {
	filePath := MustExpandHome(ID_ROSTER_FILE)

	roster := &IDRoster{
		IDs: make(map[string]uint64),
	}

	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		MustReadGob(filePath, roster)
	}

	return roster
}

func (r *IDRoster) GetId(uuid string) uint64 {
	if r.IDs[uuid] == 0 {
		if len(r.recycledIDs) > 0 {
			id := r.recycledIDs[0]
			r.recycledIDs = r.recycledIDs[1:]
			r.IDs[uuid] = id
		} else {
			r.lastId = r.lastId + 1
			r.IDs[uuid] = r.lastId
		}
	}

	return r.IDs[uuid]
}

func (r *IDRoster) GetUuid(id uint64) string {
	for k, v := range(r.IDs) {
		if v == id {
			return k
		}
	}

	return ""
}

func (r *IDRoster) RecycleId(uuid string) {
	if id := r.IDs[uuid]; id != 0 {
		delete(r.IDs, uuid)
		r.recycledIDs = append(r.recycledIDs, id)
	}
}

func (r *IDRoster) SaveToDisk() {
	filePath := MustExpandHome(ID_ROSTER_FILE)
	MustWriteGob(filePath, r)
}
