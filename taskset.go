package dstask

// main task data structures

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type TaskSet struct {
	tasks []*Task

	// indices
	tasksByID   map[int]*Task
	tasksByUUID map[string]*Task

	// program metadata
	idsFilePath string
	repoPath    string
}

type Project struct {
	Name          string
	Tasks         int
	TasksResolved int
	// if any task is in the active state
	Active bool
	// first task created
	Created time.Time
	// last task resolved
	Resolved time.Time

	// highest non-resolved priority within project
	Priority string
}

// LoadTaskSet constructs a TaskSet from a repo path..
func LoadTaskSet(repoPath, idsFilePath string, includeResolved bool) (*TaskSet, error) {

	// Initialise an empty TaskSet
	var ts TaskSet
	ts.tasksByUUID = make(map[string]*Task)
	ts.tasksByID = make(map[int]*Task)

	ts.idsFilePath = idsFilePath
	ts.repoPath = repoPath

	// Construct our options struct by calling our passed-in TaskSetOpt functions.
	ids := LoadIds(idsFilePath)

	var statuses []string

	if includeResolved {
		// expensive to load -- resolved tasks are unbounded
		statuses = ALL_STATUSES
	} else {
		// non-resolved tasks are bounded, so it's OK to load them even if
		// some are redundant due to query. It's also important to load all
		// non-resolved tasks at once for consistent IDs in case
		// SavePendingChanges is not called...!
		statuses = NON_RESOLVED_STATUSES
	}

	for _, status := range statuses {
		dir := filepath.Join(repoPath, status)
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			if os.IsNotExist(err) {
				// Continuing here is necessary, because we do not guarantee
				// that all status directories exist on program startup.
				continue
			}
			return nil, err
		}
		for _, finfo := range files {
			path := filepath.Join(dir, finfo.Name())
			t, err := unmarshalTask(path, finfo, ids, status)
			if err != nil {
				log.Printf("error loading task: %v\n", err)
				continue
			}
			ts.LoadTask(t)
		}
	}

	// hide some tasks by default. This is useful for things like templates and
	// recurring tasks which are shown either directly or with show- commands
	for _, task := range ts.tasks {
		if StrSliceContains(HIDDEN_STATUSES, task.Status) {
			task.filtered = true
		}
	}

	return &ts, nil
}

func (ts *TaskSet) UnHide() {
	for _, task := range ts.tasks {
		if StrSliceContains(HIDDEN_STATUSES, task.Status) {
			task.filtered = false
		}
	}
}

func (ts *TaskSet) SortByCreated(dir SortByDirection) {
	switch dir {
	case Ascending:
		// Oldest first
		sort.SliceStable(ts.tasks, func(i, j int) bool { return ts.tasks[i].Created.Before(ts.tasks[j].Created) })
	case Descending:
		// Newest first
		sort.SliceStable(ts.tasks, func(i, j int) bool { return ts.tasks[i].Created.After(ts.tasks[j].Created) })
	}
}

func (ts *TaskSet) SortByPriority(dir SortByDirection) {
	switch dir {
	case Ascending:
		// P1 first
		sort.SliceStable(ts.tasks, func(i, j int) bool { return ts.tasks[i].Priority < ts.tasks[j].Priority })
	case Descending:
		// P1 last
		sort.SliceStable(ts.tasks, func(i, j int) bool { return ts.tasks[i].Priority > ts.tasks[j].Priority })
	}
}

func (ts *TaskSet) SortByResolved(dir SortByDirection) {
	switch dir {
	case Ascending:
		// Oldest resolved first
		sort.SliceStable(ts.tasks, func(i, j int) bool { return ts.tasks[i].Resolved.Before(ts.tasks[j].Resolved) })
	case Descending:
		// Newest resolved first
		sort.SliceStable(ts.tasks, func(i, j int) bool { return ts.tasks[i].Resolved.After(ts.tasks[j].Resolved) })
	}
}

// MustLoadTask is the same as LoadTask, except it exits on error.
func (ts *TaskSet) MustLoadTask(task Task) Task {
	newTask, err := ts.LoadTask(task)
	if err != nil {
		ExitFail("%s, task %s", err, task.UUID)
	}
	return newTask
}

// LoadTask adds a task to the TaskSet, but only if it has a new uuid or no uuid.
// Return annotated task.
func (ts *TaskSet) LoadTask(task Task) (Task, error) {
	task.Normalise()

	if task.UUID == "" {
		task.UUID = MustGetUUID4String()
	}

	if err := task.Validate(); err != nil {
		return Task{}, err
	}

	if ts.tasksByUUID[task.UUID] != nil {
		// load tasks, do not overwrite
		// TODO ??? (maybe return a nil pointer instead?)
		return Task{}, nil
	}

	// remove ID if already taken
	if task.ID > 0 && ts.tasksByID[task.ID] != nil {
		task.ID = 0
	}

	// pick one if task isn't resolved and ID isn't there
	if task.ID == 0 && task.Status != STATUS_RESOLVED {
		for id := 1; id <= MAX_TASKS_OPEN; id++ {
			if ts.tasksByID[id] == nil {
				task.ID = id
				break
			}
		}
	}

	if task.Created.IsZero() {
		task.Created = time.Now()
		task.WritePending = true
	}

	ts.tasks = append(ts.tasks, &task)
	ts.tasksByUUID[task.UUID] = &task
	ts.tasksByID[task.ID] = &task

	return task, nil
}

// TODO maybe this is the place to check for invalid state transitions instead
// of the main switch statement. Though, a future 3rdparty sync system could
// need this to work regardless.
func (ts *TaskSet) MustUpdateTask(task Task) {
	if err := ts.UpdateTask(task); err != nil {
		ExitFail(err.Error())
	}
}

func (ts *TaskSet) UpdateTask(task Task) error {
	task.Normalise()

	if err := task.Validate(); err != nil {
		return fmt.Errorf("%s, task %s", err, task.UUID)
	}

	if ts.tasksByUUID[task.UUID] == nil {
		return fmt.Errorf("Could not find given task to update by UUID")
	}

	if !IsValidPriority(task.Priority) {
		return fmt.Errorf("Invalid priority specified")
	}

	old := ts.tasksByUUID[task.UUID]

	if old.Status != task.Status && !IsValidStateTransition(old.Status, task.Status) {
		return fmt.Errorf("Invalid state transition: %s -> %s", old.Status, task.Status)
	}

	if old.Status != task.Status && task.Status == STATUS_RESOLVED && strings.Contains(task.Notes, "- [ ] ") {
		return fmt.Errorf("Refusing to resolve task with incomplete tasklist")
	}

	if task.Status == STATUS_RESOLVED {
		task.ID = 0
	}

	if task.Status == STATUS_RESOLVED && task.Resolved.IsZero() {
		task.Resolved = time.Now()
	}

	task.WritePending = true
	// existing pointer must point to address of new task copied
	*ts.tasksByUUID[task.UUID] = task
	return nil
}

func (ts *TaskSet) Filter(query Query) {
	for _, task := range ts.tasks {
		if !task.MatchesFilter(query) {
			task.filtered = true
		}
	}
}

func (ts *TaskSet) FilterByStatus(status string) {
	for _, task := range ts.tasks {
		if task.Status != status {
			task.filtered = true
		}
	}
}

func (ts *TaskSet) FilterOrganised() {
	for _, task := range ts.tasks {
		if len(task.Tags) > 0 || task.Project != "" {
			task.filtered = true
		}
	}
}

func (ts *TaskSet) MustGetByID(id int) Task {
	task, err := ts.GetByID(id)
	if err != nil {
		ExitFail(err.Error())
	}
	return task
}

func (ts *TaskSet) GetByID(id int) (Task, error) {
	if ts.tasksByID[id] == nil {
		return Task{}, fmt.Errorf("no open task with ID %v exists", id)
	}

	return *ts.tasksByID[id], nil
}

func (ts *TaskSet) Tasks() []Task {
	tasks := make([]Task, 0, len(ts.tasks))
	for _, task := range ts.tasks {
		if !task.filtered {
			tasks = append(tasks, *task)
		}
	}
	return tasks
}

func (ts *TaskSet) AllTasks() []Task {
	tasks := make([]Task, 0, len(ts.tasks))
	for _, task := range ts.tasks {
		tasks = append(tasks, *task)
	}
	return tasks
}

func (ts *TaskSet) GetTags() map[string]bool {
	tagset := make(map[string]bool)

	for _, task := range ts.Tasks() {
		for _, tag := range task.Tags {
			tagset[tag] = true
		}
	}

	return tagset
}

func (ts *TaskSet) GetProjects() []*Project {
	projectsMap := make(map[string]*Project)

	for _, task := range ts.Tasks() {
		name := task.Project

		if name == "" {
			continue
		}

		if projectsMap[name] == nil {
			projectsMap[name] = &Project{
				Name:     name,
				Priority: PRIORITY_LOW,
			}
		}

		project := projectsMap[name]

		project.Tasks++

		if project.Created.IsZero() || task.Created.Before(project.Created) {
			project.Created = task.Created
		}

		if task.Resolved.After(project.Resolved) {
			project.Resolved = task.Resolved
		}

		if task.Status == STATUS_RESOLVED {
			project.TasksResolved++
		}

		if task.Status == STATUS_ACTIVE {
			project.Active = true
		}

		if task.Status != STATUS_RESOLVED && task.Priority < project.Priority {
			project.Priority = task.Priority
		}
	}

	// collect keys to produce ordered output (rather than randomised)
	names := make([]string, 0, len(projectsMap))
	projects := make([]*Project, 0, len(projectsMap))

	for name := range projectsMap {
		names = append(names, name)
	}

	sort.Strings(names)

	for _, name := range names {
		projects = append(projects, projectsMap[name])
	}

	return projects
}

func (ts *TaskSet) NumTotal() int {
	return len(ts.tasks)
}

// save pending changes to disk
// TODO return files that have been added/deleted/modified/renamed so they can
// be passed to git add for performance, instead of doing git add .
func (ts *TaskSet) SavePendingChanges() {
	ids := make(IdsMap, len(ts.Tasks()))

	for _, task := range ts.tasks {
		if task.WritePending {
			task.SaveToDisk(ts.repoPath)
		}

		if task.ID > 0 {
			ids[task.UUID] = task.ID
		}
	}

	// saving generally only happens when tasks are mutated. This is OK, and
	// important. Generally the ID assignment process is deterministic such
	// that a DB is not required. However, if tasks are listed and then tasks
	// are closed or created, it can have a ripple effect such that it is
	// possible for every ID to change. Therefore, tasks must retain their IDs
	// locally. This replaced a system where tasks recorded their IDs, which
	// can create merge conflicts in some (uncommon) cases.
	ids.Save(ts.idsFilePath)
}

type SortByDirection string

const (
	Ascending  SortByDirection = "ascending"
	Descending SortByDirection = "descending"
)
