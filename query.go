package dstask

// main task data structures

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

// when referring to tasks by ID, NON_RESOLVED_STATUSES must be loaded exclusively --
// even if the filter is set to show issues that have only some statuses.
type Query struct {
	Cmd           string
	IDs           []int
	Tags          []string
	AntiTags      []string
	Project       string
	AntiProjects  []string
	Due           time.Time
	DateFilter    string    // "", after, before, on, in -- or none/any (Due is zero then)
	DueEnd        time.Time // inclusive upper bound of a range, DateFilter is "after" then
	Priority      string
	Priorities    []string // all if more than one is given (any matches), else nil
	Template      int
	Text          string
	IgnoreContext bool
	// any words after the note operator: /
	Note string
}

// reconstruct args string.
func (query Query) String() string {
	var args []string

	for _, id := range query.IDs {
		args = append(args, strconv.Itoa(id))
	}

	for _, tag := range query.Tags {
		args = append(args, "+"+tag)
	}

	for _, tag := range query.AntiTags {
		args = append(args, "-"+tag)
	}

	if query.Project != "" {
		args = append(args, "project:"+query.Project)
	}

	for _, project := range query.AntiProjects {
		args = append(args, "-project:"+project)
	}

	if !query.Due.IsZero() {
		dueArg := "due"
		if query.DateFilter != "" {
			dueArg += "." + query.DateFilter
		}
		dueArg += ":" + query.Due.Format("2006-01-02")
		args = append(args, dueArg)

		if !query.DueEnd.IsZero() {
			args = append(args, "due.before:"+query.DueEnd.Format("2006-01-02"))
		}
	} else if query.DateFilter != "" {
		args = append(args, "due:"+query.DateFilter)
	}

	args = append(args, query.priorities()...)

	if query.Template > 0 {
		args = append(args, fmt.Sprintf("template:%v", query.Template))
	}

	if query.Text != "" {
		args = append(args, "\""+query.Text+"\"")
	}

	return strings.Join(args, " ")
}

func (query Query) PrintContextDescription() {
	var envVarNotification string
	if os.Getenv("DSTASK_CONTEXT") != "" {
		envVarNotification = " (set by DSTASK_CONTEXT)"
	}

	if query.String() != "" {
		fmt.Printf("\033[33mActive context%s: %s\033[0m\n", envVarNotification, query)
	}
}

// HasOperators returns true if the query has positive or negative projects/tags,
// dueDate, priorities or template.
func (query Query) HasOperators() bool {
	return (len(query.Tags) > 0 ||
		len(query.AntiTags) > 0 ||
		query.Project != "" ||
		len(query.AntiProjects) > 0 ||
		query.Due != time.Time{} ||
		query.DateFilter != "" ||
		query.Priority != "" ||
		query.Template > 0)
}

// ParseQuery parses the raw command line typed by the user.
func ParseQuery(args ...string) Query {
	var cmd string

	var ids []int

	var tags []string

	var antiTags []string

	var project string

	var antiProjects []string

	var dueDate time.Time

	var dueEnd time.Time

	var dateFilterType string

	var priority string

	var priorities []string

	var template int

	var words []string

	var notesModeActivated bool

	var notes []string

	var ignoreContext bool

	// something other than an ID has been parsed -- accept no more IDs
	var IDsExhausted bool

	var dueDateSet bool

	for _, item := range args {
		lcItem := strings.ToLower(item)

		if notesModeActivated {
			// no more parsing syntax
			notes = append(notes, item)

			continue
		}

		if cmd == "" && StrSliceContains(ALL_CMDS, lcItem) {
			cmd = lcItem

			continue
		}

		if s, err := strconv.ParseInt(item, 10, 64); !IDsExhausted && err == nil {
			ids = append(ids, int(s))

			continue
		}

		if item == IGNORE_CONTEXT_KEYWORD {
			ignoreContext = true
		} else if item == NOTE_MODE_KEYWORD {
			notesModeActivated = true
		} else if project == "" && strings.HasPrefix(lcItem, "project:") {
			project = lcItem[8:]
			// invalid project filter, but a common mistake and it's obvious what
			// the user really means.
		} else if project == "" && strings.HasPrefix(lcItem, "+project:") {
			project = lcItem[9:]
		} else if strings.HasPrefix(lcItem, "-project:") {
			antiProjects = append(antiProjects, lcItem[9:])
		} else if strings.HasPrefix(lcItem, "due.") || strings.HasPrefix(lcItem, "due:") {
			if start, end, ok := parseDuePeriod(lcItem, time.Now()); ok {
				if dueDateSet {
					ExitFail("Query should only have one due date")
				}
				dateFilterType, dueDate, dueEnd = "after", start, end
				dueDateSet = true
			} else {
				filter, date := ParseDueDateArg(lcItem)
				start, end, isRange := dueRange(dateFilterType, dueDate, filter, date)

				switch {
				case !dueDateSet:
					dateFilterType, dueDate = filter, date
					dueDateSet = true
				case isRange && dueEnd.IsZero() && !end.Before(start):
					dateFilterType, dueDate, dueEnd = "after", start, end
				case isRange && dueEnd.IsZero():
					ExitFail("due.after must not be later than due.before")
				default:
					ExitFail("Query should only have one due date or due range")
				}
			}
		} else if strings.HasPrefix(lcItem, "template:") {
			if s, err := strconv.ParseInt(lcItem[9:], 10, 64); err == nil {
				template = int(s)
			}
		} else if len(item) > 1 && lcItem[0:1] == "+" {
			tags = append(tags, lcItem[1:])
		} else if len(item) > 1 && lcItem[0:1] == "-" {
			antiTags = append(antiTags, lcItem[1:])
		} else if IsValidPriority(item) {
			if priority == "" {
				priority = item
			}
			if !StrSliceContains(priorities, item) {
				priorities = append(priorities, item)
			}
		} else {
			words = append(words, item)
		}

		IDsExhausted = true
	}

	if len(priorities) < 2 {
		priorities = nil
	}

	return Query{
		Cmd:           cmd,
		IDs:           ids,
		Tags:          tags,
		AntiTags:      antiTags,
		Project:       project,
		AntiProjects:  antiProjects,
		DateFilter:    dateFilterType,
		Due:           dueDate,
		DueEnd:        dueEnd,
		Priority:      priority,
		Priorities:    priorities,
		Template:      template,
		Text:          strings.Join(words, " "),
		Note:          strings.Join(notes, " "),
		IgnoreContext: ignoreContext,
	}
}

// Merge applies a context to a new task. Returns new Query, does not mutate.
func (query *Query) Merge(q2 Query) Query {
	// dereference to make a copy of this query
	q := *query

	for _, tag := range q2.Tags {
		if !StrSliceContains(q.Tags, tag) {
			q.Tags = append(q.Tags, tag)
		}
	}

	for _, tag := range q2.AntiTags {
		if !StrSliceContains(q.AntiTags, tag) {
			q.AntiTags = append(q.AntiTags, tag)
		}
	}

	if q2.Project != "" {
		if q.Project != "" && q.Project != q2.Project {
			ExitFail("Could not apply q2, project conflict")
		} else {
			q.Project = q2.Project
		}
	}

	for _, project := range q2.AntiProjects {
		if !StrSliceContains(q.AntiProjects, project) {
			q.AntiProjects = append(q.AntiProjects, project)
		}
	}

	if !q2.Due.IsZero() || q2.DateFilter != "" {
		start, end, isRange := dueRange(q.DateFilter, q.Due, q2.DateFilter, q2.Due)
		isRange = isRange && q.DueEnd.IsZero() && q2.DueEnd.IsZero() && !end.Before(start)

		switch {
		case q.Due.IsZero() && q.DateFilter == "":
			q.Due, q.DueEnd, q.DateFilter = q2.Due, q2.DueEnd, q2.DateFilter
		case q.DateFilter == "any" && !q2.Due.IsZero():
			// a due date (filter) is more specific than due:any
			q.Due, q.DueEnd, q.DateFilter = q2.Due, q2.DueEnd, q2.DateFilter
		case q2.DateFilter == "any" && !q.Due.IsZero():
			// a due date (filter) is more specific than due:any
		case isRange:
			q.Due, q.DueEnd, q.DateFilter = start, end, "after"
		case q.Due != q2.Due || q.DueEnd != q2.DueEnd ||
			!sameDateFilter(q.DateFilter, q2.DateFilter):
			ExitFail("Could not apply q2, date filter conflict")
		}
	}

	if q2Priorities := q2.priorities(); len(q2Priorities) > 0 {
		switch {
		case q.Priority == "":
			q.Priority, q.Priorities = q2.Priority, q2.Priorities
		case len(q2Priorities) > 1 && isSubset(q.priorities(), q2Priorities):
			// the query narrows down the priorities of the context
		default:
			ExitFail("Could not apply q2, priority conflict")
		}
	}

	return q
}

// HasWildcardTags returns true if the query contains a tag filter like +team:*
// Such tags are filters only and cannot be added to or removed from tasks.
func (query Query) HasWildcardTags() bool {
	return slices.ContainsFunc(query.Tags, isWildcardTag) ||
		slices.ContainsFunc(query.AntiTags, isWildcardTag)
}

func isWildcardTag(tag string) bool {
	return strings.HasSuffix(tag, "*")
}

// withoutWildcardTags drops wildcard tags (e.g. from the context) when
// creating a task.
func withoutWildcardTags(tags []string) []string {
	var result []string

	for _, tag := range tags {
		if !isWildcardTag(tag) {
			result = append(result, tag)
		}
	}

	return result
}

// dueRange combines a due.after and a due.before filter (in any order) to an
// inclusive range.
func dueRange(
	filterA string, dateA time.Time, filterB string, dateB time.Time,
) (start, end time.Time, ok bool) {
	switch {
	case filterA == "after" && filterB == "before":
		return dateA, dateB, true
	case filterA == "before" && filterB == "after":
		return dateB, dateA, true
	}

	return start, end, false
}

// withoutDueFilter drops a due date range (e.g. due:this-week), due:none and
// due:any from the context when creating a task. These are filters only and
// have no single due date.
func withoutDueFilter(ctx Query) Query {
	if !ctx.DueEnd.IsZero() || ctx.DateFilter == "none" || ctx.DateFilter == "any" {
		ctx.Due, ctx.DueEnd, ctx.DateFilter = time.Time{}, time.Time{}, ""
	}

	return ctx
}

// sameDateFilter treats due:X, due.on:X and due.in:X as the same filter.
func sameDateFilter(a, b string) bool {
	normalize := func(filter string) string {
		if filter == "on" || filter == "in" {
			return ""
		}

		return filter
	}

	return normalize(a) == normalize(b)
}

// priorities returns all priorities of the query, a task matches any of them.
func (query Query) priorities() []string {
	if len(query.Priorities) > 0 {
		return query.Priorities
	}

	if query.Priority != "" {
		return []string{query.Priority}
	}

	return nil
}

// withoutPriorityFilter drops several priorities (e.g. P0 P1 from the context)
// after merging the context when creating a task. They are a filter only. A
// single priority is kept and set on the task.
func withoutPriorityFilter(q Query) Query {
	if len(q.Priorities) > 1 {
		q.Priority, q.Priorities = "", nil
	}

	return q
}

func errMultiplePriorities(cmd string) error {
	return fmt.Errorf(
		"cannot use more than one priority with %s command, quote text like \"level P1\"",
		cmd,
	)
}

// isSubset reports whether every item of a is in b.
func isSubset(a, b []string) bool {
	for _, item := range a {
		if !StrSliceContains(b, item) {
			return false
		}
	}

	return true
}
